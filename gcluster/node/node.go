package node

import (
	"context"
	"fmt"
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gmodules"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/gtransport"
	"gitee.com/monobytes/gcore/gutils/gcall"
	"gitee.com/monobytes/gcore/gwrap/info"
	"gitee.com/monobytes/gcore/internal/transporter/node"
	"golang.org/x/sync/errgroup"
	"sync"
	"sync/atomic"
)

type HookHandler func(proxy *Proxy)

type serviceEntity struct {
	name     string      // 服务名称;用于定位服务发现
	desc     interface{} // 服务描述(grpc为desc描述对象; rpcx为服务路径)
	provider interface{} // 服务提供者
}

type Node struct {
	gmodules.Base
	opts        *options
	ctx         context.Context
	cancel      context.CancelFunc
	state       atomic.Int32
	evtPool     *sync.Pool
	reqPool     *sync.Pool
	router      *Router
	trigger     *Trigger
	proxy       *Proxy
	services    []*serviceEntity
	instances   []*gregistry.ServiceInstance
	linker      *node.Server
	fnChan      chan func()
	scheduler   *Scheduler
	transporter gtransport.Server
	wg          *sync.WaitGroup
	rw          sync.RWMutex
	hooks       map[gcluster.Hook][]HookHandler
}

func NewNode(opts ...Option) *Node {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	n := &Node{}
	n.opts = o
	n.ctx, n.cancel = context.WithCancel(o.ctx)
	n.proxy = newProxy(n)
	n.router = newRouter(n)
	n.trigger = newTrigger(n)
	n.scheduler = newScheduler(n)
	n.hooks = make(map[gcluster.Hook][]HookHandler)
	n.services = make([]*serviceEntity, 0)
	n.instances = make([]*gregistry.ServiceInstance, 0)
	n.fnChan = make(chan func(), 4096)
	n.state.Store(int32(gcluster.Shut))
	n.wg = &sync.WaitGroup{}
	n.evtPool = &sync.Pool{New: func() interface{} {
		return &event{
			ctx:  context.Background(),
			node: n,
		}
	}}
	n.reqPool = &sync.Pool{New: func() interface{} {
		return &request{
			ctx:     context.Background(),
			node:    n,
			message: &gcluster.Message{},
		}
	}}

	return n
}

// Name 组件名称
func (n *Node) Name() string {
	return n.opts.name
}

// Init 初始化节点
func (n *Node) Init() {
	if n.opts.id == "" {
		glog.Fatal("instance id can not be empty")
	}

	if n.opts.name == "" {
		glog.Fatal("instance name can not be empty")
	}

	if n.opts.codec == nil {
		glog.Fatal("codec modules is not injected")
	}

	if n.opts.locator == nil {
		glog.Fatal("locator modules is not injected")
	}

	if n.opts.registry == nil {
		glog.Fatal("registry modules is not injected")
	}

	n.runHookFunc(gcluster.Init)
}

// Start 启动节点
func (n *Node) Start() {
	if !n.state.CompareAndSwap(int32(gcluster.Shut), int32(gcluster.Work)) {
		return
	}

	n.startLinkServer()

	n.startTransportServer()

	n.registerServiceInstances()

	n.proxy.watch()

	go n.dispatch()

	n.printInfo()

	n.runHookFunc(gcluster.Start)
}

// Close 关闭节点
func (n *Node) Close() {
	if !n.state.CompareAndSwap(int32(gcluster.Work), int32(gcluster.Hang)) {
		if !n.state.CompareAndSwap(int32(gcluster.Busy), int32(gcluster.Hang)) {
			return
		}
	}

	n.refreshServiceInstances()

	n.runHookFunc(gcluster.Close)

	n.wg.Wait()
}

// Destroy 销毁节点服务器
func (n *Node) Destroy() {
	if !n.state.CompareAndSwap(int32(gcluster.Hang), int32(gcluster.Shut)) {
		return
	}

	n.runHookFunc(gcluster.Destroy)

	n.deregisterServiceInstances()

	n.stopLinkServer()

	n.stopTransportServer()

	n.router.close()

	n.trigger.close()

	close(n.fnChan)

	n.cancel()
}

// Proxy 获取节点代理
func (n *Node) Proxy() *Proxy {
	return n.proxy
}

// 分发处理消息
func (n *Node) dispatch() {
	for {
		select {
		case evt, ok := <-n.trigger.receive():
			if !ok {
				return
			}
			gcall.Call(func() {
				n.trigger.handle(evt)
			})
		case req, ok := <-n.router.receive():
			if !ok {
				return
			}
			gcall.Call(func() {
				n.router.handle(req)
			})
		case handle, ok := <-n.fnChan:
			if !ok {
				return
			}
			gcall.Call(func() {
				handle()
				n.doneWait()
			})
		}
	}
}

// 启动连接服务器
func (n *Node) startLinkServer() {
	linker, err := node.NewServer(n.opts.addr, &provider{node: n})
	if err != nil {
		glog.Fatalf("link server create failed: %v", err)
	}

	n.linker = linker

	go func() {
		if err = n.linker.Start(); err != nil {
			glog.Fatalf("link server start failed: %v", err)
		}
	}()
}

// 停止连接服务器
func (n *Node) stopLinkServer() {
	if err := n.linker.Stop(); err != nil {
		glog.Errorf("link server stop failed: %v", err)
	}
}

// 启动传输服务器
func (n *Node) startTransportServer() {
	if n.opts.transporter == nil {
		return
	}

	n.opts.transporter.SetDefaultDiscovery(n.opts.registry)

	if len(n.services) == 0 {
		return
	}

	transporter, err := n.opts.transporter.NewServer()
	if err != nil {
		glog.Fatalf("transport server create failed: %v", err)
	}

	n.transporter = transporter

	for _, entity := range n.services {
		if err = n.transporter.RegisterService(entity.desc, entity.provider); err != nil {
			glog.Fatalf("register service failed: %v", err)
		}
	}

	go func() {
		if err = n.transporter.Start(); err != nil {
			glog.Fatalf("transport server start failed: %v", err)
		}
	}()
}

// 停止传输服务器
func (n *Node) stopTransportServer() {
	if n.transporter == nil {
		return
	}

	if err := n.transporter.Stop(); err != nil {
		glog.Errorf("transport server stop failed: %v", err)
	}
}

// 注册服务实例
func (n *Node) registerServiceInstances() {
	routes := make([]gregistry.Route, 0, len(n.router.routes))
	events := make([]int, 0, len(n.trigger.events))

	for _, entity := range n.router.routes {
		routes = append(routes, gregistry.Route{
			ID:       entity.route,
			Stateful: entity.stateful,
			Internal: entity.internal,
		})
	}

	for evt := range n.trigger.events {
		events = append(events, int(evt))
	}

	n.instances = append(n.instances, &gregistry.ServiceInstance{
		ID:       n.opts.id,
		Name:     gcluster.Node.String(),
		Kind:     gcluster.Node.String(),
		Alias:    n.opts.name,
		State:    n.getState().String(),
		Routes:   routes,
		Events:   events,
		Endpoint: n.linker.Endpoint().String(),
		Weight:   n.opts.weight,
	})

	if n.transporter != nil {
		services := make([]string, 0, len(n.services))
		for _, item := range n.services {
			services = append(services, item.name)
		}

		n.instances = append(n.instances, &gregistry.ServiceInstance{
			ID:       n.opts.id,
			Name:     gcluster.Mesh.String(),
			Kind:     gcluster.Mesh.String(),
			Alias:    n.opts.name,
			State:    n.getState().String(),
			Services: services,
			Endpoint: n.transporter.Endpoint().String(),
			Weight:   n.opts.weight,
		})
	}

	if err := n.doRegisterServiceInstances(); err != nil {
		glog.Fatalf("register cluster instances failed: %v", err)
	}
}

// 刷新服务实例状态
func (n *Node) refreshServiceInstances() {
	if err := n.doRefreshServiceInstances(); err != nil {
		glog.Errorf("refresh cluster instances failed: %v", err)
	}
}

// 解注册服务实例
func (n *Node) deregisterServiceInstances() {
	eg, ctx := errgroup.WithContext(n.ctx)
	for i := range n.instances {
		instance := n.instances[i]
		eg.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
			defer cancel()
			return n.opts.registry.Deregister(ctx, instance)
		})
	}

	if err := eg.Wait(); err != nil {
		glog.Errorf("deregister cluster instances failed: %v", err)
	}
}

// 执行注册操作
func (n *Node) doRegisterServiceInstances() error {
	eg, ctx := errgroup.WithContext(n.ctx)

	for i := range n.instances {
		instance := n.instances[i]
		eg.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
			defer cancel()
			return n.opts.registry.Register(ctx, instance)
		})
	}

	return eg.Wait()
}

// 执行刷新实例状态操作
func (n *Node) doRefreshServiceInstances() error {
	for _, instance := range n.instances {
		instance.State = n.getState().String()
	}

	return n.doRegisterServiceInstances()
}

// 获取状态
func (n *Node) getState() gcluster.State {
	return gcluster.State(n.state.Load())
}

// 更新状态
func (n *Node) setState(state gcluster.State) error {
	n.state.Store(int32(state))

	return n.doRefreshServiceInstances()
}

// 执行钩子函数
func (n *Node) runHookFunc(hook gcluster.Hook) {
	n.rw.RLock()

	if handlers, ok := n.hooks[hook]; ok {
		wg := &sync.WaitGroup{}
		wg.Add(len(handlers))

		for i := range handlers {
			handler := handlers[i]
			gcall.Go(func() {
				handler(n.proxy)
				wg.Done()
			})
		}

		n.rw.RUnlock()

		wg.Wait()
	} else {
		n.rw.RUnlock()
	}
}

// 添加钩子监听器
func (n *Node) addHookListener(hook gcluster.Hook, handler HookHandler) {
	switch hook {
	case gcluster.Destroy:
		n.rw.Lock()
		n.hooks[hook] = append(n.hooks[hook], handler)
		n.rw.Unlock()
	default:
		if n.getState() == gcluster.Shut {
			n.hooks[hook] = append(n.hooks[hook], handler)
		} else {
			glog.Warnf("server is working, can't add hook handler")
		}
	}
}

// 添加服务提供者
func (n *Node) addServiceProvider(name string, desc, provider any) {
	if n.getState() == gcluster.Shut {
		n.services = append(n.services, &serviceEntity{
			name:     name,
			desc:     desc,
			provider: provider,
		})
	} else {
		glog.Warnf("server is working, can't add service provider")
	}
}

// 打印组件信息
func (n *Node) printInfo() {
	infos := make([]string, 0)
	infos = append(infos, fmt.Sprintf("Name: %s", n.Name()))
	infos = append(infos, fmt.Sprintf("Link: %s", n.linker.ExposeAddr()))
	infos = append(infos, fmt.Sprintf("Codec: %s", n.opts.codec.Name()))
	infos = append(infos, fmt.Sprintf("Locator: %s", n.opts.locator.Name()))
	infos = append(infos, fmt.Sprintf("Registry: %s", n.opts.registry.Name()))

	if n.opts.encryptor != nil {
		infos = append(infos, fmt.Sprintf("Encryptor: %s", n.opts.encryptor.Name()))
	} else {
		infos = append(infos, "Encryptor: -")
	}

	if n.opts.transporter != nil {
		infos = append(infos, fmt.Sprintf("Transporter: %s", n.opts.transporter.Name()))
	} else {
		infos = append(infos, "Transporter: -")
	}

	info.PrintBoxInfo("Node", infos...)
}

func (n *Node) doneWait() {
	if n.getState() != gcluster.Shut {
		n.wg.Done()
	}
}

func (n *Node) addWait() {
	if n.getState() != gcluster.Shut {
		n.wg.Add(1)
	}
}

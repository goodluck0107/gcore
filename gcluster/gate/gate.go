package gate

import (
	"context"
	"fmt"
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gmodules"
	"gitee.com/monobytes/gcore/gnetwork"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/gsession"
	"gitee.com/monobytes/gcore/gwrap/info"
	"gitee.com/monobytes/gcore/gwrap/net"
	"gitee.com/monobytes/gcore/internal/transporter/gate"
	"sync"
	"sync/atomic"
)

type Gate struct {
	gmodules.Base
	opts     *options
	ctx      context.Context
	cancel   context.CancelFunc
	state    atomic.Int32
	proxy    *proxy
	instance *gregistry.ServiceInstance
	session  *gsession.Session
	linker   *gate.Server
	wg       *sync.WaitGroup
}

func NewGate(opts ...Option) *Gate {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	g := &Gate{}
	g.opts = o
	g.ctx, g.cancel = context.WithCancel(o.ctx)
	g.proxy = newProxy(g)
	g.session = gsession.NewSession()
	g.state.Store(int32(gcluster.Shut))
	g.wg = &sync.WaitGroup{}

	return g
}

// Name 组件名称
func (g *Gate) Name() string {
	return g.opts.name
}

// Init 初始化
func (g *Gate) Init() {
	if g.opts.id == "" {
		glog.Fatal("instance id can not be empty")
	}

	if g.opts.server == nil {
		glog.Fatal("server modules is not injected")
	}

	if g.opts.locator == nil {
		glog.Fatal("locator modules is not injected")
	}

	if g.opts.registry == nil {
		glog.Fatal("registry modules is not injected")
	}
}

// Start 启动组件
func (g *Gate) Start() {
	if !g.state.CompareAndSwap(int32(gcluster.Shut), int32(gcluster.Work)) {
		return
	}

	g.startNetworkServer()

	g.startLinkerServer()

	g.registerServiceInstance()

	g.proxy.watch()

	g.printInfo()
}

// Close 关闭节点
func (g *Gate) Close() {
	if !g.state.CompareAndSwap(int32(gcluster.Work), int32(gcluster.Hang)) {
		if !g.state.CompareAndSwap(int32(gcluster.Busy), int32(gcluster.Hang)) {
			return
		}
	}

	g.refreshServiceInstance()

	g.wg.Wait()
}

// Destroy 销毁组件
func (g *Gate) Destroy() {
	if !g.state.CompareAndSwap(int32(gcluster.Hang), int32(gcluster.Shut)) {
		return
	}

	g.deregisterServiceInstance()

	g.stopNetworkServer()

	g.stopLinkerServer()

	g.cancel()
}

// 启动网络服务器
func (g *Gate) startNetworkServer() {
	g.opts.server.OnConnect(g.handleConnect)
	g.opts.server.OnDisconnect(g.handleDisconnect)
	g.opts.server.OnReceive(g.handleReceive)

	if err := g.opts.server.Start(); err != nil {
		glog.Fatalf("network server start failed: %v", err)
	}
}

// 停止网关服务器
func (g *Gate) stopNetworkServer() {
	if err := g.opts.server.Stop(); err != nil {
		glog.Errorf("network server stop failed: %v", err)
	}
}

// 处理连接打开
func (g *Gate) handleConnect(conn gnetwork.Conn) {
	g.wg.Add(1)

	g.session.AddConn(conn)

	cid, uid := conn.ID(), conn.UID()

	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	g.proxy.trigger(ctx, gcluster.Connect, cid, uid)
	cancel()
}

// 处理断开连接
func (g *Gate) handleDisconnect(conn gnetwork.Conn) {
	g.session.RemConn(conn)

	if cid, uid := conn.ID(), conn.UID(); uid != 0 {
		ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
		_ = g.proxy.unbindGate(ctx, cid, uid)
		g.proxy.trigger(ctx, gcluster.Disconnect, cid, uid)
		cancel()
	} else {
		ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
		g.proxy.trigger(ctx, gcluster.Disconnect, cid, uid)
		cancel()
	}

	g.wg.Done()
}

// 处理接收到的消息
func (g *Gate) handleReceive(conn gnetwork.Conn, data []byte) {
	cid, uid := conn.ID(), conn.UID()
	ctx, cancel := context.WithTimeout(g.ctx, g.opts.timeout)
	g.proxy.deliver(ctx, cid, uid, data)
	cancel()
}

// 启动传输服务器
func (g *Gate) startLinkerServer() {
	transporter, err := gate.NewServer(g.opts.addr, &provider{gate: g})
	if err != nil {
		glog.Fatalf("link server create failed: %v", err)
	}

	g.linker = transporter

	go func() {
		if err = g.linker.Start(); err != nil {
			glog.Errorf("link server start failed: %v", err)
		}
	}()
}

// 停止传输服务器
func (g *Gate) stopLinkerServer() {
	if err := g.linker.Stop(); err != nil {
		glog.Errorf("link server stop failed: %v", err)
	}
}

// 注册服务实例
func (g *Gate) registerServiceInstance() {
	g.instance = &gregistry.ServiceInstance{
		ID:       g.opts.id,
		Name:     gcluster.Gate.String(),
		Kind:     gcluster.Gate.String(),
		Alias:    g.opts.name,
		State:    g.getState().String(),
		Weight:   g.opts.weight,
		Endpoint: g.linker.Endpoint().String(),
	}

	ctx, cancel := context.WithTimeout(g.ctx, defaultTimeout)
	defer cancel()

	if err := g.opts.registry.Register(ctx, g.instance); err != nil {
		glog.Fatalf("register cluster instance failed: %v", err)
	}
}

// 刷新服务实例状态
func (g *Gate) refreshServiceInstance() {
	if g.instance == nil {
		return
	}

	g.instance.State = g.getState().String()

	ctx, cancel := context.WithTimeout(g.ctx, defaultTimeout)
	defer cancel()

	if err := g.opts.registry.Register(ctx, g.instance); err != nil {
		glog.Fatalf("refresh cluster instance failed: %v", err)
	}
}

// 解注册服务实例
func (g *Gate) deregisterServiceInstance() {
	ctx, cancel := context.WithTimeout(g.ctx, defaultTimeout)
	defer cancel()

	if err := g.opts.registry.Deregister(ctx, g.instance); err != nil {
		glog.Errorf("deregister cluster instance failed: %v", err)
	}
}

// 获取状态
func (g *Gate) getState() gcluster.State {
	return gcluster.State(g.state.Load())
}

// 打印组件信息
func (g *Gate) printInfo() {
	infos := make([]string, 0)
	infos = append(infos, fmt.Sprintf("Name: %s", g.Name()))
	infos = append(infos, fmt.Sprintf("Link: %s", g.linker.ExposeAddr()))
	infos = append(infos, fmt.Sprintf("Server: [%s] %s", g.opts.server.Protocol(), net.FulfillAddr(g.opts.server.Addr())))
	infos = append(infos, fmt.Sprintf("Locator: %s", g.opts.locator.Name()))
	infos = append(infos, fmt.Sprintf("Registry: %s", g.opts.registry.Name()))

	info.PrintBoxInfo("Gate", infos...)
}

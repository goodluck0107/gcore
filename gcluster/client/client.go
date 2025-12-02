package client

import (
	"context"
	"fmt"
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gerrors"
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/gmodules"
	"github.com/goodluck0107/gcore/gnetwork"
	"github.com/goodluck0107/gcore/gpacket"
	"github.com/goodluck0107/gcore/gutils/gcall"
	"github.com/goodluck0107/gcore/gwrap/info"
	"sync"
	"sync/atomic"
)

type HookHandler func(proxy *Proxy)

type RouteHandler func(ctx *Context)

type EventHandler func(conn *Conn)

type Client struct {
	gmodules.Base
	opts                *options
	ctx                 context.Context
	cancel              context.CancelFunc
	routes              map[int32][]RouteHandler
	events              map[gcluster.Event][]EventHandler
	defaultRouteHandler RouteHandler
	proxy               *Proxy
	state               int32
	conns               sync.Map
	rw                  sync.RWMutex
	hooks               map[gcluster.Hook][]HookHandler
}

func NewClient(opts ...Option) *Client {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	c := &Client{}
	c.opts = o
	c.proxy = newProxy(c)
	c.routes = make(map[int32][]RouteHandler)
	c.events = make(map[gcluster.Event][]EventHandler)
	c.hooks = make(map[gcluster.Hook][]HookHandler)
	c.ctx, c.cancel = context.WithCancel(o.ctx)
	c.state = int32(gcluster.Shut)

	return c
}

// Name 组件名称
func (c *Client) Name() string {
	return c.opts.name
}

// Init 初始化节点
func (c *Client) Init() {
	if c.opts.client == nil {
		glog.Fatal("client plugin is not injected")
	}

	if c.opts.codec == nil {
		glog.Fatal("codec plugin is not injected")
	}

	c.runHookFunc(gcluster.Init)
}

// Start 启动组件
func (c *Client) Start() {
	c.setState(gcluster.Work)

	c.opts.client.OnDisconnect(c.handleDisconnect)
	c.opts.client.OnReceive(c.handleReceive)

	c.printInfo()

	c.runHookFunc(gcluster.Start)
}

// Destroy 销毁组件
func (c *Client) Destroy() {
	c.setState(gcluster.Shut)

	c.runHookFunc(gcluster.Destroy)
}

// Proxy 获取节点代理
func (c *Client) Proxy() *Proxy {
	return c.proxy
}

// 处理断开连接
func (c *Client) handleDisconnect(conn gnetwork.Conn) {
	val, ok := c.conns.Load(conn)
	if !ok {
		return
	}

	c.conns.Delete(conn)

	handlers, ok := c.events[gcluster.Disconnect]
	if !ok {
		return
	}

	for _, handler := range handlers {
		gcall.Call(func() {
			handler(val.(*Conn))
		})
	}
}

// 处理接收到的消息
func (c *Client) handleReceive(conn gnetwork.Conn, data []byte) {
	val, ok := c.conns.Load(conn)
	if !ok {
		return
	}

	message, err := gpacket.UnpackMessage(data)
	if err != nil {
		glog.Errorf("unpack message failed: %v", err)
		return
	}

	handlers, ok := c.routes[message.Route]
	if ok {
		for _, handler := range handlers {
			gcall.Call(func() {
				handler(&Context{
					ctx:     context.Background(),
					conn:    val.(*Conn),
					message: message,
				})
			})
		}
	} else if c.defaultRouteHandler != nil {
		c.defaultRouteHandler(&Context{
			ctx:     context.Background(),
			conn:    val.(*Conn),
			message: message,
		})
	} else {
		glog.Debugf("route handler is not registered, route: %v", message.Route)
	}
}

// 拨号
func (c *Client) dial(opts ...DialOption) (*Conn, error) {
	if c.getState() == gcluster.Shut {
		return nil, gerrors.ErrClientShut
	}

	o := &dialOptions{attrs: make(map[string]any)}
	for _, opt := range opts {
		opt(o)
	}

	conn, err := c.opts.client.Dial(o.addr)
	if err != nil {
		return nil, err
	}

	cc := &Conn{conn: conn, client: c}

	for key, value := range o.attrs {
		cc.SetAttr(key, value)
	}

	c.conns.Store(conn, cc)

	if handlers, ok := c.events[gcluster.Connect]; ok {
		for _, handler := range handlers {
			gcall.Call(func() {
				handler(cc)
			})
		}
	}

	return cc, nil
}

// 添加路由处理器
func (c *Client) addRouteHandler(route int32, handler RouteHandler) {
	if c.getState() == gcluster.Shut {
		c.routes[route] = append(c.routes[route], handler)
	} else {
		glog.Warnf("client is working, can't add route handler")
	}
}

// 默认路由处理器
func (c *Client) setDefaultRouteHandler(handler RouteHandler) {
	if c.getState() == gcluster.Shut {
		c.defaultRouteHandler = handler
	} else {
		glog.Warnf("client is working, can't set default route handler")
	}
}

// 添加事件处理器
func (c *Client) addEventListener(event gcluster.Event, handler EventHandler) {
	if c.getState() == gcluster.Shut {
		c.events[event] = append(c.events[event], handler)
	} else {
		glog.Warnf("client is working, can't add event handler")
	}
}

// 添加钩子监听器
func (c *Client) addHookListener(hook gcluster.Hook, handler HookHandler) {
	switch hook {
	case gcluster.Destroy:
		c.rw.Lock()
		c.hooks[hook] = append(c.hooks[hook], handler)
		c.rw.Unlock()
	default:
		if c.getState() == gcluster.Shut {
			c.hooks[hook] = append(c.hooks[hook], handler)
		} else {
			glog.Warnf("server is working, can't add hook handler")
		}
	}
}

// 设置状态
func (c *Client) setState(state gcluster.State) {
	atomic.StoreInt32(&c.state, int32(state))
}

// 获取状态
func (c *Client) getState() gcluster.State {
	return gcluster.State(atomic.LoadInt32(&c.state))
}

// 执行钩子函数
func (c *Client) runHookFunc(hook gcluster.Hook) {
	c.rw.RLock()

	if handlers, ok := c.hooks[hook]; ok {
		wg := &sync.WaitGroup{}
		wg.Add(len(handlers))

		for i := range handlers {
			handler := handlers[i]
			gcall.Go(func() {
				handler(c.proxy)
				wg.Done()
			})
		}

		c.rw.RUnlock()

		wg.Wait()
	} else {
		c.rw.RUnlock()
	}
}

// 打印组件信息
func (c *Client) printInfo() {
	infos := make([]string, 0)
	infos = append(infos, fmt.Sprintf("Name: %s", c.Name()))
	infos = append(infos, fmt.Sprintf("Codec: %s", c.opts.codec.Name()))
	infos = append(infos, fmt.Sprintf("Protocol: %s", c.opts.client.Protocol()))

	if c.opts.encryptor != nil {
		infos = append(infos, fmt.Sprintf("Encryptor: %s", c.opts.encryptor.Name()))
	} else {
		infos = append(infos, "Encryptor: -")
	}

	info.PrintBoxInfo("Client", infos...)
}

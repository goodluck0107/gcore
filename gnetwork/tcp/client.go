package tcp

import (
	"gitee.com/monobytes/gcore/gnetwork"
	"net"
	"sync/atomic"
)

type client struct {
	opts              *clientOptions             // 配置
	id                int64                      // 连接ID
	connectHandler    gnetwork.ConnectHandler    // 连接打开hook函数
	disconnectHandler gnetwork.DisconnectHandler // 连接关闭hook函数
	receiveHandler    gnetwork.ReceiveHandler    // 接收消息hook函数
}

var _ gnetwork.Client = &client{}

func NewClient(opts ...ClientOption) gnetwork.Client {
	o := defaultClientOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &client{opts: o}
}

// Dial 拨号连接
func (c *client) Dial(addr ...string) (gnetwork.Conn, error) {
	var address string
	if len(addr) > 0 && addr[0] != "" {
		address = addr[0]
	} else {
		address = c.opts.addr
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTimeout(tcpAddr.Network(), tcpAddr.String(), c.opts.timeout)
	if err != nil {
		return nil, err
	}

	return newClientConn(c, atomic.AddInt64(&c.id, 1), conn), nil
}

// Protocol 协议
func (c *client) Protocol() string {
	return protocol
}

// OnConnect 监听连接打开
func (c *client) OnConnect(handler gnetwork.ConnectHandler) {
	c.connectHandler = handler
}

// OnDisconnect 监听连接关闭
func (c *client) OnDisconnect(handler gnetwork.DisconnectHandler) {
	c.disconnectHandler = handler
}

// OnReceive 监听接收到消息
func (c *client) OnReceive(handler gnetwork.ReceiveHandler) {
	c.receiveHandler = handler
}

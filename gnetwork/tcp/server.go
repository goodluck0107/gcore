package tcp

import (
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/gnetwork"
	"net"
	"time"
)

type server struct {
	opts              *serverOptions             // 配置
	listener          net.Listener               // 监听器
	connMgr           *serverConnMgr             // 连接管理器
	startHandler      gnetwork.StartHandler      // 服务器启动hook函数
	stopHandler       gnetwork.CloseHandler      // 服务器关闭hook函数
	connectHandler    gnetwork.ConnectHandler    // 连接打开hook函数
	disconnectHandler gnetwork.DisconnectHandler // 连接关闭hook函数
	receiveHandler    gnetwork.ReceiveHandler    // 接收消息hook函数
}

var _ gnetwork.Server = &server{}

func NewServer(opts ...ServerOption) gnetwork.Server {
	o := defaultServerOptions()
	for _, opt := range opts {
		opt(o)
	}

	s := &server{}
	s.opts = o
	s.connMgr = newServerConnMgr(s)

	return s
}

// Addr 监听地址
func (s *server) Addr() string {
	return s.opts.addr
}

// Start 启动服务器
func (s *server) Start() error {
	if err := s.init(); err != nil {
		return err
	}

	if s.startHandler != nil {
		s.startHandler()
	}

	go s.serve()

	return nil
}

// Stop 关闭服务器
func (s *server) Stop() error {
	if err := s.listener.Close(); err != nil {
		return err
	}

	s.connMgr.close()

	if s.stopHandler != nil {
		s.stopHandler()
	}

	return nil
}

// Protocol 协议
func (s *server) Protocol() string {
	return protocol
}

// OnStart 监听服务器启动
func (s *server) OnStart(handler gnetwork.StartHandler) {
	s.startHandler = handler
}

// OnStop 监听服务器关闭
func (s *server) OnStop(handler gnetwork.CloseHandler) {
	s.stopHandler = handler
}

// OnConnect 监听连接打开
func (s *server) OnConnect(handler gnetwork.ConnectHandler) {
	s.connectHandler = handler
}

// OnDisconnect 监听连接关闭
func (s *server) OnDisconnect(handler gnetwork.DisconnectHandler) {
	s.disconnectHandler = handler
}

// OnReceive 监听接收到消息
func (s *server) OnReceive(handler gnetwork.ReceiveHandler) {
	s.receiveHandler = handler
}

// 初始化TCP服务器
func (s *server) init() error {
	addr, err := net.ResolveTCPAddr("tcp", s.opts.addr)
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return err
	}

	s.listener = ln

	return nil
}

// 等待连接
func (s *server) serve() {
	var tempDelay time.Duration

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				glog.Warnf("tcp accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			glog.Warnf("tcp accept error: %v", err)
			return
		}

		tempDelay = 0

		if err = s.connMgr.allocate(conn); err != nil {
			glog.Errorf("connection allocate error: %v", err)
			_ = conn.Close()
		}
	}
}

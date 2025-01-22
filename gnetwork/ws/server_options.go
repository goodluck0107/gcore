package ws

import (
	"gitee.com/monobytes/gcore/getc"
	"net/http"
	"time"
)

const (
	defaultServerAddr               = ":3553"
	defaultServerPath               = "/"
	defaultServerMaxConnNum         = 5000
	defaultServerCheckOrigin        = "*"
	defaultServerHandshakeTimeout   = "10s"
	defaultServerHeartbeatInterval  = "10s"
	defaultServerHeartbeatMechanism = "resp"
)

const (
	defaultServerAddrKey               = "etc.network.ws.server.addr"
	defaultServerPathKey               = "etc.network.ws.server.path"
	defaultServerMaxConnNumKey         = "etc.network.ws.server.maxConnNum"
	defaultServerCheckOriginsKey       = "etc.network.ws.server.origins"
	defaultServerKeyFileKey            = "etc.network.ws.server.keyFile"
	defaultServerCertFileKey           = "etc.network.ws.server.certFile"
	defaultServerHandshakeTimeoutKey   = "etc.network.ws.server.handshakeTimeout"
	defaultServerHeartbeatIntervalKey  = "etc.network.ws.server.heartbeatInterval"
	defaultServerHeartbeatMechanismKey = "etc.network.ws.server.heartbeatMechanism"
)

const (
	RespHeartbeat HeartbeatMechanism = "resp" // 响应式心跳
	TickHeartbeat HeartbeatMechanism = "tick" // 主动定时心跳
)

type HeartbeatMechanism string

type ServerOption func(o *serverOptions)

type CheckOriginFunc func(r *http.Request) bool

type serverOptions struct {
	addr               string             // 监听地址
	maxConnNum         int                // 最大连接数
	certFile           string             // 证书文件
	keyFile            string             // 秘钥文件
	path               string             // 路径，默认为"/"
	checkOrigin        CheckOriginFunc    // 跨域检测
	handshakeTimeout   time.Duration      // 握手超时时间，默认10s
	heartbeatInterval  time.Duration      // 心跳间隔时间，默认10s
	heartbeatMechanism HeartbeatMechanism // 心跳机制，默认resp
}

func defaultServerOptions() *serverOptions {
	origins := getc.Get(defaultServerCheckOriginsKey, []string{defaultServerCheckOrigin}).Strings()
	checkOrigin := func(r *http.Request) bool {
		if len(origins) == 0 {
			return false
		}

		origin := r.Header.Get("Origin")
		for _, v := range origins {
			if v == defaultServerCheckOrigin || origin == v {
				return true
			}
		}

		return false
	}

	return &serverOptions{
		addr:               getc.Get(defaultServerAddrKey, defaultServerAddr).String(),
		maxConnNum:         getc.Get(defaultServerMaxConnNumKey, defaultServerMaxConnNum).Int(),
		path:               getc.Get(defaultServerPathKey, defaultServerPath).String(),
		checkOrigin:        checkOrigin,
		keyFile:            getc.Get(defaultServerKeyFileKey).String(),
		certFile:           getc.Get(defaultServerCertFileKey).String(),
		handshakeTimeout:   getc.Get(defaultServerHandshakeTimeoutKey, defaultServerHandshakeTimeout).Duration(),
		heartbeatInterval:  getc.Get(defaultServerHeartbeatIntervalKey, defaultServerHeartbeatInterval).Duration(),
		heartbeatMechanism: HeartbeatMechanism(getc.Get(defaultServerHeartbeatMechanismKey, defaultServerHeartbeatMechanism).String()),
	}
}

// WithServerListenAddr 设置监听地址
func WithServerListenAddr(addr string) ServerOption {
	return func(o *serverOptions) { o.addr = addr }
}

// WithServerMaxConnNum 设置连接的最大连接数
func WithServerMaxConnNum(maxConnNum int) ServerOption {
	return func(o *serverOptions) { o.maxConnNum = maxConnNum }
}

// WithServerPath 设置Websocket的连接路径
func WithServerPath(path string) ServerOption {
	return func(o *serverOptions) { o.path = path }
}

// WithServerCredentials 设置证书和秘钥
func WithServerCredentials(certFile, keyFile string) ServerOption {
	return func(o *serverOptions) { o.keyFile, o.certFile = keyFile, certFile }
}

// WithServerCheckOrigin 设置Websocket跨域检测函数
func WithServerCheckOrigin(checkOrigin CheckOriginFunc) ServerOption {
	return func(o *serverOptions) { o.checkOrigin = checkOrigin }
}

// WithServerHandshakeTimeout 设置握手超时时间
func WithServerHandshakeTimeout(handshakeTimeout time.Duration) ServerOption {
	return func(o *serverOptions) { o.handshakeTimeout = handshakeTimeout }
}

// WithServerHeartbeatInterval 设置心跳检测间隔时间
func WithServerHeartbeatInterval(heartbeatInterval time.Duration) ServerOption {
	return func(o *serverOptions) { o.heartbeatInterval = heartbeatInterval }
}

// WithServerHeartbeatMechanism 设置心跳机制
func WithServerHeartbeatMechanism(heartbeatMechanism HeartbeatMechanism) ServerOption {
	return func(o *serverOptions) { o.heartbeatMechanism = heartbeatMechanism }
}

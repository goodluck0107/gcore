package kcp

import (
	"gitee.com/monobytes/gcore/getc"
	"time"
)

const (
	defaultServerAddr               = ":3553"
	defaultServerMaxMsgLen          = 1024
	defaultServerMaxConnNum         = 5000
	defaultServerHeartbeatInterval  = "10s"
	defaultServerHeartbeatMechanism = "resp"
)

const (
	defaultServerAddrKey               = "etc.network.kcp.server.addr"
	defaultServerMaxMsgLenKey          = "etc.network.kcp.server.maxMsgLen"
	defaultServerMaxConnNumKey         = "etc.network.kcp.server.maxConnNum"
	defaultServerHeartbeatIntervalKey  = "etc.network.kcp.server.heartbeatInterval"
	defaultServerHeartbeatMechanismKey = "etc.network.kcp.server.heartbeatMechanism"
)

const (
	RespHeartbeat HeartbeatMechanism = "resp" // 响应式心跳
	TickHeartbeat HeartbeatMechanism = "tick" // 主动定时心跳
)

type HeartbeatMechanism string

type ServerOption func(o *serverOptions)

type serverOptions struct {
	addr               string             // 监听地址
	maxConnNum         int                // 最大连接数
	heartbeatInterval  time.Duration      // 心跳检测间隔时间，默认10s
	heartbeatMechanism HeartbeatMechanism // 心跳机制，默认resp
}

func defaultServerOptions() *serverOptions {
	return &serverOptions{
		addr:               getc.Get(defaultServerAddrKey, defaultServerAddr).String(),
		maxConnNum:         getc.Get(defaultServerMaxConnNumKey, defaultServerMaxConnNum).Int(),
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

// WithServerHeartbeatInterval 设置心跳检测间隔时间
func WithServerHeartbeatInterval(heartbeatInterval time.Duration) ServerOption {
	return func(o *serverOptions) { o.heartbeatInterval = heartbeatInterval }
}

// WithServerHeartbeatMechanism 设置心跳机制
func WithServerHeartbeatMechanism(heartbeatMechanism HeartbeatMechanism) ServerOption {
	return func(o *serverOptions) { o.heartbeatMechanism = heartbeatMechanism }
}

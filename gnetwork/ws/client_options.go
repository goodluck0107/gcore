package ws

import (
	"gitee.com/monobytes/gcore/getc"
	"time"
)

const (
	defaultClientDialUrl           = "ws://127.0.0.1:3553"
	defaultClientHandshakeTimeout  = "10s"
	defaultClientHeartbeatInterval = "10s"
)

const (
	defaultClientDialUrlKey           = "etc.network.ws.client.url"
	defaultClientHandshakeTimeoutKey  = "etc.network.ws.client.handshakeTimeout"
	defaultClientHeartbeatIntervalKey = "etc.network.ws.client.heartbeatInterval"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	url               string        // 拨号地址
	msgType           string        // 默认消息类型，text | binary
	handshakeTimeout  time.Duration // 握手超时时间
	heartbeatInterval time.Duration // 心跳间隔时间，默认10s
}

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		url:               getc.Get(defaultClientDialUrlKey, defaultClientDialUrl).String(),
		handshakeTimeout:  getc.Get(defaultClientHandshakeTimeoutKey, defaultClientHandshakeTimeout).Duration(),
		heartbeatInterval: getc.Get(defaultClientHeartbeatIntervalKey, defaultClientHeartbeatInterval).Duration(),
	}
}

// WithClientDialUrl 设置拨号链接
func WithClientDialUrl(url string) ClientOption {
	return func(o *clientOptions) { o.url = url }
}

// WithClientHandshakeTimeout 设置握手超时时间
func WithClientHandshakeTimeout(handshakeTimeout time.Duration) ClientOption {
	return func(o *clientOptions) { o.handshakeTimeout = handshakeTimeout }
}

// WithClientHeartbeatInterval 设置心跳间隔时间
func WithClientHeartbeatInterval(heartbeatInterval time.Duration) ClientOption {
	return func(o *clientOptions) { o.heartbeatInterval = heartbeatInterval }
}

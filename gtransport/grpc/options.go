package grpc

import (
	"github.com/goodluck0107/gcore/getc"
	"github.com/goodluck0107/gcore/gregistry"
	"github.com/goodluck0107/gcore/gtransport/grpc/internal/client"
	"github.com/goodluck0107/gcore/gtransport/grpc/internal/server"
	"google.golang.org/grpc"
)

const (
	defaultServerAddr     = ":0" // 默认服务器地址
	defaultClientPoolSize = 10   // 默认客户端连接池大小
)

const (
	defaultServerAddrKey       = "etc.transport.grpc.server.addr"
	defaultServerKeyFileKey    = "etc.transport.grpc.server.keyFile"
	defaultServerCertFileKey   = "etc.transport.grpc.server.certFile"
	defaultClientPoolSizeKey   = "etc.transport.grpc.client.poolSize"
	defaultClientCertFileKey   = "etc.transport.grpc.client.certFile"
	defaultClientServerNameKey = "etc.transport.grpc.client.serverName"
)

type Option func(o *options)

type options struct {
	server server.Options
	client client.Options
}

func defaultOptions() *options {
	opts := &options{}
	opts.server.Addr = getc.Get(defaultServerAddrKey, defaultServerAddr).String()
	opts.server.KeyFile = getc.Get(defaultServerKeyFileKey).String()
	opts.server.CertFile = getc.Get(defaultServerCertFileKey).String()
	opts.client.CertFile = getc.Get(defaultClientCertFileKey).String()
	opts.client.ServerName = getc.Get(defaultClientServerNameKey).String()

	return opts
}

// WithServerListenAddr 设置服务器监听地址
func WithServerListenAddr(addr string) Option {
	return func(o *options) { o.server.Addr = addr }
}

// WithServerCredentials 设置服务器证书和秘钥
func WithServerCredentials(certFile, keyFile string) Option {
	return func(o *options) { o.server.KeyFile, o.server.CertFile = keyFile, certFile }
}

// WithServerOptions 设置服务器选项
func WithServerOptions(opts ...grpc.ServerOption) Option {
	return func(o *options) { o.server.ServerOpts = opts }
}

// WithClientCredentials 设置客户端证书和校验域名
func WithClientCredentials(certFile string, serverName string) Option {
	return func(o *options) { o.client.CertFile, o.client.ServerName = certFile, serverName }
}

// WithClientDiscovery 设置客户端服务发现组件
func WithClientDiscovery(discovery gregistry.Discovery) Option {
	return func(o *options) { o.client.Discovery = discovery }
}

// WithClientDialOptions 设置客户端拨号选项
func WithClientDialOptions(opts ...grpc.DialOption) Option {
	return func(o *options) { o.client.DialOpts = opts }
}

package etcd

import (
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/getc"
	"go.etcd.io/etcd/client/v3"
	"time"
)

const (
	defaultAddr        = "127.0.0.1:2379"
	defaultDialTimeout = "5s"
	defaultPath        = "/config"
	defaultMode        = gconfig.ReadOnly
)

const (
	defaultAddrsKey       = "etc.config.etcd.addrs"
	defaultDialTimeoutKey = "etc.config.etcd.dialTimeout"
	defaultPathKey        = "etc.config.etcd.path"
	defaultModeKey        = "etc.config.etcd.mode"
)

type Option func(o *options)

type options struct {
	// 客户端连接地址
	// 内建客户端配置，默认为[]string{"localhost:2379"}
	addrs []string

	// 客户端拨号超时时间
	// 内建客户端配置，默认为5秒
	dialTimeout time.Duration

	// 外部客户端
	// 外部客户端配置，存在外部客户端时，优先使用外部客户端，默认为nil
	client *clientv3.Client

	// 路径
	// 默认为 /config
	path string

	// 读写模式
	// 支持read-only、write-only和read-write三种模式，默认为read-only模式
	mode gconfig.Mode
}

func defaultOptions() *options {
	return &options{
		addrs:       getc.Get(defaultAddrsKey, []string{defaultAddr}).Strings(),
		dialTimeout: getc.Get(defaultDialTimeoutKey, defaultDialTimeout).Duration(),
		path:        getc.Get(defaultPathKey, defaultPath).String(),
		mode:        gconfig.Mode(getc.Get(defaultModeKey, defaultMode).String()),
	}
}

// WithAddrs 设置客户端连接地址
func WithAddrs(addrs ...string) Option {
	return func(o *options) { o.addrs = addrs }
}

// WithDialTimeout 设置客户端拨号超时时间
func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *options) { o.dialTimeout = dialTimeout }
}

// WithClient 设置外部客户端
func WithClient(client *clientv3.Client) Option {
	return func(o *options) { o.client = client }
}

// WithPath 设置命名空间
func WithPath(path string) Option {
	return func(o *options) { o.path = path }
}

// WithMode 设置读写模式
func WithMode(mode gconfig.Mode) Option {
	return func(o *options) { o.mode = mode }
}

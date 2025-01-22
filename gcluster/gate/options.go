package gate

import (
	"context"
	"gitee.com/monobytes/gcore/getc"
	"gitee.com/monobytes/gcore/glocate"
	"gitee.com/monobytes/gcore/gutils/guuid"
	"time"

	"gitee.com/monobytes/gcore/gnetwork"
	"gitee.com/monobytes/gcore/gregistry"
)

const (
	defaultName    = "gate"          // 默认名称
	defaultAddr    = ":0"            // 连接器监听地址
	defaultTimeout = 3 * time.Second // 默认超时时间
	defaultWeight  = 1               // 默认权重
)

const (
	defaultIDKey      = "etc.cluster.gate.id"
	defaultNameKey    = "etc.cluster.gate.name"
	defaultAddrKey    = "etc.cluster.gate.addr"
	defaultTimeoutKey = "etc.cluster.gate.timeout"
	defaultWeightKey  = "etc.cluster.gate.weight"
)

type Option func(o *options)

type options struct {
	ctx      context.Context    // 上下文
	id       string             // 实例ID
	name     string             // 实例名称
	addr     string             // 监听地址
	timeout  time.Duration      // RPC调用超时时间
	weight   int                // 权重
	server   gnetwork.Server    // 网关服务器
	locator  glocate.Locator    // 用户定位器
	registry gregistry.Registry // 服务注册器
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		addr:    defaultAddr,
		timeout: defaultTimeout,
		weight:  defaultWeight,
	}

	if id := getc.Get(defaultIDKey).String(); id != "" {
		opts.id = id
	} else {
		opts.id = guuid.UUID()
	}

	if name := getc.Get(defaultNameKey).String(); name != "" {
		opts.name = name
	}

	if addr := getc.Get(defaultAddrKey).String(); addr != "" {
		opts.addr = addr
	}

	if timeout := getc.Get(defaultTimeoutKey).Duration(); timeout > 0 {
		opts.timeout = timeout
	}

	if weight := getc.Get(defaultWeightKey).Int(); weight > 0 {
		opts.weight = weight
	}

	return opts
}

// WithID 设置实例ID
func WithID(id string) Option {
	return func(o *options) { o.id = id }
}

// WithName 设置实例名称
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithServer 设置服务器
func WithServer(server gnetwork.Server) Option {
	return func(o *options) { o.server = server }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置用户定位器
func WithLocator(locator glocate.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithRegistry 设置服务注册器
func WithRegistry(r gregistry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// WithWeight 设置权重
func WithWeight(weight int) Option {
	return func(o *options) { o.weight = weight }
}

package mesh

import (
	"context"
	"gitee.com/monobytes/gcore/gcrypto"
	"gitee.com/monobytes/gcore/gencoding"
	"gitee.com/monobytes/gcore/getc"
	"gitee.com/monobytes/gcore/glocate"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/gtransport"
	"gitee.com/monobytes/gcore/gutils/guuid"
	"time"
)

const (
	defaultName    = "mesh"          // 默认节点名称
	defaultCodec   = "proto"         // 默认编解码器名称
	defaultTimeout = 3 * time.Second // 默认超时时间
	defaultWeight  = 1               // 默认权重
)

const (
	defaultIDKey      = "etc.cluster.mesh.id"
	defaultNameKey    = "etc.cluster.mesh.name"
	defaultCodecKey   = "etc.cluster.mesh.codec"
	defaultTimeoutKey = "etc.cluster.mesh.timeout"
	defaultWeightKey  = "etc.cluster.mesh.weight"
)

type Option func(o *options)

type options struct {
	id          string                 // 实例ID
	name        string                 // 实例名称
	ctx         context.Context        // 上下文
	codec       gencoding.Codec        // 编解码器
	timeout     time.Duration          // RPC调用超时时间
	locator     glocate.Locator        // 用户定位器
	registry    gregistry.Registry     // 服务注册器
	encryptor   gcrypto.Encryptor      // 消息加密器
	transporter gtransport.Transporter // 消息传输器
	weight      int                    // 权重
}

func defaultOptions() *options {
	opts := &options{
		ctx:     context.Background(),
		name:    defaultName,
		codec:   gencoding.Invoke(defaultCodec),
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

	if codec := getc.Get(defaultCodecKey).String(); codec != "" {
		opts.codec = gencoding.Invoke(codec)
	}

	if timeout := getc.Get(defaultTimeoutKey).Duration(); timeout > 0 {
		opts.timeout = timeout
	}

	if weight := getc.Get(defaultWeightKey).Int(); weight > 0 {
		opts.weight = weight
	}

	return opts
}

// WithName 设置实例名称
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithCodec 设置编解码器
func WithCodec(codec gencoding.Codec) Option {
	return func(o *options) { o.codec = codec }
}

// WithContext 设置上下文
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithTimeout 设置RPC调用超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.timeout = timeout }
}

// WithLocator 设置定位器
func WithLocator(locator glocate.Locator) Option {
	return func(o *options) { o.locator = locator }
}

// WithRegistry 设置服务注册器
func WithRegistry(r gregistry.Registry) Option {
	return func(o *options) { o.registry = r }
}

// WithEncryptor 设置消息加密器
func WithEncryptor(encryptor gcrypto.Encryptor) Option {
	return func(o *options) { o.encryptor = encryptor }
}

// WithTransporter 设置消息传输器
func WithTransporter(transporter gtransport.Transporter) Option {
	return func(o *options) { o.transporter = transporter }
}

// WithWeight 设置权重
func WithWeight(weight int) Option {
	return func(o *options) { o.weight = weight }
}

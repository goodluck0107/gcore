package tencent

import (
	"gitee.com/monobytes/gcore/getc"
	"gitee.com/monobytes/gcore/glog"
)

const (
	defaultLevel          = glog.InfoLevel
	defaultStdout         = true
	defaultSyncout        = true
	defaultTimeFormat     = "2006/01/02 15:04:05.000000"
	defaultCallerFullPath = true
)

const (
	defaultLevelKey          = "etc.log.level"
	defaultTimeFormatKey     = "etc.log.timeFormat"
	defaultStackLevelKey     = "etc.log.stackLevel"
	defaultStdoutKey         = "etc.log.stdout"
	defaultSyncoutKey        = "etc.log.syncout"
	defaultCallerFullPathKey = "etc.log.callerFullPath"
)

const (
	tencentEndpointKey        = "etc.log.tencent.endpoint"
	tencentAccessKeyIDKey     = "etc.log.tencent.accessKeyID"
	tencentAccessKeySecretKey = "etc.log.tencent.accessKeySecret"
	tencentTopicIDKey         = "etc.log.tencent.topicID"
	tencentLevelKey           = "etc.log.tencent.level"
	tencentTimeFormatKey      = "etc.log.tencent.timeFormat"
	tencentStackLevelKey      = "etc.log.tencent.stackLevel"
	tencentStdoutKey          = "etc.log.tencent.stdout"
	tencentSyncoutKey         = "etc.log.tencent.syncout"
	tencentCallerFullPathKey  = "etc.log.tencent.callerFullPath"
)

type Option func(o *options)

type options struct {
	topicID         string     // 腾讯云CLS主题ID
	endpoint        string     // 腾讯云CLS服务域名，公网使用公网域名，内网使用私网域名
	accessKeyID     string     // 腾讯云CLS访问密钥ID
	accessKeySecret string     // 腾讯云CLS访问密钥密码
	stdout          bool       // 是否输出到终端，debug模式下默认输出到终端
	syncout         bool       // 是否同步输出到远端，debug模式下默认不输出到远端
	level           glog.Level // 输出的最低日志级别，默认Info
	stackLevel      glog.Level // 堆栈的最低输出级别，默认不输出堆栈
	timeFormat      string     // 时间格式，标准库时间格式，默认2006/01/02 15:04:05.000000
	callerSkip      int        // 调用者跳过的层级深度，默认为0
	callerFullPath  bool       // 是否启用调用文件全路径，默认全路径
}

func defaultOptions() *options {
	opts := &options{
		level:          defaultLevel,
		stdout:         defaultStdout,
		syncout:        defaultSyncout,
		timeFormat:     defaultTimeFormat,
		callerFullPath: defaultCallerFullPath,
	}

	level := getc.Get(tencentLevelKey, getc.Get(defaultLevelKey).String()).String()
	if lvl := glog.ParseLevel(level); lvl != glog.NoneLevel {
		opts.level = lvl
	}

	timeFormat := getc.Get(tencentTimeFormatKey, getc.Get(defaultTimeFormatKey).String()).String()
	if timeFormat != "" {
		opts.timeFormat = timeFormat
	}

	stackLevel := getc.Get(tencentStackLevelKey, getc.Get(defaultStackLevelKey).String()).String()
	if lvl := glog.ParseLevel(stackLevel); lvl != glog.NoneLevel {
		opts.stackLevel = lvl
	}

	opts.stdout = getc.Get(tencentStdoutKey, getc.Get(defaultStdoutKey, defaultStdout).Bool()).Bool()
	opts.syncout = getc.Get(tencentSyncoutKey, getc.Get(defaultSyncoutKey, defaultSyncout).Bool()).Bool()
	opts.callerFullPath = getc.Get(tencentCallerFullPathKey, getc.Get(defaultCallerFullPathKey, defaultCallerFullPath).Bool()).Bool()
	opts.endpoint = getc.Get(tencentEndpointKey).String()
	opts.accessKeyID = getc.Get(tencentAccessKeyIDKey).String()
	opts.accessKeySecret = getc.Get(tencentAccessKeySecretKey).String()
	opts.topicID = getc.Get(tencentTopicIDKey).String()

	return opts
}

// WithTopicID 设置主题ID
func WithTopicID(topicID string) Option {
	return func(o *options) { o.topicID = topicID }
}

// WithEndpoint 设置端口
func WithEndpoint(endpoint string) Option {
	return func(o *options) { o.endpoint = endpoint }
}

// WithAccessKeyID 设置访问密钥ID
func WithAccessKeyID(accessKeyID string) Option {
	return func(o *options) { o.accessKeyID = accessKeyID }
}

// WithAccessKeySecret 设置访问密钥密码
func WithAccessKeySecret(accessKeySecret string) Option {
	return func(o *options) { o.accessKeySecret = accessKeySecret }
}

// WithStdout 设置是否输出到终端
func WithStdout(enable bool) Option {
	return func(o *options) { o.stdout = enable }
}

// WithSyncout 设置是否同步输出到远端
func WithSyncout(enable bool) Option {
	return func(o *options) { o.syncout = enable }
}

// WithLevel 设置输出的最低日志级别
func WithLevel(level glog.Level) Option {
	return func(o *options) { o.level = level }
}

// WithStackLevel 设置堆栈的最小输出级别
func WithStackLevel(level glog.Level) Option {
	return func(o *options) { o.stackLevel = level }
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(format string) Option {
	return func(o *options) { o.timeFormat = format }
}

// WithCallerSkip 设置调用者跳过的层级深度
func WithCallerSkip(skip int) Option {
	return func(o *options) { o.callerSkip = skip }
}

// WithCallerFullPath 设置是否启用调用文件全路径
func WithCallerFullPath(enable bool) Option {
	return func(o *options) { o.callerFullPath = enable }
}

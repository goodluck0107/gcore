package gpprof

import (
	"github.com/goodluck0107/gcore/getc"
)

const (
	defaultAddr = ":0" // 监听地址
)

const (
	defaultAddrKey = "etc.pprof.addr"
)

type Option func(o *options)

type options struct {
	addr string // 监听地址
}

func defaultOptions() *options {
	opts := &options{
		addr: defaultAddr,
	}

	if addr := getc.Get(defaultAddrKey).String(); addr != "" {
		opts.addr = addr
	}

	return opts
}

// WithAddr 设置连接地址
func WithAddr(addr string) Option {
	return func(o *options) { o.addr = addr }
}

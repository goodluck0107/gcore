package file

import (
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/gconfig/file/core"
	"gitee.com/monobytes/gcore/glog"
)

const Name = core.Name

type Source struct {
	opts *options
}

func NewSource(opts ...Option) gconfig.Source {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.path == "" {
		glog.Fatal("no config file path specified")
	}

	return core.NewSource(o.path, o.mode)
}

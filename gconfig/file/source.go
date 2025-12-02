package file

import (
	"github.com/goodluck0107/gcore/gconfig"
	"github.com/goodluck0107/gcore/gconfig/file/core"
	"github.com/goodluck0107/gcore/glog"
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

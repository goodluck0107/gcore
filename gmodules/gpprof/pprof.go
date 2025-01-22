package gpprof

import (
	"fmt"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gmodules"
	"gitee.com/monobytes/gcore/gwrap/info"
	xnet "gitee.com/monobytes/gcore/gwrap/net"
	"net/http"
	_ "net/http/pprof"
)

var _ gmodules.Module = &PProf{}

type PProf struct {
	gmodules.Base
	opts *options
}

func NewPProf(opts ...Option) *PProf {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	return &PProf{opts: o}
}

func (*PProf) Name() string {
	return "mpprof"
}

func (p *PProf) Start() {
	listenAddr, exposeAddr, err := xnet.ParseAddr(p.opts.addr)
	if err != nil {
		glog.Fatalf("mpprof addr parse failed: %v", err)
	}

	go func() {
		if err := http.ListenAndServe(listenAddr, nil); err != nil {
			glog.Fatalf("mpprof server start failed: %v", err)
		}
	}()

	info.PrintBoxInfo("PProf",
		fmt.Sprintf("Url: http://%s/debug/pprof/", exposeAddr),
	)
}

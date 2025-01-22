package main

import (
	"gitee.com/monobytes/gcore/gengine/hextech"
	"gitee.com/monobytes/gcore/gmodules/ghttp"
	"gitee.com/monobytes/gcore/gmodules/gpprof"
)

func main() {
	engine := hextech.NewEngine()
	engine.Injection(ghttp.NewServer())
	engine.Injection(gpprof.NewPProf())
	engine.Up()
}

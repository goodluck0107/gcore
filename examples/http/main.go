package main

import (
	"gitee.com/monobytes/gcore/gengine/hextech"
	"gitee.com/monobytes/gcore/gmodules/ghttp"
	"gitee.com/monobytes/gcore/gmodules/gpprof"
)

func main() {
	engine := hextech.NewEngine()
	module := ghttp.NewServer()

	initApp(module.Proxy())

	engine.Injection(module)
	engine.Injection(gpprof.NewPProf())
	engine.Up()
}

func initApp(proxy *ghttp.Proxy) {
	router := proxy.Router()
	router.Get("/", func(ctx ghttp.Context) error {
		return ctx.Success("Hello World")
	})
}

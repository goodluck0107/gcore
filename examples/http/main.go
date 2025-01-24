package main

import (
	"gitee.com/monobytes/gcore/examples/http/service"
	"gitee.com/monobytes/gcore/examples/protocol/pb"
	"gitee.com/monobytes/gcore/gengine/hextech"
	"gitee.com/monobytes/gcore/gmodules/ghttp"
	"gitee.com/monobytes/gcore/gmodules/gpprof"
	"gitee.com/monobytes/gcore/gprotocol/handler"
	"gitee.com/monobytes/gcore/gprotocol/interfaces"
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
	interfaces.SetClusterProxy(proxy)
	pb.AddGreeterServerProvider(service.NewGreeter())
	router := proxy.Router()
	router.Get("/", func(ctx ghttp.Context) error {
		return ctx.Success("Hello World")
	})
	router.AddHandlers(handler.GetGRPCHandlersMgr())
}

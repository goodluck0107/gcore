package main

import (
	"github.com/goodluck0107/gcore/examples/http/service"
	"github.com/goodluck0107/gcore/examples/protocol/pb"
	"github.com/goodluck0107/gcore/gengine/hextech"
	"github.com/goodluck0107/gcore/gmodules/ghttp"
	"github.com/goodluck0107/gcore/gmodules/gpprof"
	"github.com/goodluck0107/gcore/gprotocol/handler"
	"github.com/goodluck0107/gcore/gprotocol/interfaces"
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

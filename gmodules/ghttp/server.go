package ghttp

import (
	"fmt"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gmodules"
	"gitee.com/monobytes/gcore/gmodules/ghttp/swagger"
	"gitee.com/monobytes/gcore/gwrap/info"
	xnet "gitee.com/monobytes/gcore/gwrap/net"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"strings"
)

type Server struct {
	gmodules.Base
	opts  *options
	app   *fiber.App
	proxy *Proxy
}

func NewServer(opts ...Option) *Server {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	s := &Server{}
	s.opts = o
	s.proxy = newProxy(s)
	s.app = fiber.New(fiber.Config{
		ServerHeader: o.name,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	s.app.Use(logger.New())

	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	if s.opts.corsOpts.Enable {
		s.app.Use(cors.New(cors.Config{
			AllowOrigins:        s.opts.corsOpts.AllowOrigins,
			AllowMethods:        s.opts.corsOpts.AllowMethods,
			AllowHeaders:        s.opts.corsOpts.AllowHeaders,
			AllowCredentials:    s.opts.corsOpts.AllowCredentials,
			ExposeHeaders:       s.opts.corsOpts.ExposeHeaders,
			MaxAge:              s.opts.corsOpts.MaxAge,
			AllowPrivateNetwork: s.opts.corsOpts.AllowPrivateNetwork,
		}))
	}

	if s.opts.swagOpts.Enable {
		s.app.Use(swagger.New(swagger.Config{
			Title:    s.opts.swagOpts.Title,
			BasePath: s.opts.swagOpts.BasePath,
			FilePath: s.opts.swagOpts.FilePath,
		}))
	}

	for i := range o.middlewares {
		switch handler := o.middlewares[i].(type) {
		case Handler:
			s.app.Use(func(ctx fiber.Ctx) error {
				return handler(&context{Ctx: ctx})
			})
		case fiber.Handler:
			s.app.Use(handler)
		}
	}

	return s
}

// Name 组件名称
func (s *Server) Name() string {
	return s.opts.name
}

// Init 初始化组件
func (s *Server) Init() {}

// Proxy 获取HTTP代理API
func (s *Server) Proxy() *Proxy {
	return s.proxy
}

// Start 启动组件
func (s *Server) Start() {
	listenAddr, exposeAddr, err := xnet.ParseAddr(s.opts.addr)
	if err != nil {
		glog.Fatalf("http addr parse failed: %v", err)
	}

	if s.opts.transporter != nil && s.opts.registry != nil {
		s.opts.transporter.SetDefaultDiscovery(s.opts.registry)
	}

	s.printInfo(exposeAddr)

	go func() {
		if err = s.app.Listen(listenAddr, fiber.ListenConfig{
			CertFile:              s.opts.certFile,
			CertKeyFile:           s.opts.keyFile,
			DisableStartupMessage: true,
			EnablePrintRoutes:     true,
		}); err != nil {
			glog.Fatal("http server startup failed: %v", err)
		}
	}()
}

func (s *Server) printInfo(addr string) {
	infos := make([]string, 0, 3)
	infos = append(infos, fmt.Sprintf("Name: %s", s.Name()))

	var baseUrl string
	if s.opts.certFile != "" && s.opts.keyFile != "" {
		baseUrl = fmt.Sprintf("https://%s", addr)
	} else {
		baseUrl = fmt.Sprintf("http://%s", addr)
	}

	infos = append(infos, fmt.Sprintf("Url: %s", baseUrl))

	if s.opts.swagOpts.Enable {
		infos = append(infos, fmt.Sprintf("Swagger: %s/%s", baseUrl, strings.TrimPrefix(s.opts.swagOpts.BasePath, "/")))
	}

	if s.opts.registry != nil {
		infos = append(infos, fmt.Sprintf("Registry: %s", s.opts.registry.Name()))
	} else {
		infos = append(infos, "Registry: -")
	}

	if s.opts.transporter != nil {
		infos = append(infos, fmt.Sprintf("Transporter: %s", s.opts.transporter.Name()))
	} else {
		infos = append(infos, "Transporter: -")
	}

	info.PrintBoxInfo("Http", infos...)
}

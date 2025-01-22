package hextech

import (
	"context"
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/getc"
	"gitee.com/monobytes/gcore/geventbus"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gmodules"
	"gitee.com/monobytes/gcore/gtask"
	"gitee.com/monobytes/gcore/gutils/gcall"
	"gitee.com/monobytes/gcore/gutils/gfile"
	"gitee.com/monobytes/gcore/gwrap/info"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const (
	defaultHextechPIDKey                 = "etc.engine.pid"                 // 进程文件路径
	defaultHextechShutdownMaxWaitTimeKey = "etc.engine.shutdownMaxWaitTime" // 容器关闭最大等待时间
)

type Engine struct {
	mods []gmodules.Module
}

func NewEngine() *Engine {
	return &Engine{}
}

// Injection 给引擎注入模块
func (c *Engine) Injection(mods ...gmodules.Module) {
	c.mods = append(c.mods, mods...)
}

// Up 启动引擎
func (c *Engine) Up() {
	c.doSaveProcessID()

	c.doPrintFrameworkInfo()

	c.doInitModules()

	c.doStartModules()

	c.doWaitSystemSignal()

	c.doCloseModules()

	c.doDestroyModules()

	c.doClearInnerModules()
}

// 初始化所有组件
func (c *Engine) doInitModules() {
	for _, mod := range c.mods {
		mod.Init()
	}
}

// 启动所有组件
func (c *Engine) doStartModules() {
	for _, mod := range c.mods {
		mod.Start()
	}
}

// 关闭所有组件
func (c *Engine) doCloseModules() {
	g := gcall.NewGoroutines()

	for _, mod := range c.mods {
		g.Add(mod.Close)
	}

	g.Run(context.Background(), getc.Get(defaultHextechShutdownMaxWaitTimeKey).Duration())
}

// 销毁所有组件
func (c *Engine) doDestroyModules() {
	g := gcall.NewGoroutines()

	for _, mod := range c.mods {
		g.Add(mod.Destroy)
	}

	g.Run(context.Background(), 5*time.Second)
}

// 等待系统信号
func (c *Engine) doWaitSystemSignal() {
	sig := make(chan os.Signal)

	switch runtime.GOOS {
	case `windows`:
		signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	default:
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	}

	s := <-sig

	signal.Stop(sig)

	glog.Warnf("process got signal %v, engine will close", s)
}

// 清理所有模块
func (c *Engine) doClearInnerModules() {
	if err := geventbus.Close(); err != nil {
		glog.Warnf("eventbus close failed: %v", err)
	}

	gtask.Release()

	gconfig.Close()

	getc.Close()

	glog.Close()
}

// 保存进程号
func (c *Engine) doSaveProcessID() {
	filename := getc.Get(defaultHextechPIDKey).String()
	if filename == "" {
		return
	}

	if err := gfile.WriteFile(filename, []byte(strconv.Itoa(syscall.Getpid()))); err != nil {
		glog.Fatalf("pid save failed: %v", err)
	}
}

// 打印框架信息
func (c *Engine) doPrintFrameworkInfo() {
	info.PrintFrameworkInfo()

	info.PrintGlobalInfo()
}

package getc

import (
	"github.com/goodluck0107/gcore/gconfig"
	"github.com/goodluck0107/gcore/gconfig/file/core"
	"github.com/goodluck0107/gcore/genv"
	"github.com/goodluck0107/gcore/gflag"
	"github.com/goodluck0107/gcore/gwrap/value"
)

// etc主要被当做项目启动配置存在；常用于集群配置、服务组件配置等。
// etc只能通过配置文件进行配置；并且无法通过master管理服进行修改。
// 如想在业务使用配置，推荐使用config配置中心进行实现。
// config配置中心的配置信息可通过master管理服进行动态修改。

const (
	gcoreEtcEnvName = "gcore_ETC"
	gcoreEtcArgName = "etc"
	defaultEtcPath  = "./etc"
)

var globalConfigurator gconfig.Configurator

func init() {
	path := genv.Get(gcoreEtcEnvName, defaultEtcPath).String()
	path = gflag.String(gcoreEtcArgName, path)
	globalConfigurator = gconfig.NewConfigurator(gconfig.WithSources(core.NewSource(path, gconfig.ReadOnly)))
}

// SetConfigurator 设置配置器
func SetConfigurator(configurator gconfig.Configurator) {
	if globalConfigurator != nil {
		globalConfigurator.Close()
	}

	globalConfigurator = configurator
}

// GetConfigurator 获取配置器
func GetConfigurator() gconfig.Configurator {
	return globalConfigurator
}

// Has 是否存在配置
func Has(pattern string) bool {
	return globalConfigurator.Has(pattern)
}

// Get 获取配置值
func Get(pattern string, def ...interface{}) value.Value {
	return globalConfigurator.Get(pattern, def...)
}

// Set 设置配置值
func Set(pattern string, value interface{}) error {
	return globalConfigurator.Set(pattern, value)
}

// Match 匹配多个规则
func Match(patterns ...string) gconfig.Matcher {
	return globalConfigurator.Match(patterns...)
}

// Close 关闭配置监听
func Close() {
	globalConfigurator.Close()
}

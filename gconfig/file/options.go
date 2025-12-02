package file

import (
	"github.com/goodluck0107/gcore/gconfig"
	"github.com/goodluck0107/gcore/getc"
)

const (
	defaultPath = "./config"
	defaultMode = gconfig.ReadOnly
)

const (
	defaultPathKey = "etc.config.file.path"
	defaultModeKey = "etc.config.file.mode"
)

type Option func(o *options)

type options struct {
	// 配置文件或配置目录路径
	path string

	// 读写模式
	// 支持read-only、write-only和read-write三种模式，默认为read-only模式
	mode gconfig.Mode
}

func defaultOptions() *options {
	return &options{
		path: getc.Get(defaultPathKey, defaultPath).String(),
		mode: gconfig.Mode(getc.Get(defaultModeKey, defaultMode).String()),
	}
}

// WithPath 设置配置文件或配置目录路径
func WithPath(path string) Option {
	return func(o *options) { o.path = path }
}

// WithMode 设置读写模式
func WithMode(mode gconfig.Mode) Option {
	return func(o *options) { o.mode = mode }
}

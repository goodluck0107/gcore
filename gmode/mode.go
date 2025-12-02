package gmode

import (
	"github.com/goodluck0107/gcore/genv"
	"github.com/goodluck0107/gcore/getc"
	"github.com/goodluck0107/gcore/gflag"
)

const (
	gcoreModeEtcName = "etc.mode"
	gcoreModeArgName = "mode"
	gcoreModeEnvName = "GCORE_MODE"
)

const (
	// DebugMode indicates gcore mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates gcore mode is release.
	ReleaseMode = "release"
	// TestMode indicates gcore mode is test.
	TestMode = "test"
)

var gcoreMode string

// 优先级： 配置文件 < 环境变量 < 运行参数 < gmode.SetMode()
func init() {
	mode := getc.Get(gcoreModeEtcName, DebugMode).String()
	mode = genv.Get(gcoreModeEnvName, mode).String()
	mode = gflag.String(gcoreModeArgName, mode)
	SetMode(mode)
}

// SetMode 设置运行模式
func SetMode(m string) {
	if m == "" {
		m = DebugMode
	}

	switch m {
	case DebugMode, TestMode, ReleaseMode:
		gcoreMode = m
	default:
		panic("gcore mode unknown: " + m + " (available mode: debug test release)")
	}
}

// GetMode 获取运行模式
func GetMode() string {
	return gcoreMode
}

// IsDebugMode 是否Debug模式
func IsDebugMode() bool {
	return gcoreMode == DebugMode
}

// IsTestMode 是否Test模式
func IsTestMode() bool {
	return gcoreMode == TestMode
}

// IsReleaseMode 是否Release模式
func IsReleaseMode() bool {
	return gcoreMode == ReleaseMode
}

package zap_test

import (
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/glog/zap"
	"testing"
)

var logger = zap.NewLogger(
	zap.WithStackLevel(glog.DebugLevel),
	zap.WithFormat(glog.JsonFormat),
)

func TestNewLogger(t *testing.T) {
	//logger.Print(log.ErrorLevel, "print")
	//logger.Info("info")
	//logger.Warn("warn")
	//logger.Error("error")
	//logger.Error("error")
	//logger.Fatal("fatal")
	logger.Panic("panic")
}

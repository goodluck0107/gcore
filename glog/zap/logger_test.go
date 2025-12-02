package zap_test

import (
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/glog/zap"
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

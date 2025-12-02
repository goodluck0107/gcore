package glog_test

import (
	"github.com/goodluck0107/gcore/glog"
	"testing"
)

func TestLog(t *testing.T) {
	logger := glog.NewLogger(glog.WithFormat(glog.JsonFormat))

	logger.Debug("welcome to gcore-framework")
	logger.Info("welcome to gcore-framework")
	logger.Warn("welcome to gcore-framework")
	logger.Error("welcome to gcore-framework")
}

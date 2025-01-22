package logrus_test

import (
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/glog/logrus"
	"testing"
)

var logger = logrus.NewLogger()

func TestNewLogger(t *testing.T) {
	//logger.Warn(`log: warn`)
	logger.Error(`log: error`)
	logger.Print(glog.ErrorLevel, `log: error`)
}

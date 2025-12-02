package logrus_test

import (
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/glog/logrus"
	"testing"
)

var logger = logrus.NewLogger()

func TestNewLogger(t *testing.T) {
	//logger.Warn(`log: warn`)
	logger.Error(`log: error`)
	logger.Print(glog.ErrorLevel, `log: error`)
}

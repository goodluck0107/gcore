package aliyun_test

import (
	"github.com/goodluck0107/gcore/glog/aliyun"
	"testing"
)

var logger = aliyun.NewLogger()

func TestNewLogger(t *testing.T) {
	defer logger.Close()

	logger.Info("info")
}

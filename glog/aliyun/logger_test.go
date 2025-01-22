package aliyun_test

import (
	"gitee.com/monobytes/gcore/glog/aliyun"
	"testing"
)

var logger = aliyun.NewLogger()

func TestNewLogger(t *testing.T) {
	defer logger.Close()

	logger.Info("info")
}

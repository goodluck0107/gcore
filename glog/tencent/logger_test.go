package tencent_test

import (
	"gitee.com/monobytes/gcore/glog/tencent"
	"testing"
)

var logger = tencent.NewLogger()

func TestNewLogger(t *testing.T) {
	defer logger.Close()

	logger.Error("error")
}

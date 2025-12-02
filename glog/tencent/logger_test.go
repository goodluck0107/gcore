package tencent_test

import (
	"github.com/goodluck0107/gcore/glog/tencent"
	"testing"
)

var logger = tencent.NewLogger()

func TestNewLogger(t *testing.T) {
	defer logger.Close()

	logger.Error("error")
}

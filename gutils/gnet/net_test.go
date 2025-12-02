package gnet_test

import (
	"github.com/goodluck0107/gcore/gutils/gnet"
	"testing"
)

func TestIP2Long(t *testing.T) {
	str1 := "218.108.212.34"
	ip := gnet.IP2Long(str1)
	str2 := gnet.Long2IP(ip)

	t.Logf("str format: %s", str1)
	t.Logf("long format: %d", ip)
	t.Logf("str format: %s", str2)
}

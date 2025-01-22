package getc_test

import (
	"gitee.com/monobytes/gcore/getc"
	"testing"
)

func Test_Get(t *testing.T) {
	v := getc.Get("c.redis.addrs.1A", "192.168.0.1:3308").String()
	t.Log(v)
}

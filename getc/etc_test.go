package getc_test

import (
	"github.com/goodluck0107/gcore/getc"
	"testing"
)

func Test_Get(t *testing.T) {
	v := getc.Get("c.redis.addrs.1A", "192.168.0.1:3308").String()
	t.Log(v)
}

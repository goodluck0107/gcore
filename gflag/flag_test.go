package gflag_test

import (
	"gitee.com/monobytes/gcore/gflag"
	"testing"
)

func TestString(t *testing.T) {
	t.Log(gflag.Bool("test.v"))
	t.Log(gflag.String("config", "./config"))
}

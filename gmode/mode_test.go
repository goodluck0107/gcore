package gmode_test

import (
	"flag"
	"testing"

	"gitee.com/monobytes/gcore/gmode"
)

func TestGetMode(t *testing.T) {
	flag.Parse()

	t.Log(gmode.GetMode())
}

package gflag_test

import (
	"github.com/goodluck0107/gcore/gflag"
	"testing"
)

func TestString(t *testing.T) {
	t.Log(gflag.Bool("test.v"))
	t.Log(gflag.String("config", "./config"))
}

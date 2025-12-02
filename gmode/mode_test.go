package gmode_test

import (
	"flag"
	"testing"

	"github.com/goodluck0107/gcore/gmode"
)

func TestGetMode(t *testing.T) {
	flag.Parse()

	t.Log(gmode.GetMode())
}

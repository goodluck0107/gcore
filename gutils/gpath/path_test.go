package gpath_test

import (
	"github.com/goodluck0107/gcore/gutils/gpath"
	"testing"
)

func TestSplit(t *testing.T) {
	path := "/etc/my.ini"

	t.Log(gpath.Split(path))
}

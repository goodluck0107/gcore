package gpath_test

import (
	"gitee.com/monobytes/gcore/gutils/gpath"
	"testing"
)

func TestSplit(t *testing.T) {
	path := "/etc/my.ini"

	t.Log(gpath.Split(path))
}

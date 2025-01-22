package gfile_test

import (
	"gitee.com/monobytes/gcore/gutils/gfile"
	"testing"
)

func TestStat(t *testing.T) {
	fi, err := gfile.Stat("a.txt")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(fi.CreateTime())
}

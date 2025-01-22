package gfile_test

import (
	"gitee.com/monobytes/gcore/gutils/gfile"
	"testing"
)

func TestWriteFile(t *testing.T) {
	err := gfile.WriteFile("./run/test.txt", []byte("hello world"))
	if err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}

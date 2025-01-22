package gconv_test

import (
	"fmt"
	"gitee.com/monobytes/gcore/gutils/gconv"
	"testing"
)

func TestBytesToString(t *testing.T) {
	b := []byte("abc")

	s := gconv.BytesToString(b)

	fmt.Printf("%p\n", &b)
	fmt.Printf("%p\n", &s)
}

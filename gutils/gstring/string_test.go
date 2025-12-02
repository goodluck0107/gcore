package gstring_test

import (
	"github.com/goodluck0107/gcore/gutils/gstring"
	"testing"
)

func Test_PaddingPrefix(t *testing.T) {
	t.Log(gstring.PaddingPrefix("1", "0", 3))
	t.Log(gstring.PaddingPrefix("001", "0", 3))
	t.Log(gstring.PaddingPrefix("0001", "0", 3))
	t.Log(gstring.PaddingPrefix("1", "00", 3))
}

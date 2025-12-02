package greflect_test

import (
	"github.com/goodluck0107/gcore/gutils/greflect"
	"testing"
)

func TestIsNil(t *testing.T) {
	var b1 bool
	var b2 *bool
	t.Log(greflect.IsNil(b1))
	t.Log(greflect.IsNil(&b1))
	t.Log(greflect.IsNil(b2))
	t.Log(greflect.IsNil(&b2))
}

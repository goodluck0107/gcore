package gerrors_test

import (
	"fmt"
	"github.com/goodluck0107/gcore/gcodes"
	"github.com/goodluck0107/gcore/gerrors"
	"testing"
)

func TestNew(t *testing.T) {
	innerErr := gerrors.NewError(
		"db error",
		gcodes.NewCode(2, "core error"),
		gerrors.New("std not found"),
	)

	err := gerrors.NewError(
		//"not found",
		gcodes.NewCode(1, "not found"),
		innerErr,
	)

	t.Log(err)
	t.Log(err.Code())
	t.Log(err.Next())
	t.Log(err.Cause())
	fmt.Println(fmt.Sprintf("%+v", err))
}

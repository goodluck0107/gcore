package gcodes_test

import (
	"errors"
	"github.com/goodluck0107/gcore/gcodes"
	"testing"
)

func TestConvert(t *testing.T) {
	code, _ := gcodes.Convert(errors.New("rpc error: code = Unknown desc = code error: code = 10 desc = account exists"))

	t.Log(code)
}

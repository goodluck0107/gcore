package gcodes_test

import (
	"errors"
	"gitee.com/monobytes/gcore/gcodes"
	"testing"
)

func TestConvert(t *testing.T) {
	code, _ := gcodes.Convert(errors.New("rpc error: code = Unknown desc = code error: code = 10 desc = account exists"))

	t.Log(code)
}

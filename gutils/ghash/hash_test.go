package ghash_test

import (
	"gitee.com/monobytes/gcore/gutils/ghash"
	"testing"
)

func TestSHA256(t *testing.T) {
	t.Log(ghash.SHA256("abc"))
}

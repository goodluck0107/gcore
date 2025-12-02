package ghash_test

import (
	"github.com/goodluck0107/gcore/gutils/ghash"
	"testing"
)

func TestSHA256(t *testing.T) {
	t.Log(ghash.SHA256("abc"))
}

package gos_test

import (
	"github.com/goodluck0107/gcore/gutils/gos"
	"testing"
)

func TestCreate(t *testing.T) {
	_, err := gos.Create("./mpprof/server/cpu_profile")
	if err != nil {
		t.Fatal(err)
	}

}

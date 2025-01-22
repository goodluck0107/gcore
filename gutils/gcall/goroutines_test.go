package gcall_test

import (
	"context"
	"gitee.com/monobytes/gcore/gutils/gcall"
	"testing"
	"time"
)

func TestGoroutines_Run(t *testing.T) {
	g := gcall.NewGoroutines()

	g.Add(func() {
		time.Sleep(2 * time.Second)
	}).Add(func() {
		time.Sleep(3 * time.Second)
	}).Add(func() {
		time.Sleep(4 * time.Second)
	}).Add(func() {
		time.Sleep(10 * time.Second)
	}).Run(context.Background(), 5*time.Second)
}

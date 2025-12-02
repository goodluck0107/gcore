package redis_test

import (
	"context"
	"fmt"
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/glocate/redis"
	"github.com/goodluck0107/gcore/gutils/guuid"
	"testing"
	"time"
)

var locator = redis.NewLocator(
	redis.WithAddrs("127.0.0.1:6379"),
)

func TestLocator_BindGate(t *testing.T) {
	for i := 1; i <= 6; i++ {
		gid := guuid.UUID()

		err := locator.BindGate(context.Background(), int64(i), gid)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLocator_BindNode(t *testing.T) {
	for i := 1; i <= 6; i++ {
		nid := guuid.UUID()

		name := fmt.Sprintf("node-%d", i)

		err := locator.BindNode(context.Background(), int64(i), name, nid)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLocator_UnbindGate(t *testing.T) {
	for i := 1; i <= 6; i++ {
		gid := guuid.UUID()
		ctx := context.Background()
		uid := int64(i)

		err := locator.BindGate(ctx, uid, gid)
		if err != nil {
			t.Fatal(err)
		}

		err = locator.UnbindGate(ctx, uid, gid)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLocator_UnbindNode(t *testing.T) {
	for i := 1; i <= 6; i++ {
		nid := guuid.UUID()

		ctx := context.Background()
		uid := int64(i)
		name := fmt.Sprintf("node-%d", i)

		err := locator.BindNode(ctx, uid, name, nid)
		if err != nil {
			t.Fatal(err)
		}

		err = locator.UnbindNode(ctx, uid, name, nid)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestLocator_Watch(t *testing.T) {
	watcher1, err := locator.Watch(context.Background(), gcluster.Gate.String(), gcluster.Node.String())
	if err != nil {
		t.Fatal(err)
	}

	watcher2, err := locator.Watch(context.Background(), gcluster.Gate.String())
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			events, err := watcher1.Next()
			if err != nil {
				t.Errorf("goroutine 1: %v", err)
				return
			}

			fmt.Println("goroutine 1: new event entity")

			for _, event := range events {
				t.Logf("goroutine 1: %+v", event)
			}
		}
	}()

	go func() {
		for {
			events, err := watcher2.Next()
			if err != nil {
				t.Errorf("goroutine 2: %v", err)
				return
			}

			fmt.Println("goroutine 2: new event entity")

			for _, event := range events {
				t.Logf("goroutine 2: %+v", event)
			}
		}
	}()

	time.Sleep(60 * time.Second)
}

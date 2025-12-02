package etcd_test

import (
	"context"
	"github.com/goodluck0107/gcore/gconfig"
	"github.com/goodluck0107/gcore/gconfig/etcd"
	"testing"
	"time"
)

func init() {
	source := etcd.NewSource(etcd.WithMode(gconfig.ReadWrite))
	gconfig.SetConfigurator(gconfig.NewConfigurator(gconfig.WithSources(source)))
}

func TestWatch(t *testing.T) {
	ticker1 := time.NewTicker(2 * time.Second)
	ticker2 := time.After(time.Minute)

	for {
		select {
		case <-ticker1.C:
			t.Log(gconfig.Get("config.timezone").String())
			t.Log(gconfig.Get("config.pid").String())
		case <-ticker2:
			gconfig.Close()
			return
		}
	}
}

func TestLoad(t *testing.T) {
	ctx := context.Background()
	file := "config.json"
	c, err := gconfig.Load(ctx, etcd.Name, file)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(c[0].Name)
	t.Log(c[0].Path)
	t.Log(c[0].Format)
	t.Log(c[0].Content)
}

func TestStore(t *testing.T) {
	ctx := context.Background()
	file := "config.json"
	content1 := map[string]interface{}{
		"timezone": "Local",
	}

	content2 := map[string]interface{}{
		"timezone": "UTC",
		"pid":      "./run/gate.pid",
	}

	err := gconfig.Store(ctx, etcd.Name, file, content1, true)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	err = gconfig.Store(ctx, etcd.Name, file, content2)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gconfig.Get("config").Value()
	}
}

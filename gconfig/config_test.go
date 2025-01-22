package gconfig_test

import (
	"context"
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/gconfig/file"
	"testing"
	"time"
)

func init() {
	source := file.NewSource(file.WithMode(gconfig.ReadWrite))
	gconfig.SetConfigurator(gconfig.NewConfigurator(gconfig.WithSources(source)))
}

func TestWatch(t *testing.T) {
	ticker1 := time.NewTicker(2 * time.Second)
	ticker2 := time.After(5 * time.Second)

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

func TestStore(t *testing.T) {
	ctx := context.Background()
	filename := "config.json"
	content1 := map[string]interface{}{
		"timezone": "Local",
	}

	content2 := map[string]interface{}{
		"timezone": "UTC",
		"pid":      "./run/gate.pid",
	}

	err := gconfig.Store(ctx, file.Name, filename, content1, true)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	err = gconfig.Store(ctx, file.Name, filename, content2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	ctx := context.Background()
	filename := "config.json"
	c, err := gconfig.Load(ctx, file.Name, filename)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(c[0].Name)
	t.Log(c[0].Path)
	t.Log(c[0].Format)
	t.Log(c[0].Content)
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gconfig.Get("config").Value()
	}
}

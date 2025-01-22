package nacos_test

import (
	"context"
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/gconfig/nacos"
	"testing"
	"time"
)

func init() {
	source := nacos.NewSource()
	gconfig.SetConfigurator(gconfig.NewConfigurator(gconfig.WithSources(source)))
}

func TestWatch(t *testing.T) {
	ticker1 := time.NewTicker(2 * time.Second)
	ticker2 := time.After(20 * time.Minute)

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
	c, err := gconfig.Load(ctx, nacos.Name, file)
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
	file := "configs.json"
	content1 := map[string]interface{}{
		"timezone": "Local",
	}

	content2 := map[string]interface{}{
		"timezone": "UTC",
		"pid":      "./run/gate.pid",
	}

	err := gconfig.Store(ctx, nacos.Name, file, content1, true)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	err = gconfig.Store(ctx, nacos.Name, file, content2)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gconfig.Get("config").Value()
	}
}

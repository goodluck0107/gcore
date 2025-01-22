package consul

import (
	"context"
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gutils/gcall"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"path/filepath"
	"strings"
)

type watcher struct {
	ctx     context.Context
	cancel  context.CancelFunc
	source  *Source
	plan    *watch.Plan
	chWatch chan []*gconfig.Configuration
}

func newWatcher(ctx context.Context, s *Source) (gconfig.Watcher, error) {
	w := &watcher{}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.source = s
	w.chWatch = make(chan []*gconfig.Configuration, 2)

	if err := w.init(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *watcher) init() (err error) {
	var prefix string
	if w.source.opts.path != "" {
		prefix = w.source.opts.path + "/"
	}

	w.plan, err = watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": prefix,
	})
	if err != nil {
		return
	}

	w.plan.Handler = w.planHandler

	gcall.Go(func() {
		if err = w.plan.RunWithClientAndHclog(w.source.opts.client, nil); err != nil {
			glog.Fatalf("create watcher failed: %v", err)
		}
	})

	return
}

func (w *watcher) planHandler(idx uint64, raw interface{}) {
	if raw == nil {
		return // ignore
	}

	kvs, ok := raw.(api.KVPairs)
	if !ok {
		return
	}

	configs := make([]*gconfig.Configuration, 0, len(kvs))
	for _, kv := range kvs {
		configs = append(configs, w.parseKV(kv.Key, kv.Value))
	}

	w.chWatch <- configs
}

// 解析KV
func (w *watcher) parseKV(key string, value []byte) *gconfig.Configuration {
	fullPath := key
	path := strings.TrimPrefix(fullPath, w.source.opts.path)
	file := filepath.Base(fullPath)
	ext := filepath.Ext(file)

	return &gconfig.Configuration{
		Path:     path,
		File:     file,
		Name:     strings.TrimSuffix(file, ext),
		Format:   strings.TrimPrefix(ext, "."),
		Content:  value,
		FullPath: fullPath,
	}
}

// Next 返回配置列表
func (w *watcher) Next() ([]*gconfig.Configuration, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case configs, ok := <-w.chWatch:
		if !ok {
			if err := w.ctx.Err(); err != nil {
				return nil, err
			}
		}

		return configs, nil
	}
}

// Stop 停止监听
func (w *watcher) Stop() error {
	w.cancel()
	w.plan.Stop()

	return nil
}

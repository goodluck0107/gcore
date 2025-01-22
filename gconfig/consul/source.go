package consul

import (
	"context"
	"gitee.com/monobytes/gcore/gconfig"
	"github.com/hashicorp/consul/api"
	"path/filepath"
	"strings"
)

const Name = "consul"

type Source struct {
	err  error
	opts *options
}

func NewSource(opts ...Option) gconfig.Source {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	s := &Source{}
	s.opts = o
	s.opts.path = strings.TrimSuffix(strings.TrimPrefix(s.opts.path, "/"), "/")

	if o.client == nil {
		c := api.DefaultConfig()
		if o.addr != "" {
			c.Address = o.addr
		}

		s.opts.client, s.err = api.NewClient(c)
	}

	return s
}

// Name 配置源名称
func (s *Source) Name() string {
	return Name
}

// Load 加载配置项
func (s *Source) Load(ctx context.Context, file ...string) ([]*gconfig.Configuration, error) {
	if s.err != nil {
		return nil, s.err
	}

	var prefix string

	if s.opts.path != "" {
		if len(file) > 0 && file[0] != "" {
			prefix = s.opts.path + "/" + strings.TrimPrefix(file[0], "/")
		} else {
			prefix = s.opts.path + "/"
		}
	}

	kvs, _, err := s.opts.client.KV().List(prefix, nil)
	if err != nil {
		return nil, err
	}

	if len(kvs) == 0 {
		return nil, nil
	}

	configs := make([]*gconfig.Configuration, 0, len(kvs))
	for _, kv := range kvs {
		fullPath := kv.Key
		path := strings.TrimPrefix(fullPath, s.opts.path)
		file := filepath.Base(fullPath)
		ext := filepath.Ext(file)
		configs = append(configs, &gconfig.Configuration{
			Path:     path,
			File:     file,
			Name:     strings.TrimSuffix(file, ext),
			Format:   strings.TrimPrefix(ext, "."),
			Content:  kv.Value,
			FullPath: fullPath,
		})
	}

	return configs, nil
}

// Store 保存配置项
func (s *Source) Store(ctx context.Context, file string, content []byte) error {
	if s.err != nil {
		return s.err
	}

	var key string

	if s.opts.path != "" {
		key = s.opts.path + "/" + strings.TrimPrefix(file, "/")
	} else {
		key = strings.TrimPrefix(file, "/")
	}

	_, err := s.opts.client.KV().Put(&api.KVPair{
		Key:   key,
		Value: content,
	}, nil)

	return err
}

// Watch 监听配置项
func (s *Source) Watch(ctx context.Context) (gconfig.Watcher, error) {
	if s.err != nil {
		return nil, s.err
	}

	return newWatcher(ctx, s)
}

// Close 关闭配置源
func (s *Source) Close() error {
	return nil
}

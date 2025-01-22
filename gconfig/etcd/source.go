package etcd

import (
	"context"
	"fmt"
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/gerrors"
	"gitee.com/monobytes/gcore/gutils/gconv"
	"go.etcd.io/etcd/client/v3"
	"path/filepath"
	"strings"
)

const Name = "etcd"

type Source struct {
	err     error
	opts    *options
	builtin bool
}

func NewSource(opts ...Option) gconfig.Source {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	s := &Source{}
	s.opts = o
	s.opts.path = fmt.Sprintf("/%s", strings.TrimSuffix(strings.TrimPrefix(s.opts.path, "/"), "/"))

	if o.client == nil {
		s.builtin = true
		o.client, s.err = clientv3.New(clientv3.Config{
			Endpoints:   o.addrs,
			DialTimeout: o.dialTimeout,
		})
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

	var (
		key  = s.opts.path
		opts []clientv3.OpOption
	)

	if len(file) > 0 && file[0] != "" {
		key += "/" + strings.TrimPrefix(file[0], "/")
	} else {
		opts = append(opts, clientv3.WithPrefix())
	}

	res, err := s.opts.client.Get(ctx, key, opts...)
	if err != nil {
		return nil, err
	}

	configs := make([]*gconfig.Configuration, 0, len(res.Kvs))
	for _, kv := range res.Kvs {
		fullPath := string(kv.Key)
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

	if s.opts.mode != gconfig.WriteOnly && s.opts.mode != gconfig.ReadWrite {
		return gerrors.ErrNoOperationPermission
	}

	key := s.opts.path + "/" + strings.TrimPrefix(file, "/")
	_, err := s.opts.client.Put(ctx, key, gconv.String(content))
	return err
}

// Watch 监听配置项
func (s *Source) Watch(ctx context.Context) (gconfig.Watcher, error) {
	if s.err != nil {
		return nil, s.err
	}

	return newWatcher(ctx, s)
}

// Close 关闭资源
func (s *Source) Close() error {
	if s.err != nil {
		return s.err
	}

	if s.builtin {
		return s.opts.client.Close()
	}

	return nil
}

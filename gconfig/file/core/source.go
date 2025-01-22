package core

import (
	"context"
	"gitee.com/monobytes/gcore/gconfig"
	"gitee.com/monobytes/gcore/gerrors"
	"gitee.com/monobytes/gcore/gutils/gfile"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const Name = "file"

type Source struct {
	path string
	mode gconfig.Mode
}

var _ gconfig.Source = &Source{}

func NewSource(path string, mode gconfig.Mode) *Source {
	return &Source{path: strings.TrimSuffix(path, "/"), mode: mode}
}

// Name 配置源名称
func (s *Source) Name() string {
	return Name
}

// Load 加载配置
func (s *Source) Load(ctx context.Context, file ...string) ([]*gconfig.Configuration, error) {
	path := s.path

	if len(file) > 0 && file[0] != "" {
		info, err := os.Stat(s.path)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			return nil, gerrors.New("the specified file cannot be loaded at the file path")
		}

		path = filepath.Join(s.path, file[0])
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return s.loadDir(path)
	}

	c, err := s.loadFile(path)
	if err != nil {
		return nil, err
	}

	return []*gconfig.Configuration{c}, nil
}

// Store 保存配置项
func (s *Source) Store(ctx context.Context, file string, content []byte) error {
	if s.mode != gconfig.WriteOnly && s.mode != gconfig.ReadWrite {
		return gerrors.ErrNoOperationPermission
	}

	info, err := os.Stat(s.path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return gerrors.New("the specified file cannot be modified under the file path")
	}

	return gfile.WriteFile(filepath.Join(s.path, file), content)
}

// Watch 监听配置变化
func (s *Source) Watch(ctx context.Context) (gconfig.Watcher, error) {
	return newWatcher(ctx, s)
}

func (s *Source) Close() error {
	return nil
}

// 加载文件配置
func (s *Source) loadFile(path string) (*gconfig.Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(info.Name())
	path1, _ := filepath.Abs(path)
	path2, _ := filepath.Abs(s.path)
	path = strings.TrimPrefix(path1, path2)
	fullPath := s.path + path

	return &gconfig.Configuration{
		Path:     path,
		File:     info.Name(),
		Name:     strings.TrimSuffix(info.Name(), ext),
		Format:   strings.TrimPrefix(ext, "."),
		Content:  content,
		FullPath: fullPath,
	}, nil
}

// 加载目录配置
func (s *Source) loadDir(path string) (cs []*gconfig.Configuration, err error) {
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || strings.HasSuffix(d.Name(), ".") {
			return nil
		}

		c, err := s.loadFile(path)
		if err != nil {
			return err
		}
		cs = append(cs, c)

		return nil
	})

	return
}

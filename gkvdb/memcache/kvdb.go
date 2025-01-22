// Package memcache 还在开发中
package memcache

import (
	"context"
	"gitee.com/monobytes/gcore/gerrors"
	"gitee.com/monobytes/gcore/gkvdb"
	"gitee.com/monobytes/gcore/gutils/gconv"
	"gitee.com/monobytes/gcore/gutils/grand"
	"github.com/bradfitz/gomemcache/memcache"
	"golang.org/x/sync/singleflight"
	"reflect"
	"time"
)

type KvDB struct {
	opts    *options
	builtin bool
	sfg     singleflight.Group
}

func NewKvDB(opts ...Option) *KvDB {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	c := &KvDB{}
	c.opts = o

	if o.client == nil {
		c.builtin = true
		o.client = memcache.New(o.addrs...)
	}

	return c
}

// Has 检测缓存是否存在
func (c *KvDB) Has(ctx context.Context, key string) (bool, error) {
	key = c.AddPrefix(key)

	val, err, _ := c.sfg.Do(key, func() (interface{}, error) {
		item, err := c.opts.client.Get(key)
		if err != nil {
			return nil, err
		}

		return gconv.String(item.Value), nil
	})
	if err != nil {
		if gerrors.Is(err, memcache.ErrCacheMiss) {
			return false, nil
		}
		return false, err
	}

	if val.(string) == c.opts.nilValue {
		return false, nil
	}

	return true, nil
}

// Get 获取缓存值
func (c *KvDB) Get(ctx context.Context, key string, def ...interface{}) gkvdb.Result {
	key = c.AddPrefix(key)

	val, err, _ := c.sfg.Do(key, func() (interface{}, error) {
		item, err := c.opts.client.Get(key)
		if err != nil {
			return nil, err
		}

		return gconv.String(item.Value), nil
	})
	if err != nil && !gerrors.Is(err, memcache.ErrCacheMiss) {
		return gkvdb.NewResult(nil, err)
	}

	if gerrors.Is(err, memcache.ErrCacheMiss) || val.(string) == c.opts.nilValue {
		if len(def) > 0 {
			return gkvdb.NewResult(def[0])
		} else {
			return gkvdb.NewResult(nil, gerrors.ErrNil)
		}
	}

	return gkvdb.NewResult(val)
}

// Set 设置缓存值
func (c *KvDB) Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error {
	if len(expiration) > 0 && expiration[0] > 0 {
		return c.opts.client.Set(&memcache.Item{
			Key:        c.AddPrefix(key),
			Value:      gconv.Bytes(value),
			Expiration: int32(expiration[0] / time.Second),
		})
	} else {
		return c.opts.client.Set(&memcache.Item{
			Key:   c.AddPrefix(key),
			Value: gconv.Bytes(value),
		})
	}
}

// GetSet 获取设置缓存值
func (c *KvDB) GetSet(ctx context.Context, key string, fn gkvdb.SetValueFunc) gkvdb.Result {
	key = c.AddPrefix(key)

	val, err, _ := c.sfg.Do(key, func() (interface{}, error) {
		item, err := c.opts.client.Get(key)
		if err != nil {
			return nil, err
		}

		return gconv.String(item.Value), nil
	})
	if err != nil && !gerrors.Is(err, memcache.ErrCacheMiss) {
		return gkvdb.NewResult(nil, err)
	}

	if err == nil {
		if val == c.opts.nilValue {
			return gkvdb.NewResult(nil, gerrors.ErrNil)
		} else {
			return gkvdb.NewResult(val)
		}
	}

	rst, _, _ := c.sfg.Do(key+":set", func() (interface{}, error) {
		val, err = fn()
		if err != nil {
			return gkvdb.NewResult(nil, err), nil
		}

		if val == nil || reflect.ValueOf(val).IsNil() {
			if err = c.opts.client.Set(&memcache.Item{
				Key:        key,
				Value:      gconv.Bytes(c.opts.nilValue),
				Expiration: int32(c.opts.nilExpiration / time.Second),
			}); err != nil {
				return gkvdb.NewResult(nil, err), nil
			}
			return gkvdb.NewResult(nil, gerrors.ErrNil), nil
		}
		expiration := time.Duration(grand.Int64(int64(c.opts.minExpiration), int64(c.opts.maxExpiration)))
		if err = c.opts.client.Set(&memcache.Item{
			Key:        key,
			Value:      gconv.Bytes(val),
			Expiration: int32(expiration / time.Second),
		}); err != nil {
			return gkvdb.NewResult(nil, err), nil
		}

		return gkvdb.NewResult(val, nil), nil
	})

	return rst.(gkvdb.Result)
}

// Delete 删除缓存
func (c *KvDB) Delete(ctx context.Context, key string) (bool, error) {
	err := c.opts.client.Delete(c.AddPrefix(key))
	if err != nil {
		if gerrors.Is(err, memcache.ErrCacheMiss) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// IncrInt 整数自增
func (c *KvDB) IncrInt(ctx context.Context, key string, value int64) (int64, error) {
	if value < 0 {
		return c.DecrInt(ctx, key, 0-value)
	}

	key = c.AddPrefix(key)

	newValue, err := c.opts.client.Increment(key, uint64(value))
	if err != nil {
		if gerrors.Is(err, memcache.ErrCacheMiss) {
			err = c.opts.client.Add(&memcache.Item{
				Key:   key,
				Value: gconv.Bytes(gconv.String(value)),
			})
			if err != nil {
				if gerrors.Is(err, memcache.ErrNotStored) {
					newValue, err = c.opts.client.Increment(key, uint64(value))
					if err != nil {
						return 0, err
					}

					return int64(newValue), nil
				}
				return 0, err
			}

			return value, nil
		} else {
			return 0, err
		}
	}

	return int64(newValue), nil
}

// IncrFloat 浮点数自增
func (c *KvDB) IncrFloat(ctx context.Context, key string, value float64) (float64, error) {
	return 0, nil
}

// DecrInt 整数自减
func (c *KvDB) DecrInt(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

// DecrFloat 浮点数自减
func (c *KvDB) DecrFloat(ctx context.Context, key string, value float64) (float64, error) {
	return 0, nil
}

// AddPrefix 添加Key前缀
func (c *KvDB) AddPrefix(key string) string {
	if c.opts.prefix == "" {
		return key
	} else {
		return c.opts.prefix + ":" + key
	}
}

// Client 获取客户端
func (c *KvDB) Client() interface{} {
	return c.opts.client
}

// Close 关闭客户端
func (c *KvDB) Close() error {
	if !c.builtin {
		return nil
	}

	return c.opts.client.Close()
}

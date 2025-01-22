package nacos

import (
	"context"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gregistry"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sync"
	"sync/atomic"
)

type watcher struct {
	idx        int64
	state      int32
	ctx        context.Context
	cancel     context.CancelFunc
	watcherMgr *watcherMgr
	chWatch    chan []*gregistry.ServiceInstance
}

func newWatcher(wm *watcherMgr, idx int64) *watcher {
	w := &watcher{}
	w.ctx, w.cancel = context.WithCancel(wm.ctx)
	w.idx = idx
	w.watcherMgr = wm
	w.chWatch = make(chan []*gregistry.ServiceInstance, 16)

	return w
}

func (w *watcher) notify(services []*gregistry.ServiceInstance) {
	if atomic.LoadInt32(&w.state) == 0 {
		return
	}

	w.chWatch <- services
}

// Next 返回服务实例列表
func (w *watcher) Next() ([]*gregistry.ServiceInstance, error) {
	if atomic.LoadInt32(&w.state) == 0 {
		atomic.StoreInt32(&w.state, 1)
		return w.watcherMgr.services(), nil
	}

	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case services, ok := <-w.chWatch:
		if !ok {
			if err := w.ctx.Err(); err != nil {
				return nil, err
			}
		}

		return services, nil
	}
}

// Stop 停止监听
func (w *watcher) Stop() error {
	w.cancel()
	close(w.chWatch)
	return w.watcherMgr.recycle(w.idx)
}

type watcherMgr struct {
	err              error
	ctx              context.Context
	cancel           context.CancelFunc
	registry         *Registry
	serviceName      string
	serviceInstances *atomic.Value
	serviceWaitIndex uint64
	idx              int64
	rw               sync.RWMutex
	watchers         map[int64]*watcher
}

func newWatcherMgr(registry *Registry, ctx context.Context, serviceName string) (*watcherMgr, error) {
	services, err := registry.services(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	wm := &watcherMgr{}
	wm.ctx, wm.cancel = context.WithCancel(registry.ctx)
	wm.registry = registry
	wm.serviceName = serviceName
	wm.serviceInstances = &atomic.Value{}
	wm.serviceInstances.Store(services)
	wm.watchers = make(map[int64]*watcher)

	if err = wm.subscribe(); err != nil {
		return nil, err
	}

	return wm, nil
}

func (wm *watcherMgr) subscribe() error {
	return wm.registry.opts.client.Subscribe(&vo.SubscribeParam{
		ServiceName: wm.serviceName,
		Clusters:    []string{wm.registry.opts.clusterName},
		GroupName:   wm.registry.opts.groupName,
		SubscribeCallback: func(instances []model.Instance, err error) {
			if err != nil {
				glog.Warnf("%s subscribe callback failed: %v", wm.serviceName, err)
				return
			}

			services, err := parseInstances(instances)
			if err != nil {
				glog.Warnf("%s instances parse failed: %v", wm.serviceName, err)
				return
			}

			wm.serviceInstances.Store(services)
			wm.broadcast(services)
		},
	})
}

func (wm *watcherMgr) fork() gregistry.Watcher {
	wm.rw.Lock()
	defer wm.rw.Unlock()

	w := newWatcher(wm, atomic.AddInt64(&wm.idx, 1))
	wm.watchers[w.idx] = w

	return w
}

func (wm *watcherMgr) recycle(idx int64) error {
	wm.rw.Lock()
	defer wm.rw.Unlock()

	delete(wm.watchers, idx)

	if len(wm.watchers) == 0 {
		wm.cancel()
		wm.registry.watchers.Delete(wm.serviceName)
	}

	return nil
}

func (wm *watcherMgr) broadcast(services []*gregistry.ServiceInstance) {
	wm.rw.RLock()
	defer wm.rw.RUnlock()

	for _, w := range wm.watchers {
		w.notify(services)
	}
}

func (wm *watcherMgr) services() []*gregistry.ServiceInstance {
	return wm.serviceInstances.Load().([]*gregistry.ServiceInstance)
}

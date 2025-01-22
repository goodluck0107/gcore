package dispatcher

import (
	"gitee.com/monobytes/gcore/gerrors"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/gwrap/endpoint"
	"sync"
)

type BalanceStrategy string

const (
	Random           BalanceStrategy = "random" // 随机
	RoundRobin       BalanceStrategy = "rr"     // 轮询
	WeightRoundRobin BalanceStrategy = "wrr"    // 加权轮询
)

type Dispatcher struct {
	strategy  BalanceStrategy
	rw        sync.RWMutex
	routes    map[int32]*Route
	events    map[int]*Event
	endpoints map[string]*endpoint.Endpoint
	instances map[string]*gregistry.ServiceInstance
}

func NewDispatcher(strategy BalanceStrategy) *Dispatcher {
	return &Dispatcher{strategy: strategy}
}

// FindEndpoint 查找服务端口
func (d *Dispatcher) FindEndpoint(insID string) (*endpoint.Endpoint, error) {
	d.rw.RLock()
	defer d.rw.RUnlock()

	ep, ok := d.endpoints[insID]
	if !ok {
		return nil, gerrors.ErrNotFoundEndpoint
	}

	return ep, nil
}

// IterateEndpoint 迭代服务端口
func (d *Dispatcher) IterateEndpoint(fn func(insID string, ep *endpoint.Endpoint) bool) {
	d.rw.RLock()
	defer d.rw.RUnlock()

	for insID, ep := range d.endpoints {
		if fn(insID, ep) == false {
			break
		}
	}
}

// FindRoute 查找节点路由
func (d *Dispatcher) FindRoute(route int32) (*Route, error) {
	d.rw.RLock()
	defer d.rw.RUnlock()

	r, ok := d.routes[route]
	if !ok {
		return nil, gerrors.ErrNotFoundRoute
	}

	return r, nil
}

// FindEvent 查找节点事件
func (d *Dispatcher) FindEvent(event int) (*Event, error) {
	d.rw.RLock()
	defer d.rw.RUnlock()

	e, ok := d.events[event]
	if !ok {
		return nil, gerrors.ErrNotFoundEvent
	}

	return e, nil
}

// ReplaceServices 替换服务
func (d *Dispatcher) ReplaceServices(services ...*gregistry.ServiceInstance) {
	routes := make(map[int32]*Route, len(services))
	events := make(map[int]*Event, len(services))
	endpoints := make(map[string]*endpoint.Endpoint)
	instances := make(map[string]*gregistry.ServiceInstance, len(services))

	for _, service := range services {
		ep, err := endpoint.ParseEndpoint(service.Endpoint)
		if err != nil {
			glog.Errorf("service endpoint parse failed, insID: %s kind: %s name: %s alias: %s endpoint: %s err: %v",
				service.ID, service.Kind, service.Name, service.Alias, service.Endpoint, err)
			continue
		}

		endpoints[service.ID] = ep
		instances[service.ID] = service

		for _, item := range service.Routes {
			route, ok := routes[item.ID]
			if !ok {
				route = newRoute(d, item.ID, service.Alias, item.Stateful, item.Internal)
				routes[item.ID] = route
			}
			route.addEndpoint(service.ID, service.State, ep)
		}

		for _, evt := range service.Events {
			event, ok := events[evt]
			if !ok {
				event = newEvent(d, evt)
				events[evt] = event
			}
			event.addEndpoint(service.ID, service.State, ep)
		}
	}

	d.rw.Lock()
	d.routes = routes
	d.events = events
	d.endpoints = endpoints
	d.instances = instances

	if d.strategy == WeightRoundRobin {
		for _, route := range routes {
			route.initWRRQueue()
		}
		for _, event := range events {
			event.initWRRQueue()
		}
	}
	d.rw.Unlock()
}

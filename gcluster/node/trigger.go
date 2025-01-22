package node

import (
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/glog"
)

type EventHandler func(ctx Context)

type Trigger struct {
	node    *Node
	events  map[gcluster.Event]EventHandler
	evtChan chan *event
}

func newTrigger(node *Node) *Trigger {
	return &Trigger{
		node:    node,
		events:  make(map[gcluster.Event]EventHandler, 3),
		evtChan: make(chan *event, 4096),
	}
}

func (e *Trigger) trigger(kind gcluster.Event, gid string, cid, uid int64) {
	evt := e.node.evtPool.Get().(*event)
	evt.event = kind
	evt.gid = gid
	evt.cid = cid
	evt.uid = uid
	e.evtChan <- evt
}

func (e *Trigger) receive() <-chan *event {
	return e.evtChan
}

func (e *Trigger) close() {
	close(e.evtChan)
}

// 处理事件消息
func (e *Trigger) handle(evt *event) {
	version := evt.incrVersion()

	defer evt.compareVersionRecycle(version)

	handler, ok := e.events[evt.event]
	if !ok {
		return
	}

	handler(evt)

	evt.compareVersionExecDefer(version)
}

// AddEventHandler 添加事件处理器
func (e *Trigger) AddEventHandler(event gcluster.Event, handler EventHandler) {
	if e.node.getState() != gcluster.Shut {
		glog.Warnf("the node server is working, can't add Event handler")
		return
	}

	e.events[event] = handler
}

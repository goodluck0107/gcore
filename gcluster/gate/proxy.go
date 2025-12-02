package gate

import (
	"context"
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gerrors"
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/gmode"
	"github.com/goodluck0107/gcore/gpacket"
	"github.com/goodluck0107/gcore/internal/link"
)

type proxy struct {
	gate       *Gate            // 网关服
	nodeLinker *link.NodeLinker // 节点链接器
}

func newProxy(gate *Gate) *proxy {
	return &proxy{gate: gate, nodeLinker: link.NewNodeLinker(gate.ctx, &link.Options{
		InsID:    gate.opts.id,
		InsKind:  gcluster.Gate,
		Locator:  gate.opts.locator,
		Registry: gate.opts.registry,
	})}
}

// 绑定用户与网关间的关系
func (p *proxy) bindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.BindGate(ctx, uid, p.gate.opts.id)
	if err != nil {
		return err
	}

	p.trigger(ctx, gcluster.Reconnect, cid, uid)

	return nil
}

// 解绑用户与网关间的关系
func (p *proxy) unbindGate(ctx context.Context, cid, uid int64) error {
	err := p.gate.opts.locator.UnbindGate(ctx, uid, p.gate.opts.id)
	if err != nil {
		glog.Errorf("user unbind failed, gid: %s, cid: %d, uid: %d, err: %v", p.gate.opts.id, cid, uid, err)
	}

	return err
}

// 触发事件
func (p *proxy) trigger(ctx context.Context, event gcluster.Event, cid, uid int64) {
	if gmode.IsDebugMode() {
		glog.Debugf("trigger event, event: %v cid: %d uid: %d", event.String(), cid, uid)
	}

	if err := p.nodeLinker.Trigger(ctx, &link.TriggerArgs{
		Event: event,
		CID:   cid,
		UID:   uid,
	}); err != nil {
		switch {
		case gerrors.Is(err, gerrors.ErrNotFoundEvent), gerrors.Is(err, gerrors.ErrNotFoundUserLocation):
			glog.Warnf("trigger event failed, cid: %d, uid: %d, event: %v, err: %v", cid, uid, event.String(), err)
		default:
			glog.Errorf("trigger event failed, cid: %d, uid: %d, event: %v, err: %v", cid, uid, event.String(), err)
		}
	}
}

// 投递消息
func (p *proxy) deliver(ctx context.Context, cid, uid int64, message []byte) {
	msg, err := gpacket.UnpackMessage(message)
	if err != nil {
		glog.Errorf("unpack message failed: %v", err)
		return
	}

	if gmode.IsDebugMode() {
		glog.Debugf("deliver message, cid: %d uid: %d seq: %d route: %d buffer: %s", cid, uid, msg.Seq, msg.Route, string(msg.Buffer))
	}

	if err = p.nodeLinker.Deliver(ctx, &link.DeliverArgs{
		CID:     cid,
		UID:     uid,
		Route:   msg.Route,
		Message: message,
	}); err != nil {
		switch {
		case gerrors.Is(err, gerrors.ErrNotFoundRoute), gerrors.Is(err, gerrors.ErrNotFoundEndpoint):
			glog.Warnf("deliver message failed, cid: %d uid: %d seq: %d route: %d err: %v", cid, uid, msg.Seq, msg.Route, err)
		default:
			glog.Errorf("deliver message failed, cid: %d uid: %d seq: %d route: %d err: %v", cid, uid, msg.Seq, msg.Route, err)
		}
	}
}

// 开始监听
func (p *proxy) watch() {
	p.nodeLinker.WatchUserLocate()

	p.nodeLinker.WatchClusterInstance()
}

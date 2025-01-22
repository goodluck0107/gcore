package node

import (
	"context"
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/gerrors"
	"gitee.com/monobytes/gcore/gpacket"
)

type provider struct {
	node *Node
}

// Trigger 触发事件
func (p *provider) Trigger(ctx context.Context, gid string, cid, uid int64, event gcluster.Event) error {
	p.node.trigger.trigger(event, gid, cid, uid)

	return nil
}

// Deliver 投递消息
func (p *provider) Deliver(ctx context.Context, gid, nid string, cid, uid int64, message []byte) error {
	msg, err := gpacket.UnpackMessage(message)
	if err != nil {
		return err
	}

	stateful, ok := p.node.router.CheckRouteStateful(msg.Route)
	if !ok {
		if ok = p.node.router.HasDefaultRouteHandler(); !ok {
			return nil
		}
	}

	if stateful {
		if uid == 0 {
			return gerrors.ErrInvalidArgument
		}

		_, ok, err = p.node.proxy.AskNode(ctx, uid, p.node.opts.name, p.node.opts.id)
		if err != nil {
			return err
		}

		if !ok {
			return gerrors.ErrNotFoundSession
		}
	}

	p.node.router.deliver(gid, nid, "", cid, uid, msg.Seq, msg.Route, msg.Buffer)

	return nil
}

// GetState 获取状态
func (p *provider) GetState() (gcluster.State, error) {
	return p.node.getState(), nil
}

// SetState 设置状态
func (p *provider) SetState(state gcluster.State) error {
	return p.node.setState(state)
}

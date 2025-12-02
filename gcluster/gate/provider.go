package gate

import (
	"context"
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gerrors"
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/gsession"
	"github.com/goodluck0107/gcore/gutils/gcall"
)

type provider struct {
	gate *Gate
}

// Bind 绑定用户与网关间的关系
func (p *provider) Bind(ctx context.Context, cid, uid int64) error {
	if cid <= 0 || uid <= 0 {
		return gerrors.ErrInvalidArgument
	}

	err := p.gate.session.Bind(cid, uid)
	if err != nil {
		return err
	}

	err = p.gate.proxy.bindGate(ctx, cid, uid)
	if err != nil {
		_, _ = p.gate.session.Unbind(uid)
	}

	return err
}

// Unbind 解绑用户与网关间的关系
func (p *provider) Unbind(ctx context.Context, uid int64) error {
	if uid == 0 {
		return gerrors.ErrInvalidArgument
	}

	cid, err := p.gate.session.Unbind(uid)
	if err != nil {
		return err
	}

	return p.gate.proxy.unbindGate(ctx, cid, uid)
}

// GetIP 获取客户端IP地址
func (p *provider) GetIP(ctx context.Context, kind gsession.Kind, target int64) (string, error) {
	return p.gate.session.RemoteIP(kind, target)
}

// IsOnline 检测是否在线
func (p *provider) IsOnline(ctx context.Context, kind gsession.Kind, target int64) (bool, error) {
	return p.gate.session.Has(kind, target)
}

// Stat 统计会话总数
func (p *provider) Stat(ctx context.Context, kind gsession.Kind) (int64, error) {
	return p.gate.session.Stat(kind)
}

// Disconnect 断开连接
func (p *provider) Disconnect(ctx context.Context, kind gsession.Kind, target int64, force bool) error {
	return p.gate.session.Close(kind, target, force)
}

// Push 发送消息
func (p *provider) Push(ctx context.Context, kind gsession.Kind, target int64, message []byte) error {
	err := p.gate.session.Push(kind, target, message)

	if kind == gsession.User && gerrors.Is(err, gerrors.ErrNotFoundSession) {
		gcall.Go(func() {
			if err := p.gate.opts.locator.UnbindGate(ctx, target, p.gate.opts.id); err != nil {
				glog.Errorf("unbind gate failed, uid = %d gid = %s err = %v", target, p.gate.opts.id, err)
			}
		})
	}

	return err
}

// Multicast 推送组播消息
func (p *provider) Multicast(ctx context.Context, kind gsession.Kind, targets []int64, message []byte) (int64, error) {
	return p.gate.session.Multicast(kind, targets, message)
}

// Broadcast 推送广播消息
func (p *provider) Broadcast(ctx context.Context, kind gsession.Kind, message []byte) (int64, error) {
	return p.gate.session.Broadcast(kind, message)
}

// GetState 获取状态
func (p *provider) GetState() (gcluster.State, error) {
	return gcluster.Work, nil
}

// SetState 设置状态
func (p *provider) SetState(state gcluster.State) error {
	return nil
}

package gate

import (
	"context"
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/gsession"
)

type Provider interface {
	// Bind 绑定用户与网关间的关系
	Bind(ctx context.Context, cid, uid int64) error
	// Unbind 解绑用户与网关间的关系
	Unbind(ctx context.Context, uid int64) error
	// GetIP 获取客户端IP地址
	GetIP(ctx context.Context, kind gsession.Kind, target int64) (ip string, err error)
	// IsOnline 检测是否在线
	IsOnline(ctx context.Context, kind gsession.Kind, target int64) (isOnline bool, err error)
	// Stat 统计会话总数
	Stat(ctx context.Context, kind gsession.Kind) (total int64, err error)
	// Disconnect 断开连接
	Disconnect(ctx context.Context, kind gsession.Kind, target int64, force bool) error
	// Push 发送消息
	Push(ctx context.Context, kind gsession.Kind, target int64, message []byte) error
	// Multicast 推送组播消息
	Multicast(ctx context.Context, kind gsession.Kind, targets []int64, message []byte) (total int64, err error)
	// Broadcast 推送广播消息
	Broadcast(ctx context.Context, kind gsession.Kind, message []byte) (total int64, err error)
	// GetState 获取状态
	GetState() (gcluster.State, error)
	// SetState 设置状态
	SetState(state gcluster.State) error
}

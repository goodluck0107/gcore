package gate_test

import (
	"context"
	"fmt"
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gsession"
	"github.com/goodluck0107/gcore/internal/transporter/gate"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	server, err := gate.NewServer(":49899", &provider{})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("server listen on: %s", server.ListenAddr())

	go server.Start()

	<-time.After(20 * time.Second)
}

type provider struct {
}

// Bind 绑定用户与网关间的关系
func (p *provider) Bind(ctx context.Context, cid, uid int64) error {
	return nil
}

// Unbind 解绑用户与网关间的关系
func (p *provider) Unbind(ctx context.Context, uid int64) error {
	return nil
}

// GetIP 获取客户端IP地址
func (p *provider) GetIP(ctx context.Context, kind gsession.Kind, target int64) (ip string, err error) {
	fmt.Println(kind, target)
	ip = "192.168.0.88"
	return
}

// IsOnline 检测是否在线
func (p *provider) IsOnline(ctx context.Context, kind gsession.Kind, target int64) (isOnline bool, err error) {
	return
}

// Push 发送消息（异步）
func (p *provider) Push(ctx context.Context, kind gsession.Kind, target int64, message []byte) error {
	//fmt.Println(kind, target, message)

	return nil
}

// Multicast 推送组播消息（异步）
func (p *provider) Multicast(ctx context.Context, kind gsession.Kind, targets []int64, message []byte) (total int64, err error) {
	return
}

// Broadcast 推送广播消息（异步）
func (p *provider) Broadcast(ctx context.Context, kind gsession.Kind, message []byte) (total int64, err error) {
	return
}

// DoStat 统计会话总数
func (p *provider) Stat(ctx context.Context, kind gsession.Kind) (total int64, err error) {
	return
}

// Disconnect 断开连接
func (p *provider) Disconnect(ctx context.Context, kind gsession.Kind, target int64, force bool) error {
	return nil
}

// GetState 获取状态
func (p *provider) GetState() (gcluster.State, error) {
	return gcluster.Work, nil
}

// SetState 设置状态
func (p *provider) SetState(state gcluster.State) error {
	return nil
}

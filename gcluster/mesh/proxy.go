package mesh

import (
	"context"
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/gsession"
	"gitee.com/monobytes/gcore/gtransport"
	"gitee.com/monobytes/gcore/internal/link"
)

type Proxy struct {
	mesh       *Mesh            // 微服务器
	gateLinker *link.GateLinker // 网关链接器
	nodeLinker *link.NodeLinker // 节点链接器
}

func newProxy(mesh *Mesh) *Proxy {
	opts := &link.Options{
		InsID:     mesh.opts.id,
		InsKind:   gcluster.Mesh,
		Codec:     mesh.opts.codec,
		Locator:   mesh.opts.locator,
		Registry:  mesh.opts.registry,
		Encryptor: mesh.opts.encryptor,
	}

	return &Proxy{
		mesh:       mesh,
		gateLinker: link.NewGateLinker(mesh.opts.ctx, opts),
		nodeLinker: link.NewNodeLinker(mesh.opts.ctx, opts),
	}
}

// GetID 获取当前实例ID
func (p *Proxy) GetID() string {
	return p.mesh.opts.id
}

// GetName 获取当前实例名称
func (p *Proxy) GetName() string {
	return p.mesh.opts.name
}

// AddServiceProvider 添加服务提供者
func (p *Proxy) AddServiceProvider(name string, desc interface{}, provider interface{}) {
	p.mesh.addServiceProvider(name, desc, provider)
}

// AddHookListener 添加钩子监听器
func (p *Proxy) AddHookListener(hook gcluster.Hook, handler HookHandler) {
	p.mesh.addHookListener(hook, handler)
}

// NewMeshClient 新建微服务客户端
// target参数可分为三种模式:
// 服务直连模式: 	direct://127.0.0.1:8011
// 服务直连模式: 	direct://711baf8d-8a06-11ef-b7df-f4f19e1f0070
// 服务发现模式: 	discovery://service_name
func (p *Proxy) NewMeshClient(target string) (gtransport.Client, error) {
	return p.mesh.opts.transporter.NewClient(target)
}

// BindGate 绑定网关
func (p *Proxy) BindGate(ctx context.Context, gid string, cid, uid int64) error {
	return p.gateLinker.Bind(ctx, gid, cid, uid)
}

// UnbindGate 解绑网关
func (p *Proxy) UnbindGate(ctx context.Context, uid int64) error {
	return p.gateLinker.Unbind(ctx, uid)
}

// BindNode 绑定节点
// 单个用户可以绑定到多个节点服务器上，相同名称的节点服务器只能绑定一个，多次绑定会到相同名称的节点服务器会覆盖之前的绑定。
// 绑定操作会通过发布订阅方式同步到网关服务器和其他相关节点服务器上。
func (p *Proxy) BindNode(ctx context.Context, uid int64, name, nid string) error {
	return p.nodeLinker.Bind(ctx, uid, name, nid)
}

// UnbindNode 解绑节点
// 解绑时会对对应名称的节点服务器进行解绑，解绑时会对解绑节点ID进行校验，不匹配则解绑失败。
// 解绑操作会通过发布订阅方式同步到网关服务器和其他相关节点服务器上。
func (p *Proxy) UnbindNode(ctx context.Context, uid int64, name, nid string) error {
	return p.nodeLinker.Unbind(ctx, uid, name, nid)
}

// LocateGate 定位用户所在网关
func (p *Proxy) LocateGate(ctx context.Context, uid int64) (string, error) {
	return p.gateLinker.Locate(ctx, uid)
}

// AskGate 检测用户是否在给定的网关上
func (p *Proxy) AskGate(ctx context.Context, gid string, uid int64) (string, bool, error) {
	return p.gateLinker.Ask(ctx, gid, uid)
}

// LocateNode 定位用户所在节点
func (p *Proxy) LocateNode(ctx context.Context, uid int64, name string) (string, error) {
	return p.nodeLinker.Locate(ctx, uid, name)
}

// AskNode 检测用户是否在给定的节点上
func (p *Proxy) AskNode(ctx context.Context, uid int64, name, nid string) (string, bool, error) {
	return p.nodeLinker.Ask(ctx, uid, name, nid)
}

// FetchGateList 拉取网关列表
func (p *Proxy) FetchGateList(ctx context.Context, states ...gcluster.State) ([]*gregistry.ServiceInstance, error) {
	return p.gateLinker.FetchGateList(ctx, states...)
}

// FetchNodeList 拉取节点列表
func (p *Proxy) FetchNodeList(ctx context.Context, states ...gcluster.State) ([]*gregistry.ServiceInstance, error) {
	return p.nodeLinker.FetchNodeList(ctx, states...)
}

// PackMessage 打包消息
func (p *Proxy) PackMessage(message *gcluster.Message) ([]byte, error) {
	buf, err := p.gateLinker.PackMessage(message, true)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// PackBuffer 打包Buffer
func (p *Proxy) PackBuffer(message any) ([]byte, error) {
	return p.gateLinker.PackBuffer(message, true)
}

// GetIP 获取客户端IP
func (p *Proxy) GetIP(ctx context.Context, args *gcluster.GetIPArgs) (string, error) {
	return p.gateLinker.GetIP(ctx, args)
}

// Stat 统计会话总数
func (p *Proxy) Stat(ctx context.Context, kind gsession.Kind) (int64, error) {
	return p.gateLinker.Stat(ctx, kind)
}

// IsOnline 检测是否在线
func (p *Proxy) IsOnline(ctx context.Context, args *gcluster.IsOnlineArgs) (bool, error) {
	return p.gateLinker.IsOnline(ctx, args)
}

// Disconnect 断开连接
func (p *Proxy) Disconnect(ctx context.Context, args *gcluster.DisconnectArgs) error {
	return p.gateLinker.Disconnect(ctx, args)
}

// Push 推送消息
func (p *Proxy) Push(ctx context.Context, args *gcluster.PushArgs) error {
	return p.gateLinker.Push(ctx, args)
}

// Multicast 推送组播消息
func (p *Proxy) Multicast(ctx context.Context, args *gcluster.MulticastArgs) error {
	return p.gateLinker.Multicast(ctx, args)
}

// Broadcast 推送广播消息
func (p *Proxy) Broadcast(ctx context.Context, args *gcluster.BroadcastArgs) error {
	return p.gateLinker.Broadcast(ctx, args)
}

// Deliver 投递消息给节点处理
func (p *Proxy) Deliver(ctx context.Context, args *gcluster.DeliverArgs) error {
	return p.nodeLinker.Deliver(ctx, &link.DeliverArgs{
		NID:     args.NID,
		UID:     args.UID,
		Route:   args.Message.Route,
		Message: args.Message,
	})
}

// 开始监听
func (p *Proxy) watch() {
	p.gateLinker.WatchUserLocate()

	p.gateLinker.WatchClusterInstance()

	p.nodeLinker.WatchUserLocate()

	p.nodeLinker.WatchClusterInstance()
}

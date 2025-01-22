package link

import (
	"gitee.com/monobytes/gcore/gcluster"
)

type (
	Message        = gcluster.Message
	GetIPArgs      = gcluster.GetIPArgs
	IsOnlineArgs   = gcluster.IsOnlineArgs
	DisconnectArgs = gcluster.DisconnectArgs
	PushArgs       = gcluster.PushArgs
	MulticastArgs  = gcluster.MulticastArgs
	BroadcastArgs  = gcluster.BroadcastArgs
)

type DeliverArgs struct {
	NID     string      // 接收节点。存在接收节点时，消息会直接投递给接收节点；不存在接收节点时，系统定位用户所在节点，然后投递。
	CID     int64       // 连接ID
	UID     int64       // 用户ID
	Route   int32       // 路由
	Message interface{} // 消息
}

type TriggerArgs struct {
	Event gcluster.Event // 事件
	CID   int64          // 连接ID
	UID   int64          // 用户ID
}

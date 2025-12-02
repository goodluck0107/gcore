package client

import "github.com/goodluck0107/gcore/gcluster"

type Options struct {
	Addr         string        // 连接地址
	InsID        string        // 实例ID
	InsKind      gcluster.Kind // 实例类型
	CloseHandler func()        // 关闭处理器
}

package link

import (
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/gcrypto"
	"gitee.com/monobytes/gcore/gencoding"
	"gitee.com/monobytes/gcore/glocate"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/internal/dispatcher"
)

type Options struct {
	InsID           string                     // 实例ID
	InsKind         gcluster.Kind              // 实例类型
	Codec           gencoding.Codec            // 编解码器
	Locator         glocate.Locator            // 定位器
	Registry        gregistry.Registry         // 注册器
	Encryptor       gcrypto.Encryptor          // 加密器
	BalanceStrategy dispatcher.BalanceStrategy // 负载均衡策略
}

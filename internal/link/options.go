package link

import (
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gcrypto"
	"github.com/goodluck0107/gcore/gencoding"
	"github.com/goodluck0107/gcore/glocate"
	"github.com/goodluck0107/gcore/gregistry"
	"github.com/goodluck0107/gcore/internal/dispatcher"
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

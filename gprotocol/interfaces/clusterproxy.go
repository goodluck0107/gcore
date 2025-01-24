package interfaces

import (
	"errors"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gtransport"
)

var (
	rpcProvider   RpcProvider   = &defaultClusterProxy{}
	rpcDiscoverer RpcDiscoverer = &defaultClusterProxy{}
)

type RpcProvider interface {
	AddServiceProvider(name string, desc interface{}, provider interface{})
}

type RpcDiscoverer interface {
	NewMeshClient(target string) (gtransport.Client, error)
}

func SetClusterProxy(proxy any) {
	if provider, ok := proxy.(RpcProvider); ok {
		rpcProvider = provider
	}
	if discoverer, ok := proxy.(RpcDiscoverer); ok {
		rpcDiscoverer = discoverer
	}
}

func GetRpcProvider() RpcProvider {
	return rpcProvider
}

func GetRpcDiscoverer() RpcDiscoverer {
	return rpcDiscoverer
}

type defaultClusterProxy struct{}

func (d *defaultClusterProxy) AddServiceProvider(name string, desc interface{}, provider interface{}) {
	glog.Warnf("rpc provider not init yet, please call gprotocol.SetClusterProxy() first")
}

func (d *defaultClusterProxy) NewMeshClient(target string) (gtransport.Client, error) {
	return nil, errors.New("rpc discoverer not init yet, please call gprotocol.SetClusterProxy() first")
}

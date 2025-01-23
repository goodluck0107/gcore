package interfaces

import (
	"errors"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gtransport"
	"reflect"
)

var rpcProvider RpcProvider = &defaultRpcProvider{}

var rpcClients map[string]reflect.Value

func init() {
	rpcClients = make(map[string]reflect.Value)
}

type RpcProvider interface {
	AddServiceProvider(name string, desc interface{}, provider interface{})
	NewMeshClient(target string) (gtransport.Client, error)
}

func InitRpcProvider(p RpcProvider) {
	rpcProvider = p
}

func GetRpcProvider() RpcProvider {
	return rpcProvider
}

func AddRpcClient(name string, clientNewFunc interface{}) {
	rpcClients[name] = reflect.ValueOf(clientNewFunc)
}

type defaultRpcProvider struct{}

func (d *defaultRpcProvider) AddServiceProvider(name string, desc interface{}, provider interface{}) {
	glog.Fatalf("rpc provider not init yet, please call gprotocol.InitRpcProvider() first")
}

func (d *defaultRpcProvider) NewMeshClient(target string) (gtransport.Client, error) {
	return nil, errors.New("rpc provider not init yet, please call gprotocol.InitRpcProvider() first")
}

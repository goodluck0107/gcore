package interfaces

import "reflect"

var rpcClients map[string]reflect.Value

func init() {
	rpcClients = make(map[string]reflect.Value)
}

func AddRpcClientGenerator(name string, generator interface{}) {
	rpcClients[name] = reflect.ValueOf(generator)
}

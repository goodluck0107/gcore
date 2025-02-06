package handler

import (
	"context"
	"gitee.com/monobytes/gcore/gcodes"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type Handler func(ctx context.Context, dec func(interface{}) error) (interface{}, *gcodes.Code)

type Result interface {
	GetCode() int
	GetMsg() string
	GetData() proto.Message
}

type Manager interface {
	RangeURLHandlers(do func(md Metadata, handler Handler))
}

type result struct {
	Code int           `json:"code,omitempty"`
	Msg  string        `json:"msg,omitempty"`
	Data proto.Message `json:"data,omitempty"`
}

func (r *result) GetCode() int {
	if r != nil {
		return r.Code
	}
	return 0
}

func (r *result) GetMsg() string {
	if r != nil {
		return r.Msg
	}
	return ""
}

func (r *result) GetData() proto.Message {
	if r != nil {
		return r.Data
	}
	return nil
}

type Metadata struct {
	Cmd        int32
	Uri        string
	Srv        interface{}
	AuthType   int32
	HTTPMethod string
	Req        reflect.Type
	Rsp        reflect.Type
}

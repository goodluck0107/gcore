package interfaces

import "google.golang.org/protobuf/proto"

type (
	MethodItem struct {
		Method string
		HTTP   string
		Auth   int32
		Cmd    int32
	}

	ReqItem struct {
		Req  func() proto.Message
		Rsp  func() proto.Message
		Auth int32
		Name string
		HTTP string
	}
)

type MsgDef interface {
	GetMethodRouter() map[string]*MethodItem
	GetIdMsg() map[int32]*ReqItem
}

package handler

import (
	"context"
	"gitee.com/monobytes/gcore/gcodes"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gprotocol/interfaces"
	"github.com/jinzhu/copier"
	"github.com/smallnest/rpcx/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
	"sync"
)

var (
	singletonGRPCHandlersMgr *GRPCHandlersMgr
	initOnceGRPCHandlersMgr  sync.Once
)

type MethodInfo struct {
	Metadata
	Handler MethodHandler
	Desc    *grpc.ServiceDesc
}

func GetGRPCHandlersMgr() *GRPCHandlersMgr {
	initOnceGRPCHandlersMgr.Do(func() {
		singletonGRPCHandlersMgr = &GRPCHandlersMgr{
			routerMap: map[string]*interfaces.MethodItem{},
			uriMap:    map[string]*MethodInfo{},
			cmdMap:    map[int32]*MethodInfo{},
		}
	})
	return singletonGRPCHandlersMgr
}

type MethodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)

// GRPCHandlersMgr 管理接口到方法的对应关系
type GRPCHandlersMgr struct {
	routerMap map[string]*interfaces.MethodItem
	reqMap    map[int32]*interfaces.ReqItem
	uriMap    map[string]*MethodInfo
	cmdMap    map[int32]*MethodInfo
}

func (s *GRPCHandlersMgr) RegisterMsg(msg interfaces.MsgDef) {
	s.routerMap = msg.GetMethodRouter()
	s.reqMap = msg.GetIdMsg()
}

func (s *GRPCHandlersMgr) RegisterServer(desc *grpc.ServiceDesc, srv interface{}) {
	if srv != nil {
		ht := reflect.TypeOf(desc.HandlerType).Elem()
		st := reflect.TypeOf(srv)
		if !st.Implements(ht) {
			log.Fatalf("grpc handler RegisterServer found the handler of type %v that does not satisfy %v", st, ht)
		}
	}

	mapMethodItem := map[string]grpc.MethodDesc{}
	for _, methodItem := range desc.Methods {
		fullName := s.basename(desc.ServiceName) + "." + methodItem.MethodName
		mapMethodItem[fullName] = methodItem
	}

	for url, configItem := range s.routerMap {
		fullName := configItem.Method
		cmd := configItem.Cmd
		authType := configItem.Auth
		httpMethod := configItem.HTTP
		if methodItem, ok := mapMethodItem[fullName]; ok {
			req, ok := s.reqMap[cmd]
			if !ok {
				glog.Fatalf("server method %s, cmd: %d  req mismatch", fullName, cmd)
				return
			}
			methodInfo := &MethodInfo{
				Metadata: Metadata{
					Cmd:        cmd,
					Uri:        url,
					Srv:        srv,
					AuthType:   authType,
					HTTPMethod: httpMethod,
					Req:        reflect.TypeOf(req.Req()).Elem(),
					Rsp:        reflect.TypeOf(req.Rsp()).Elem(),
				},
				Handler: MethodHandler(methodItem.Handler),
				Desc:    desc,
			}
			s.uriMap[url] = methodInfo
			s.cmdMap[cmd] = methodInfo
		}
	}
}

func (s *GRPCHandlersMgr) basename(fullname string) string {
	arr := strings.Split(fullname, ".")
	length := len(arr)
	if length == 0 {
		return fullname
	}
	return arr[length-1]
}

func (s *GRPCHandlersMgr) GetHandlerByRouter(url string) (*MethodInfo, bool) {
	ret, ok := s.uriMap[url]
	return ret, ok
}

func (s *GRPCHandlersMgr) GetHandlerById(id int32) (*MethodInfo, bool) {
	ret, ok := s.cmdMap[id]
	return ret, ok
}

func (s *GRPCHandlersMgr) RangeURLHandlers(do func(md Metadata, handler Handler)) {
	for _, methodInfo := range s.uriMap {
		do(methodInfo.Metadata, s.GenerateHandler(methodInfo))
	}
}

func (s *GRPCHandlersMgr) GenerateHandler(methodInfo *MethodInfo) Handler {
	var (
		handler = methodInfo.Handler
		srv     = methodInfo.Srv
	)
	return func(ctx context.Context, req proto.Message) Result {
		ret, callErr := handler(srv, ctx, func(i interface{}) error {
			return copier.Copy(i, req)
		}, nil)
		var res result
		codex := gcodes.Convert(callErr)
		res.Code = codex.Code()
		if codex != gcodes.OK {
			res.Msg = codex.Message()
		} else {
			retMsg, ok := ret.(proto.Message)
			if !ok {
				codex = gcodes.InternalError
				res.Code = codex.Code()
				res.Msg = codex.Message()
				glog.Errorf("grpc handle return err name:%s, result struct not impl proto.Message yet", methodInfo.Uri)
			} else {
				res.Data = retMsg
			}
		}
		return &res
	}
}

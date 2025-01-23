package genmsg

import (
	"fmt"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"protoc-gen-idmsg/pb"

	"google.golang.org/protobuf/proto"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func getGenSvcName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sClientX", firstUpper(svc.GetName()))
}

func getNativeSvcName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sClient", firstUpper(svc.GetName()))
}

func getInnerSvcName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sClientX", firstLower(svc.GetName()))
}

func getInputType(meth *descriptor.MethodDescriptorProto) string {
	return typeStrip(meth.GetInputType())
}

func getOutputType(meth *descriptor.MethodDescriptorProto) string {
	return typeStrip(meth.GetOutputType())
}

func extractIdOptions(meth *descriptor.MethodDescriptorProto) (*pb.ProtoIdRule, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, pb.E_Id) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, pb.E_Id)
	opts, ok := ext.(*pb.ProtoIdRule)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an ProtoIdRule", ext)
	}
	return opts, nil
}

type RouterInfo struct {
	Method string
	Route  string
}

func getMethRouter(meth *descriptor.MethodDescriptorProto) RouterInfo {
	var routeInfo RouterInfo
	mOption, _ := extractAPIOptions(meth)
	if mOption == nil {
		return routeInfo
	}
	switch method := mOption.GetPattern().(type) {

	case *options.HttpRule_Get:
		routeInfo.Method = "GET"
		routeInfo.Route = method.Get
	case *options.HttpRule_Put:
		routeInfo.Method = "PUT"
		routeInfo.Route = method.Put
	case *options.HttpRule_Post:
		routeInfo.Method = "POST"
		routeInfo.Route = method.Post
	case *options.HttpRule_Delete:
		routeInfo.Method = "DELETE"
		routeInfo.Route = method.Delete
	case *options.HttpRule_Patch:
		routeInfo.Method = "PATCH"
		routeInfo.Route = method.Patch
	case *options.HttpRule_Custom:
		routeInfo.Method = "CUSTOM"
		routeInfo.Route = method.Custom.Path
	}
	return routeInfo
}

func getMethodCmd(meth *descriptor.MethodDescriptorProto) int32 {
	mOption, _ := extractIdOptions(meth)
	if mOption == nil {
		return 0
	}
	return int32(mOption.Cmd)
}

func getMethAuthType(meth *descriptor.MethodDescriptorProto) *pb.AuthType {
	mOption, _ := extractAuthMethodOptions(meth)
	return mOption
}

func getServiceAuth(service *descriptor.ServiceDescriptorProto) int32 {
	opts, err := extractAuthServiceOptions(service)

	if err == nil && opts != 0 {
		return int32(opts)
	}
	return 0
}

func getMethodAuth(serviceAuth int32, meth *descriptor.MethodDescriptorProto) int32 {
	opts, err := extractAuthMethodOptions(meth)
	var mthAuth int32
	if err == nil && opts != nil {
		mthAuth = int32(*opts)
	} else {
		mthAuth = serviceAuth
	}
	return mthAuth
}

func getMethodName(meth *descriptor.MethodDescriptorProto) string {
	return *meth.Name
}

func extractAuthServiceOptions(service *descriptor.ServiceDescriptorProto) (pb.AuthType, error) {
	if service.Options == nil {
		return 0, nil
	}
	if !proto.HasExtension(service.Options, pb.E_AuthService) {
		return 0, nil
	}
	ext := proto.GetExtension(service.Options, pb.E_AuthService)
	opts, ok := ext.(pb.AuthType)
	if !ok {
		return 0, fmt.Errorf("extension is %T; want an pb.AuthType", ext)
	}
	return opts, nil
}

func extractAuthMethodOptions(meth *descriptor.MethodDescriptorProto) (*pb.AuthType, error) {
	if meth.Options == nil {
		return nil, nil
	}
	if !proto.HasExtension(meth.Options, pb.E_AuthMethod) {
		return nil, nil
	}
	ext := proto.GetExtension(meth.Options, pb.E_AuthMethod)
	opts, ok := ext.(*pb.AuthType)
	if !ok {
		return nil, fmt.Errorf("extension is %T; want an HttpRule", ext)
	}
	return opts, nil
}

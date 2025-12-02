package genmsg

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

type RpcProtocol struct {
	Code       int32
	Name       string
	ReqType    string
	RspType    string
	SvcName    string
	MthName    string
	MthComment string
	MthAuth    int32
}

type FileData struct {
	File     *File
	Services []*descriptor.ServiceDescriptorProto
}

type RpcData struct {
	files  []*FileData
	suffix string
}

// 首字母大写
func firstUpper(s string) string {
	if s == "" {
		return s
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func firstLower(s string) string {
	if s == "" {
		return s
	}
	if len(s) == 1 {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[0:1]) + s[1:]
}

func (g *RpcData) GenRpc() []*plugin.CodeGeneratorResponse_File {
	var ret []*plugin.CodeGeneratorResponse_File
	for _, item := range g.files {
		if len(item.Services) > 0 {
			f := g.generateFileTpl(item)
			ret = append(ret, f)
		}
	}
	return ret
}

func (g *RpcData) GenRpcV2() []*plugin.CodeGeneratorResponse_File {
	var ret []*plugin.CodeGeneratorResponse_File
	for _, item := range g.files {
		if len(item.Services) > 0 {
			f := g.generateFileTplV2(item)
			ret = append(ret, f)
		}
	}
	return ret
}

func (g *RpcData) getGenFileName(f *File) string {
	prefix := strings.ReplaceAll(f.GetName(), ".proto", "")
	return fmt.Sprintf("%s_grpcx.pb.go", prefix)
}

func (g *RpcData) getServiceName(svc *descriptor.ServiceDescriptorProto) string {
	return svc.GetName()
}

func (g *RpcData) getNativeCliName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sClient", firstUpper(svc.GetName()))
}

func (g *RpcData) getNativeSvcName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sServer", firstUpper(svc.GetName()))
}

func (g *RpcData) getGenSvcName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sClientX", firstUpper(svc.GetName()))
}

func (g *RpcData) getInnerSvcName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sClientX", firstLower(svc.GetName()))
}

func (g *RpcData) getGenFileNameV2(f *File) string {
	prefix := strings.ReplaceAll(f.GetName(), ".proto", "")
	return fmt.Sprintf("%s_grpc.wrap.go", prefix)
}

func (g *RpcData) getGenCliNameV2(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sClientDiscovery", firstUpper(svc.GetName()))
}

func (g *RpcData) getGenSvcNameV2(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%sServerProvider", firstUpper(svc.GetName()))
}

func (g *RpcData) getSvcDescName(svc *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("%s_ServiceDesc", firstUpper(svc.GetName()))
}

func (g *RpcData) getInputType(meth *descriptor.MethodDescriptorProto) string {
	return typeStrip(meth.GetInputType())
}

func (g *RpcData) getOutputType(meth *descriptor.MethodDescriptorProto) string {
	return typeStrip(meth.GetOutputType())
}

// 使用模板渲染代码，未完成
func (g *RpcData) renderCode(fileData *FileData) string {
	tplStr :=
		`package pb
import (
	"context"
	"mist-server/core/log"
	"mist-server/core/client"
	"mist-server/core/grpcx"
	"mist-server/core/interfaces"
	"time"
)

{{- range .Services}}
type {{getGenSvcName .}} interface {
	{{- range .Method}}
	{{firstUpper .Name}}(ctx context.Context, req *{{getInputType .}}, optionSlice... grpcx.ShortOption) (*{{getOutputType .}}, error)
	{{- end}}
}
{{- end}}
{{- range .Services}}
{{- $innerSvcName := getInnerSvcName .}}
{{- $nativeCliName := getNativeCliName .}}
{{- $genSvcName := getGenSvcName .}}

type {{$innerSvcName}} struct {
	allocator interfaces.EntityAllocator
	defaultServerType string
	stubName string
}
func New{{$genSvcName}}(a interfaces.EntityAllocator, serverType string, stubName string) {{$genSvcName}} {
	return &{{$innerSvcName}}{
		allocator:      a,
		defaultServerType: serverType,
		stubName:       stubName,
	}
}
func (c *{{$innerSvcName}}) getCallContext(ctx context.Context, opts *grpcx.ShortOptions) (context.Context, context.CancelFunc) {
	var (
		destServerType string
		oldServerId    string

		pid = opts.Pid
		eid = opts.Eid
	)
	if opts.ServerType != "" {
		destServerType = opts.ServerType
	} else {
		destServerType = c.defaultServerType
	}
	if opts.ServerId != "" {
		oldServerId = opts.ServerId
	} else {
		oldServerId, _ = c.allocator.GetServer(ctx, destServerType, eid)
	}

	newServerId, err := c.allocator.CheckAndAllocate(ctx, destServerType, eid, oldServerId)
	if err != nil {
		log.Error("GetCallContext serverType:%s pid:%s eid:%s err:%v", destServerType, pid, eid, err)
		return context.WithTimeout(context.TODO(), time.Millisecond*10)
	}
	rCtx, cancel := grpcx.MakeGrpCtx(ctx,
		grpcx.WithServerId(newServerId),
		grpcx.WithFormerServerId(oldServerId),
		grpcx.WithPlayerId(pid),
		grpcx.WithUniqueId(eid),
	)
	return rCtx, cancel
}
func (c *{{$innerSvcName}}) getCli(opts *grpcx.ShortOptions) {{$nativeCliName}} {
	serverType := opts.ServerType
	if opts.ServerType == "" {
		serverType = c.defaultServerType
	}
	cli, _ := client.GRPCClient().GetStub(serverType, c.stubName).({{$nativeCliName}})
	return cli
}

{{range .Method}}
func (c *{{$innerSvcName}}) {{firstUpper .Name}}(ctx context.Context, req *{{getInputType .}}, optionSlice... grpcx.ShortOption) (*{{getOutputType .}}, error) {
	opts := grpcx.GetOpts(optionSlice)
	rCtx, cancel := c.getCallContext(ctx, opts)
	defer cancel()
	return c.getCli(opts).{{firstUpper .Name}}(rCtx, req)
}
{{end}}
{{end}}
`
	tplFuncMap := map[string]interface{}{
		"getInnerSvcName":  g.getInnerSvcName,
		"getNativeCliName": g.getNativeCliName,
		"getGenSvcName":    g.getGenSvcName,
		"getInputType":     g.getInputType,
		"getOutputType":    g.getOutputType,
		"firstUpper":       firstUpper,
	}
	tpl := template.Must(template.New("tmpl").Funcs(tplFuncMap).Parse(tplStr))
	buf := &bytes.Buffer{}

	if err := tpl.Execute(buf, fileData); err != nil {
		fmt.Printf("template err:%v\n", err)
	}
	return buf.String()
}

// 使用模板渲染代码，未完成
func (g *RpcData) renderCodeV2(fileData *FileData) string {
	tplStr :=
		`// Code generated by protoc-gen-idmsg. DO NOT EDIT.

package pb

import (
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/gprotocol/handler"
	"github.com/goodluck0107/gcore/gprotocol/interfaces"
	"google.golang.org/grpc"
)

{{- range .Services}}
{{- $serviceName := getServiceName .}}
{{- $innerSvcName := getInnerSvcName .}}
{{- $nativeCliName := getNativeCliName .}}
{{- $nativeSvcName := getNativeSvcName .}}
{{- $genCliName := getGenCliName .}}
{{- $genSvcName := getGenSvcName .}}
{{- $svcDescName := getSvcDescName .}}

func init() {
	handler.GetGRPCHandlersMgr().RegisterMsg(MsgDef)
	interfaces.AddRpcClientGenerator("{{$serviceName}}", New{{$nativeCliName}})
}

func New{{$genCliName}}() ({{$nativeCliName}}, error) {
	serviceName := "{{$serviceName}}"
	target := "discovery://" + serviceName
	client, err := interfaces.GetRpcDiscoverer().NewMeshClient(target)
	if err != nil {
		glog.Errorf("failed to discovery {{$nativeCliName}} named: %s! err:%v", serviceName, err)
		return nil, err
	}
	return New{{$nativeCliName}}(client.Client().(grpc.ClientConnInterface)), nil
}

func Add{{$genSvcName}}(svr {{$nativeSvcName}}) {
	serviceName := "{{$serviceName}}"
	handler.GetGRPCHandlersMgr().RegisterServer(&{{$svcDescName}}, svr)
	interfaces.GetRpcProvider().AddServiceProvider(serviceName, &{{$svcDescName}}, svr)
}

{{end}}
`
	tplFuncMap := map[string]interface{}{
		"getServiceName":   g.getServiceName,
		"getInnerSvcName":  g.getInnerSvcName,
		"getNativeCliName": g.getNativeCliName,
		"getNativeSvcName": g.getNativeSvcName,
		"getGenCliName":    g.getGenCliNameV2,
		"getGenSvcName":    g.getGenSvcNameV2,
		"getSvcDescName":   g.getSvcDescName,
		"getInputType":     g.getInputType,
		"getOutputType":    g.getOutputType,
		"firstUpper":       firstUpper,
	}
	tpl := template.Must(template.New("tmpl").Funcs(tplFuncMap).Parse(tplStr))
	buf := &bytes.Buffer{}

	if err := tpl.Execute(buf, fileData); err != nil {
		fmt.Printf("template err:%v\n", err)
	}
	return buf.String()
}

func (g *RpcData) generateFileTpl(fileData *FileData) *plugin.CodeGeneratorResponse_File {
	str := g.renderCode(fileData)
	fileName := g.getGenFileName(fileData.File)
	var ret = &plugin.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: &str,
	}
	return ret
}

func (g *RpcData) generateFileTplV2(fileData *FileData) *plugin.CodeGeneratorResponse_File {
	str := g.renderCodeV2(fileData)
	fileName := g.getGenFileNameV2(fileData.File)
	var ret = &plugin.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: &str,
	}
	return ret
}

func (g *RpcData) AddServices(file *File, services []*descriptor.ServiceDescriptorProto) {
	fileData := &FileData{
		File:     file,
		Services: services,
	}
	g.files = append(g.files, fileData)
}

func (g *RpcData) FromReg(reg *Registry) {
	files := reg.GetFiles()
	for _, file := range files {
		file.SourceCodeInfo.ProtoReflect()
		g.AddServices(file, file.GetService())
	}
}

func NewRpcData() *RpcData {
	return &RpcData{}
}

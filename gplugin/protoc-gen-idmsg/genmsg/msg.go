package genmsg

import (
	"bytes"
	"fmt"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"sort"
	"strings"
	"text/template"

	"github.com/golang/glog"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

type Protocol struct {
	Code       int32
	Method     string
	Name       string
	ReqType    string
	RspType    string
	SvcName    string
	MthName    string
	MthComment string
	MthAuth    int32
	IdValue    int32
	IdName     string
}

type GenData struct {
	CmdTypeName string
	MapProto    map[string]*Protocol
	ProtoList   []string // 有序
	cmdType     *Enum
	cmdIdValue  map[int32]*descriptor.EnumValueDescriptorProto
	Services    []*descriptor.ServiceDescriptorProto
}

func NewGenData() *GenData {
	ret := &GenData{
		MapProto:   map[string]*Protocol{},
		cmdIdValue: map[int32]*descriptor.EnumValueDescriptorProto{},
	}
	return ret
}

func (ts *GenData) Sort() {
	sort.Slice(ts.ProtoList, func(i, j int) bool {
		return ts.MapProto[ts.ProtoList[i]].Code < ts.MapProto[ts.ProtoList[j]].Code
	})

	_ = ts.ProtoList
}

// FindEnumType 找到定义 cmd 的枚举类型
func (ts *GenData) FindEnumType(reg *Registry) *Enum {
	// 找到所有方法的 cmd 定义
	var cmdSlice []int32
	for _, item := range reg.files {
		for _, svc := range item.Services {
			for _, mth := range svc.Methods {
				protoIdRule, _ := extractIdOptions(mth.MethodDescriptorProto)
				if protoIdRule != nil {
					cmdSlice = append(cmdSlice, int32(protoIdRule.Cmd))
				}
			}
		}
	}

	// 遍历所有 enum，需要 cmd 的所有值在这个枚举里都定义过
	for _, item := range reg.enums {
		var int32Values []int32
		for _, v := range item.Value {
			int32Values = append(int32Values, *v.Number)
		}

		if ts.CheckInt32SliceIn(cmdSlice, int32Values) {
			return item
		}
	}
	return nil
}

func (ts *GenData) CheckInt32SliceIn(sliceChild []int32, sliceParent []int32) bool {
	if len(sliceParent) < len(sliceChild) {
		return false
	}

	var m = map[int32]bool{}
	for _, item := range sliceParent {
		m[item] = true
	}

	for _, item := range sliceChild {
		if !m[item] {
			return false
		}
	}
	return true
}

// FromReg 解析文件
func (ts *GenData) FromReg(reg *Registry) {
	files := reg.GetFiles()
	enumType := ts.FindEnumType(reg)
	ts.SetCmd(enumType)

	for _, file := range files {
		file.SourceCodeInfo.ProtoReflect()
		ts.AddServices(file.GetService())
		for _, service := range file.GetService() {
			serviceAuth, err := extractAuthServiceOptions(service)
			if err != nil {
				glog.Error(err)
			}

			for _, meth := range service.GetMethod() {
				optionId, _ := extractIdOptions(meth)
				if optionId == nil {
					continue
				}
				optionHttp, _ := extractAPIOptions(meth)
				if optionHttp == nil {
					continue
				}
				methAuth, _ := extractAuthMethodOptions(meth)
				auth := int32(0)
				if methAuth != nil {
					auth = int32(*methAuth)
				}

				if auth == 0 && serviceAuth != 0 {
					auth = int32(serviceAuth)
				}

				p := &Protocol{
					Code:    int32(optionId.Cmd),
					ReqType: typeStrip(*meth.InputType),
					RspType: typeStrip(*meth.OutputType),
					SvcName: *service.Name,
					MthName: *meth.Name,
					IdValue: int32(optionId.Cmd),
					IdName:  *ts.GetEnumValue(int32(optionId.Cmd)).Name,
					MthAuth: auth,
				}

				switch method := optionHttp.GetPattern().(type) {
				case *options.HttpRule_Get:
					p.Method = "GET"
					p.Name = method.Get
				case *options.HttpRule_Put:
					p.Method = "PUT"
					p.Name = method.Put
				case *options.HttpRule_Post:
					p.Method = "POST"
					p.Name = method.Post
				case *options.HttpRule_Delete:
					p.Method = "DELETE"
					p.Name = method.Delete
				case *options.HttpRule_Patch:
					p.Method = "PATCH"
					p.Name = method.Patch
				case *options.HttpRule_Custom:
					p.Method = "CUSTOM"
					p.Name = method.Custom.Path
				default:
					continue
				}
				ts.AddProtocol(p)
			}
		}
	}
}

func typeStrip(name string) string {
	return name[1:]
}

func (ts *GenData) SetCmd(cmdType *Enum) {
	ts.cmdType = cmdType
	ts.CmdTypeName = *ts.cmdType.Name
	ts.initCmdIdValue()
}

func (ts *GenData) getMethodCmdName(meth *descriptor.MethodDescriptorProto) (ret string) {
	if meth == nil {
		return ""
	}
	cmd := getMethodCmd(meth)
	if cmd == 0 {
		return ""
	}

	enum := ts.GetEnumValue(cmd)

	if enum == nil {
		ret = *meth.Name
	} else {
		if enum.Name == nil {
			return ""
		} else {
			ret = *(enum.Name)
		}
	}
	return
}

func (ts *GenData) getProtoItem(protoName string) *Protocol {
	return ts.MapProto[protoName]
}

// AddServices 设置 Services
func (ts *GenData) AddServices(services []*descriptor.ServiceDescriptorProto) {
	ts.Services = append(ts.Services, services...)
}

func (ts *GenData) initCmdIdValue() {
	value := ts.cmdType.EnumDescriptorProto.Value
	for _, item := range value {
		ts.cmdIdValue[*item.Number] = item
	}
}

func (ts *GenData) GetEnumValue(value int32) *descriptor.EnumValueDescriptorProto {
	ret, _ := ts.cmdIdValue[value]
	return ret
}

func (ts *GenData) AddProtocol(protocol *Protocol) {
	ts.MapProto[protocol.Name] = protocol
	ts.ProtoList = append(ts.ProtoList, protocol.Name)
}

func (ts *GenData) baseName(name string) string {
	arr := strings.Split(name, ".")
	length := len(arr)
	if len(arr) == 0 {
		return name
	}
	return arr[length-1]
}

func (ts *GenData) RenderGOCodeV1(name string) *plugin.CodeGeneratorResponse_File {
	tplStr := `// Code generated by protoc-gen-idmsg. DO NOT EDIT.

package pb

import "google.golang.org/protobuf/proto"

{{- $cmdTypeName := .CmdTypeName}}
var (
	MethodRouter = map[string]map[string]interface{}{
{{- range .Services}}
{{- $serviceName :=  .Name}}
{{- $serviceAuth := getServiceAuth .}}
	{{- range .Method}}
	{{- $methodCmdName := getMethodCmdName .}}
	{{- if ne $methodCmdName ""}}
	{{- $methodName := getMethodName .}}
	{{- $methRouter := getMethRouter .}}
	{{- $methodAuth := getMethodAuth $serviceAuth .}}
		"{{$methRouter.Route}}": {"http":{{$methRouter.Method}}, "method": "{{$serviceName}}.{{$methodName}}","auth": {{$methodAuth}},"cmd": int32({{$cmdTypeName}}_{{$methodCmdName}})},
	{{- end}}
	{{- end}}
{{- end}}
	}

	IdMsg = map[int32]struct {
		Req func() proto.Message
		Rsp func() proto.Message
		Auth int32
		Name string
		HTTP string
	}{
		{{- range .ProtoList}}
		{{- $protoItem := getProtoItem .}}
		int32({{$cmdTypeName}}_{{$protoItem.IdName}}): {Req: func() proto.Message { return &{{baseName $protoItem.ReqType}}{} }, Rsp: func() proto.Message { return &{{baseName $protoItem.RspType}}{} }, Auth: {{$protoItem.MthAuth}}, Name: "{{$protoItem.Name}}", HTTP: "{{$protoItem.Method}}"},
		{{- end}}
    }
)`

	str := ts.renderCode(tplStr)
	return &plugin.CodeGeneratorResponse_File{
		Name:    &name,
		Content: &str,
	}
}

func (ts *GenData) RenderGOCode(name string) *plugin.CodeGeneratorResponse_File {
	tplStr := `// Code generated by protoc-gen-idmsg. DO NOT EDIT.

package pb

import (
	"github.com/goodluck0107/gcore/gprotocol/interfaces"
	"google.golang.org/protobuf/proto"
)

{{- $cmdTypeName := .CmdTypeName}}

var (
	MethodRouter = map[string]*interfaces.MethodItem{
{{- range .Services}}
{{- $serviceName :=  .Name}}
{{- $serviceAuth := getServiceAuth .}}
	{{- range .Method}}
	{{- $methodCmdName := getMethodCmdName .}}
	{{- if ne $methodCmdName ""}}
    {{- $methodName := getMethodName .}}
	{{- $methRouter := getMethRouter .}}
	{{- $methodAuth := getMethodAuth $serviceAuth .}}
		"{{$methRouter.Route}}": {HTTP: "{{$methRouter.Method}}", Method: "{{$serviceName}}.{{$methodName}}", Auth: {{$methodAuth}}, Cmd: int32({{$cmdTypeName}}_{{$methodCmdName}})},
	{{- end}}
	{{- end}}
{{- end}}
	}

	IdMsg = map[int32]*interfaces.ReqItem{
		{{- range .ProtoList}}
		{{- $protoItem := getProtoItem .}}
		int32({{$cmdTypeName}}_{{$protoItem.IdName}}): {Req: func() proto.Message { return &{{baseName $protoItem.ReqType}}{} }, Rsp: func() proto.Message { return &{{baseName $protoItem.RspType}}{} }, Auth: {{$protoItem.MthAuth}}, Name: "{{$protoItem.Name}}", HTTP: "{{$protoItem.Method}}"},
		{{- end}}
    }
)
type MsgDefStruct struct{}

func (m *MsgDefStruct) GetMethodRouter() map[string]*interfaces.MethodItem {
	return MethodRouter
}
func (m *MsgDefStruct) GetIdMsg() map[int32]*interfaces.ReqItem {
	return IdMsg
}

var MsgDef = &MsgDefStruct{}`

	str := ts.renderCode(tplStr)

	return &plugin.CodeGeneratorResponse_File{
		Name:    &name,
		Content: &str,
	}
}

func (ts *GenData) RenderTSCode(name string) *plugin.CodeGeneratorResponse_File {
	tplStr := `// Code generated by protoc-gen-idmsg. DO NOT EDIT.

let Protocol = {
	{{- range .ProtoList}}
	{{- $protoItem := getProtoItem .}}
  {{$protoItem.IdName}}: "{{$protoItem.Name}}",
	{{- end}}
}

let ProtocolRouter 					= {} // 协议名:[code, req, rsp]
{{- range .ProtoList}}
{{- $protoItem := getProtoItem .}}
ProtocolRouter[Protocol.{{$protoItem.IdName}}] = [{{$protoItem.Code}}, "{{$protoItem.ReqType}}", "{{$protoItem.RspType}}"]
{{- end}}

let ProtocolId   = {} // 协议id: [router, req, rsp]
{{- range .ProtoList}}
{{- $protoItem := getProtoItem .}}
ProtocolId[{{$protoItem.Code}}] = [Protocol.{{$protoItem.IdName}}, "{{$protoItem.ReqType}}", "{{$protoItem.RspType}}"]
{{- end}}

export { Protocol, ProtocolRouter, ProtocolId }`
	str := ts.renderCode(tplStr)
	return &plugin.CodeGeneratorResponse_File{
		Name:    &name,
		Content: &str,
	}
}

// 使用模板渲染代码，未完成
func (ts *GenData) renderCode(tplStr string) string {
	tplFuncMap := map[string]interface{}{
		"getInnerSvcName":  getInnerSvcName,
		"getNativeSvcName": getNativeSvcName,
		"getGenSvcName":    getGenSvcName,
		"getInputType":     getInputType,
		"getOutputType":    getOutputType,
		"firstUpper":       firstUpper,
		"getServiceAuth":   getServiceAuth,
		"getMethodAuth":    getMethodAuth,
		"getMethodName":    getMethodName,
		"getMethRouter":    getMethRouter,

		"getMethodCmdName": ts.getMethodCmdName,
		"getProtoItem":     ts.getProtoItem,
		"baseName":         ts.baseName,
	}
	tpl := template.Must(template.New("tmpl").Funcs(tplFuncMap).Parse(tplStr))
	buf := &bytes.Buffer{}

	ts.Sort()

	sort.Slice(ts.Services, func(i, j int) bool {
		return *ts.Services[i].Name < *ts.Services[j].Name
	})

	for _, item := range ts.Services {
		sort.Slice(item.Method, func(i, j int) bool {
			return *item.Method[i].Name < *item.Method[j].Name
		})
	}

	if err := tpl.Execute(buf, ts); err != nil {
		fmt.Printf("template err:%v\n", err)
	}
	return buf.String()
}

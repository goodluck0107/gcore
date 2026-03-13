package genmsg

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/golang/glog"
	options "google.golang.org/genproto/googleapis/api/annotations"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

// ImportInfo 包导入信息
type ImportInfo struct {
	Path  string // 完整导入路径
	Alias string // 别名
}

type Protocol struct {
	Code       int32
	Method     string
	Name       string
	ReqType    string // 完整限定名或裸类型名
	RspType    string // 完整限定名或裸类型名
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
	ProtoList   []string
	cmdType     *Enum
	cmdIdValue  map[int32]*descriptor.EnumValueDescriptorProto
	Services    []*descriptor.ServiceDescriptorProto
	Imports     map[string]*ImportInfo // key: 完整导入路径
	importAlias map[string]int         // key: 别名，value: 计数

	// 生成文件所在的包信息
	genFileGoPkg    string // 如 "gitlab.yq-dev-inner.com/.../ck-proto/pb"
	genFileProtoPkg string // 如 "pb"

	// 缓存所有文件的 go_package 映射，用于判断同目录
	fileGoPackages map[string]string // key: protoPkg, value: goPkg
}

func NewGenData() *GenData {
	ret := &GenData{
		MapProto:       map[string]*Protocol{},
		cmdIdValue:     map[int32]*descriptor.EnumValueDescriptorProto{},
		Imports:        make(map[string]*ImportInfo),
		importAlias:    make(map[string]int),
		fileGoPackages: make(map[string]string),
	}
	return ret
}

func (ts *GenData) Sort() {
	sort.Slice(ts.ProtoList, func(i, j int) bool {
		return ts.MapProto[ts.ProtoList[i]].Code < ts.MapProto[ts.ProtoList[j]].Code
	})
}

// initGenFilePackage 从文件列表中初始化生成文件所在的包信息
func (ts *GenData) initGenFilePackage(files []*descriptor.FileDescriptorProto, enumType *Enum) {
	var targetFile *descriptor.FileDescriptorProto

	if enumType != nil && enumType.File != nil {
		targetFile = enumType.File.FileDescriptorProto
	} else if len(files) > 0 {
		targetFile = files[0]
	}

	if targetFile == nil {
		return
	}

	if targetFile.Package != nil {
		ts.genFileProtoPkg = *targetFile.Package
	}

	if targetFile.Options != nil && targetFile.Options.GoPackage != nil {
		goPkg := *targetFile.Options.GoPackage
		if idx := strings.Index(goPkg, ";"); idx != -1 {
			ts.genFileGoPkg = goPkg[:idx]
		} else {
			ts.genFileGoPkg = goPkg
		}
	}

	// 缓存所有文件的 go_package
	for _, f := range files {
		if f.Package == nil || f.Options == nil || f.Options.GoPackage == nil {
			continue
		}
		goPkg := *f.Options.GoPackage
		if idx := strings.Index(goPkg, ";"); idx != -1 {
			goPkg = goPkg[:idx]
		}
		ts.fileGoPackages[*f.Package] = goPkg
	}
}

// getGoPackage 获取指定 proto 包的 go_package
func (ts *GenData) getGoPackage(protoPkg string) string {
	if goPkg, ok := ts.fileGoPackages[protoPkg]; ok {
		return goPkg
	}
	return ""
}

// isSameDirectory 判断目标包是否与当前文件在同一目录
func (ts *GenData) isSameDirectory(targetProtoPkg string) bool {
	// 如果是当前包，肯定同目录
	if targetProtoPkg == ts.genFileProtoPkg {
		return true
	}

	// 获取目标包的 go_package
	targetGoPkg := ts.getGoPackage(targetProtoPkg)
	if targetGoPkg == "" {
		return false
	}

	// 比较目录部分（去掉最后的包名）
	genDir := filepath.Dir(ts.genFileGoPkg)
	targetDir := filepath.Dir(targetGoPkg)

	// 如果目录相同，则在同一目录
	return genDir == targetDir
}

// getOrCreateImport 获取或创建包导入信息
// 如果同目录，返回 nil
func (ts *GenData) getOrCreateImport(protoPkg string) *ImportInfo {
	// 如果是当前包或同目录，不需要导入
	if protoPkg == ts.genFileProtoPkg || ts.isSameDirectory(protoPkg) {
		return nil
	}

	// 构造完整导入路径
	goImportPath := ts.buildImportPath(protoPkg)

	if imp, ok := ts.Imports[goImportPath]; ok {
		return imp
	}

	parts := strings.Split(protoPkg, ".")
	alias := parts[len(parts)-1]

	originalAlias := alias
	count := ts.importAlias[alias]
	if count > 0 {
		alias = fmt.Sprintf("%s%d", originalAlias, count)
	}
	ts.importAlias[originalAlias] = count + 1

	imp := &ImportInfo{
		Path:  goImportPath,
		Alias: alias,
	}
	ts.Imports[goImportPath] = imp
	return imp
}

// buildImportPath 根据 proto 包名构造完整 Go 导入路径
func (ts *GenData) buildImportPath(protoPkg string) string {
	if ts.genFileGoPkg == "" {
		return strings.ReplaceAll(protoPkg, ".", "/")
	}

	genPkgParts := strings.Split(ts.genFileProtoPkg, ".")
	targetPkgParts := strings.Split(protoPkg, ".")

	commonLen := 0
	for i := 0; i < len(genPkgParts) && i < len(targetPkgParts); i++ {
		if genPkgParts[i] == targetPkgParts[i] {
			commonLen++
		} else {
			break
		}
	}

	goPkgParts := strings.Split(ts.genFileGoPkg, "/")
	lastProtoPart := genPkgParts[len(genPkgParts)-1]

	if len(goPkgParts) > 0 && goPkgParts[len(goPkgParts)-1] == lastProtoPart {
		basePath := strings.Join(goPkgParts[:len(goPkgParts)-1], "/")
		targetPath := strings.Join(targetPkgParts, "/")
		return basePath + "/" + targetPath
	}

	return ts.genFileGoPkg + "/" + strings.Join(targetPkgParts[commonLen:], "/")
}

// baseName 获取类型名（裸名）
func (ts *GenData) baseName(fullType string) string {
	cleanType := strings.TrimPrefix(fullType, ".")
	lastDot := strings.LastIndex(cleanType, ".")
	if lastDot == -1 {
		return cleanType
	}
	return cleanType[lastDot+1:]
}

// parseType 解析类型
// 同目录返回裸类型名，不同目录返回带包别名的类型名
func (ts *GenData) parseType(fullType string, currentPkg string) (string, *ImportInfo) {
	cleanType := strings.TrimPrefix(fullType, ".")
	lastDot := strings.LastIndex(cleanType, ".")
	if lastDot == -1 {
		return cleanType, nil
	}

	protoPkg := cleanType[:lastDot]
	typeName := cleanType[lastDot+1:]

	// 如果是当前包，直接返回裸名
	if protoPkg == currentPkg {
		return typeName, nil
	}

	// 判断是否同目录
	if ts.isSameDirectory(protoPkg) {
		// 同目录，使用裸类型名
		return typeName, nil
	}

	// 不同目录，需要导入
	imp := ts.getOrCreateImport(protoPkg)
	if imp == nil {
		return typeName, nil
	}
	return imp.Alias + "." + typeName, imp
}

// FindEnumType 找到定义 cmd 的枚举类型
func (ts *GenData) FindEnumType(reg *Registry) *Enum {
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

	fileDesc := make([]*descriptor.FileDescriptorProto, 0, len(files))
	for _, file := range files {
		fileDesc = append(fileDesc, file.FileDescriptorProto)
	}

	ts.initGenFilePackage(fileDesc, enumType)

	for _, file := range files {
		currentPkg := ts.getCurrentPackage(file.FileDescriptorProto)

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
				if methAuth != 0 {
					auth = int32(methAuth)
				}
				if auth == 0 && serviceAuth != 0 {
					auth = int32(serviceAuth)
				}

				reqTypeStr, _ := ts.parseType(*meth.InputType, currentPkg)
				rspTypeStr, _ := ts.parseType(*meth.OutputType, currentPkg)

				p := &Protocol{
					Code:    int32(optionId.Cmd),
					ReqType: reqTypeStr,
					RspType: rspTypeStr,
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

// getCurrentPackage 获取当前文件的 proto 包名
func (ts *GenData) getCurrentPackage(file *descriptor.FileDescriptorProto) string {
	if file.Package != nil {
		return *file.Package
	}
	return ""
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

// RenderGOCode 生成 Go 代码（带动态导入）
func (ts *GenData) RenderGOCode(name string) *plugin.CodeGeneratorResponse_File {
	tplStr := `// Code generated by protoc-gen-idmsg. DO NOT EDIT.

package pb

import (
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/goodluck0107/gmsgdef"
	"google.golang.org/protobuf/proto"
{{- range .Imports}}
	{{.Alias}} "{{.Path}}"
{{- end}}
)

type StringValue = wrappers.StringValue
type BoolValue = wrappers.BoolValue
type BytesValue = wrappers.BytesValue
type FloatValue = wrappers.FloatValue
type DoubleValue = wrappers.DoubleValue
type Int32Value = wrappers.Int32Value
type Int64Value = wrappers.Int64Value
type UInt32Value = wrappers.UInt32Value
type UInt64Value = wrappers.UInt64Value

{{- $cmdTypeName := .CmdTypeName}}

var (
	MethodRouter = map[string]*gmsgdef.MethodItem{
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

	IdMsg = map[int32]*gmsgdef.ReqItem{
		{{- range .ProtoList}}
		{{- $protoItem := getProtoItem .}}
		int32({{$cmdTypeName}}_{{$protoItem.IdName}}): {Req: func() proto.Message { return &{{$protoItem.ReqType}}{} }, Rsp: func() proto.Message { return &{{$protoItem.RspType}}{} }, Auth: {{$protoItem.MthAuth}}, Name: "{{$protoItem.Name}}", HTTP: "{{$protoItem.Method}}", RPCMethod: {{$protoItem.SvcName}}_{{$protoItem.MthName}}_FullMethodName },
		{{- end}}
    }
)
type MsgDefStruct struct{}

func (m *MsgDefStruct) GetMethodRouter() map[string]*gmsgdef.MethodItem {
	return MethodRouter
}
func (m *MsgDefStruct) GetIdMsg() map[int32]*gmsgdef.ReqItem {
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

// 使用模板渲染代码
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

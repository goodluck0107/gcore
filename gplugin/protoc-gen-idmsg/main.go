package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"protoc-gen-idmsg/genmsg"
	_ "protoc-gen-idmsg/pb"
	"strings"

	plugin "google.golang.org/protobuf/types/pluginpb"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/codegenerator"
	"google.golang.org/protobuf/proto"
)

var (
	importPrefix               = flag.String("import_prefix", "", "prefix to be added to go package paths for imported proto files")
	file                       = flag.String("file", "-", "where to load data from")
	allowDeleteBody            = flag.Bool("allow_delete_body", false, "unless set, HTTP DELETE methods may not have a body")
	grpcAPIConfiguration       = flag.String("grpc_api_configuration", "", "path to gRPC API Configuration in YAML format")
	allowMerge                 = flag.Bool("allow_merge", false, "if set, generation one swagger file out of multiple protos")
	allowSave                  = flag.Bool("allow_save", false, "if set, save req.bin for debug")
	mergeFileName              = flag.String("merge_file_name", "apidocs", "target swagger file name prefix after merge")
	useJSONNamesForFields      = flag.Bool("json_names_for_fields", false, "if it sets Field.GetJsonName() will be used for generating swagger definitions, otherwise Field.GetName() will be used")
	repeatedPathParamSeparator = flag.String("repeated_path_param_separator", "csv", "configures how repeated fields should be split. Allowed values are `csv`, `pipes`, `ssv` and `tsv`.")
	versionFlag                = flag.Bool("version", false, "print the current verison")
	allowRepeatedFieldsInBody  = flag.Bool("allow_repeated_fields_in_body", false, "allows to use repeated field in `body` and `response_body` field of `google.api.http` annotation option")
	includePackageInTags       = flag.Bool("include_package_in_tags", false, "if unset, the gRPC service name is added to the `Tags` field of each operation. if set and the `package` directive is shown in the proto file, the package name will be prepended to the service name")
	useFQNForSwaggerName       = flag.Bool("fqn_for_swagger_name", false, "if set, the object's swagger names will use the fully qualify name from the proto definition (ie my.package.MyMessage.MyInnerMessage")
	v1                         = flag.Bool("v1", false, "if set, output file use old version")
)

// Variables set by goreleaser at build time
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func saveFile(filepath string, content []byte) {
	err := ioutil.WriteFile(filepath, content, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if *versionFlag {
		fmt.Printf("Version %v, commit %v, built at %v\n", version, commit, date)
		os.Exit(0)
	}

	glog.V(1).Info("Processing code generator request")
	f := os.Stdin
	if *file != "-" {
		var err error
		f, err = os.Open(*file)
		if err != nil {
			glog.Info(err)
		}
	}

	glog.V(1).Info("Parsing code generator request")

	mBytes, err := ioutil.ReadAll(f)

	if err != nil {
		glog.Fatal(err)
	}

	req, err := codegenerator.ParseRequest(ioutil.NopCloser(bytes.NewBuffer(mBytes)))
	pkgMap := make(map[string]string)

	var param string
	if req.Parameter != nil {
		param = *req.Parameter
	}

	if param != "" {
		err := parseReqParam(param, flag.CommandLine, pkgMap)
		if err != nil {
			glog.Fatalf("Error parsing flags: %v, %s", err, param)
		}
	}

	if *allowSave {
		saveFile("./req.bin", mBytes)
	}

	reg := genmsg.NewRegistry()
	err = reg.Load(req)

	if err != nil {
		glog.Fatal(err)
		return
	}

	genData := genmsg.NewGenData()
	genData.FromReg(reg)

	rpcData := genmsg.NewRpcData()
	rpcData.FromReg(reg)

	emitFile(genData.RenderTSCode("msg.ts"))

	if *v1 {
		emitFile(genData.RenderGOCodeV1("msg.go"))
	} else {
		emitFile(genData.RenderGOCode("msg.go"))
		emitFiles(rpcData.GenRpcV2())
	}
}

func emitFile(out *plugin.CodeGeneratorResponse_File) {
	emitResp(&plugin.CodeGeneratorResponse{File: []*plugin.CodeGeneratorResponse_File{out}})
}

func emitFiles(out []*plugin.CodeGeneratorResponse_File) {
	emitResp(&plugin.CodeGeneratorResponse{File: out})
}

func emitError(err error) {
	emitResp(&plugin.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

func emitResp(resp *plugin.CodeGeneratorResponse) {
	buf, err := proto.Marshal(resp)
	if err != nil {
		glog.Fatal(err)
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		glog.Fatal(err)
	}
}

// parseReqParam parses a CodeGeneratorRequest parameter and adds the
// extracted values to the given FlagSet and pkgMap. Returns a non-nil
// error if setting a flag failed.
func parseReqParam(param string, f *flag.FlagSet, pkgMap map[string]string) error {
	if param == "" {
		return nil
	}
	for _, p := range strings.Split(param, ",") {
		spec := strings.SplitN(p, "=", 2)
		if len(spec) == 1 {
			if spec[0] == "allow_delete_body" {
				err := f.Set(spec[0], "true")
				if err != nil {
					return fmt.Errorf("Cannot set flag %s: %v", p, err)
				}
				continue
			}
			if spec[0] == "allow_merge" {
				err := f.Set(spec[0], "true")
				if err != nil {
					return fmt.Errorf("Cannot set flag %s: %v", p, err)
				}
				continue
			}
			if spec[0] == "v1" {
				err := f.Set(spec[0], "true")
				if err != nil {
					return fmt.Errorf("Cannot set flag %s: %v", p, err)
				}
				continue
			}
			if spec[0] == "v1" {
				err := f.Set(spec[0], "true")
				if err != nil {
					return fmt.Errorf("Cannot set flag %s: %v", p, err)
				}
				continue
			}
			if spec[0] == "allow_repeated_fields_in_body" {
				err := f.Set(spec[0], "true")
				if err != nil {
					return fmt.Errorf("Cannot set flag %s: %v", p, err)
				}
				continue
			}
			if spec[0] == "include_package_in_tags" {
				err := f.Set(spec[0], "true")
				if err != nil {
					return fmt.Errorf("Cannot set flag %s: %v", p, err)
				}
				continue
			}
			err := f.Set(spec[0], "")
			if err != nil {
				return fmt.Errorf("Cannot set flag %s: %v", p, err)
			}
			continue
		}
		name, value := spec[0], spec[1]
		if strings.HasPrefix(name, "M") {
			pkgMap[name[1:]] = value
			continue
		}
		if err := f.Set(name, value); err != nil {
			return fmt.Errorf("Cannot set flag %s: %v", p, err)
		}
	}
	return nil
}

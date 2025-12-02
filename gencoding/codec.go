package gencoding

import (
	"github.com/goodluck0107/gcore/gencoding/json"
	"github.com/goodluck0107/gcore/gencoding/msgpack"
	"github.com/goodluck0107/gcore/gencoding/proto"
	"github.com/goodluck0107/gcore/gencoding/toml"
	"github.com/goodluck0107/gcore/gencoding/xml"
	"github.com/goodluck0107/gcore/gencoding/yaml"
	"github.com/goodluck0107/gcore/glog"
)

var codecs = make(map[string]Codec)

func init() {
	Register(json.DefaultCodec)
	Register(proto.DefaultCodec)
	Register(toml.DefaultCodec)
	Register(xml.DefaultCodec)
	Register(yaml.DefaultCodec)
	Register(msgpack.DefaultCodec)
}

type Codec interface {
	// Name 编解码器类型
	Name() string
	// Marshal 编码
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal 解码
	Unmarshal(data []byte, v interface{}) error
}

// Register 注册编解码器
func Register(codec Codec) {
	if codec == nil {
		glog.Fatal("can't register a invalid codec")
	}

	name := codec.Name()

	if name == "" {
		glog.Fatal("can't register a codec without name")
	}

	if _, ok := codecs[name]; ok {
		glog.Warnf("the old %s codec will be overwritten", name)
	}

	codecs[name] = codec
}

// Invoke 调用编解码器
func Invoke(name string) Codec {
	codec, ok := codecs[name]
	if !ok {
		glog.Fatalf("%s codec is not registered", name)
	}

	return codec
}

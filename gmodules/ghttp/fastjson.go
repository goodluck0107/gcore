package ghttp

import (
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ignoreOmitEmptyTagExtension struct {
	jsoniter.DummyExtension
}

type ignoreOmitEmptyTagEncoder struct {
	originDecoder jsoniter.ValEncoder
}

func (p *ignoreOmitEmptyTagEncoder) IsEmpty(ptr unsafe.Pointer) bool { //关键逻辑
	return false
}

func (p *ignoreOmitEmptyTagEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	p.originDecoder.Encode(ptr, stream)
}

func (e *ignoreOmitEmptyTagExtension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	return &ignoreOmitEmptyTagEncoder{encoder}
}

func init() {
	jsoniter.RegisterExtension(&ignoreOmitEmptyTagExtension{})
}

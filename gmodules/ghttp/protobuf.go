package ghttp

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang/protobuf/proto"
)

type ProtobufEncoder struct {
}

func (p ProtobufEncoder) Name() string {
	return "protobuf"
}

func (p ProtobufEncoder) MIMETypes() []string {
	return []string{"application/protobuf", "application/x-protobuf", "application/vnd.google.protobuf"}
}

func (p ProtobufEncoder) Parse(c fiber.Ctx, out any) error {
	if msg, ok := out.(proto.Message); ok {
		return proto.Unmarshal(c.Body(), msg)
	}
	return nil
}

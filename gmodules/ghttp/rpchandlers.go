package ghttp

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/goodluck0107/gcore/gcodes"
	"github.com/goodluck0107/gcore/glog"
	ghandler "github.com/goodluck0107/gcore/gprotocol/handler"
)

const (
	LocalKeyHandlerMetadata = "ghttp-handler-metadata"
)

type HandlersContext interface {
	// CMD 获取CMD
	CMD() int32
	// AuthType 获取认证类型
	AuthType() int32
	// HandlerMetadata 获取HandlerMetadata
	HandlerMetadata() ghandler.Metadata
}

func handlersMiddleware(md ghandler.Metadata) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		ctx.Locals(LocalKeyHandlerMetadata, md)
		return ctx.Next()
	}
}

func convertHandle(handle ghandler.Handler) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		ctxWrapper := &context{Ctx: ctx}
		dec := func(req interface{}) error {
			var err error
			switch ctx.Method() {
			case fiber.MethodGet:
				err = ctx.Bind().Query(req)
			case fiber.MethodPut, fiber.MethodPost:
				err = ctx.Bind().Body(req)
			default:
				err = fmt.Errorf("not support method: %s", ctx.Method())
			}
			return err
		}
		ret, code := handle(ctxWrapper.Context(), dec)
		if code != gcodes.OK {
			if redirect := code.Redirect(); len(redirect) > 0 {
				glog.Warnf("[redirect] [%s] %s -> %s", ctx.IP(), ctx.OriginalURL(), redirect)
				return ctxWrapper.Redirect().Status(fiber.StatusTemporaryRedirect).To(redirect)
			}
			return ctxWrapper.Failure(code)
		} else {
			return ctxWrapper.Success(ret)
		}
	}
}

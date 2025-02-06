package ghttp

import (
	"bytes"
	"gitee.com/monobytes/gcore/gcodes"
	"github.com/gofiber/fiber/v3"
	"io"
	"net/http"
	"net/url"
)

type Resp struct {
	Code    int    `json:"code"`           // 响应码
	Message string `json:"message"`        // 响应消息
	Data    any    `json:"data,omitempty"` // 响应数据
}

type Context interface {
	fiber.Ctx
	// CTX 获取fiber.Ctx
	CTX() fiber.Ctx
	// Proxy 获取代理API
	Proxy() *Proxy
	// Failure 失败响应
	Failure(rst any) error
	// Success 成功响应
	Success(data ...any) error
	// StdRequest 获取标准请求（net/http）
	StdRequest() *http.Request
}

type context struct {
	fiber.Ctx
	proxy *Proxy
}

// CTX 获取fiber.Ctx
func (c *context) CTX() fiber.Ctx {
	return c.Ctx
}

// Proxy 代理API
func (c *context) Proxy() *Proxy {
	return c.proxy
}

// Failure 失败响应
func (c *context) Failure(rst any) error {
	switch v := rst.(type) {
	case error:
		code, _ := gcodes.Convert(v)

		return c.JSON(&Resp{Code: code.Code(), Message: code.Message()})
	case *gcodes.Code:
		return c.JSON(&Resp{Code: v.Code(), Message: v.Message()})
	default:
		return c.JSON(&Resp{Code: gcodes.Unknown.Code(), Message: gcodes.Unknown.Message()})
	}
}

// Success 成功响应
func (c *context) Success(data ...any) error {
	if len(data) > 0 {
		return c.JSON(&Resp{Code: gcodes.OK.Code(), Message: gcodes.OK.Message(), Data: data[0]})
	} else {
		return c.JSON(&Resp{Code: gcodes.OK.Code(), Message: gcodes.OK.Message()})
	}
}

// StdRequest 获取标准请求（net/http）
func (c *context) StdRequest() *http.Request {
	req := c.Request()

	std := &http.Request{}
	std.Method = c.Method()
	std.URL, _ = url.Parse(req.URI().String())
	std.Proto = c.Protocol()
	std.ProtoMajor, std.ProtoMinor, _ = http.ParseHTTPVersion(std.Proto)
	std.Header = c.GetReqHeaders()
	std.Host = c.Host()
	std.ContentLength = int64(len(c.Body()))
	std.RemoteAddr = c.Context().RemoteAddr().String()
	std.RequestURI = string(req.RequestURI())

	if req.Body() != nil {
		std.Body = io.NopCloser(bytes.NewReader(req.Body()))
	}

	return std
}

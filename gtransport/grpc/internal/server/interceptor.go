package server

import (
	"context"
	"gitee.com/monobytes/gcore/glog"
	"google.golang.org/grpc"
	"runtime"
)

func recoverInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				glog.Panic(err)
			default:
				glog.Panicf("panic error: %v", err)
			}
		}
	}()

	return handler(ctx, req)
}

package service

import (
	"context"
	"gitee.com/monobytes/gcore/examples/protocol/pb"
)

type Greeter struct {
	pb.UnimplementedGreeterServer
}

func NewGreeter() *Greeter {
	return &Greeter{}
}

func (g *Greeter) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloRsp, error) {
	return &pb.HelloRsp{Msg: "hello," + req.Name}, nil
}

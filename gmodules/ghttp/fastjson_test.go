package ghttp

import (
	"gitee.com/monobytes/gcore/examples/protocol/pb"
	"testing"
)

func TestFastJson(t *testing.T) {
	user := &pb.HelloReq{}
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

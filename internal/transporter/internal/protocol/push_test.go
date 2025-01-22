package protocol_test

import (
	"gitee.com/monobytes/gcore/gpacket"
	"gitee.com/monobytes/gcore/gsession"
	"gitee.com/monobytes/gcore/gwrap/buffer"
	"gitee.com/monobytes/gcore/internal/transporter/internal/codes"
	"gitee.com/monobytes/gcore/internal/transporter/internal/protocol"
	"testing"
)

func TestEncodePushReq(t *testing.T) {
	message, err := gpacket.PackMessage(&gpacket.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodePushReq(1, gsession.User, 3, buffer.NewNocopyBuffer(message))

	t.Log(buf.Bytes())
}

func TestDecodePushReq(t *testing.T) {
	message, err := gpacket.PackMessage(&gpacket.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodePushReq(1, gsession.User, 3, buffer.NewNocopyBuffer(message))

	seq, kind, target, msg, err := protocol.DecodePushReq(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("seq: %v", seq)
	t.Logf("kind: %v", kind)
	t.Logf("target: %v", target)
	t.Logf("message: %v", msg)
}

func TestEncodePushRes(t *testing.T) {
	buffer := protocol.EncodePushRes(1, codes.OK)

	t.Log(buffer.Bytes())
}

func TestDecodePushRes(t *testing.T) {
	buffer := protocol.EncodePushRes(1, codes.OK)

	code, err := protocol.DecodePushRes(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("code: %v", code)
}

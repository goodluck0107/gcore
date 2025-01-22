package protocol_test

import (
	"gitee.com/monobytes/gcore/gpacket"
	"gitee.com/monobytes/gcore/gsession"
	"gitee.com/monobytes/gcore/gwrap/buffer"
	"gitee.com/monobytes/gcore/internal/transporter/internal/codes"
	"gitee.com/monobytes/gcore/internal/transporter/internal/protocol"
	"testing"
)

func TestEncodeBroadcastReq(t *testing.T) {
	message, err := gpacket.PackMessage(&gpacket.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodeBroadcastReq(1, gsession.User, buffer.NewNocopyBuffer(message))

	t.Log(buf.Bytes())
}

func TestDecodeBroadcastReq(t *testing.T) {
	message, err := gpacket.PackMessage(&gpacket.Message{
		Route:  1,
		Seq:    2,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	buf := protocol.EncodeBroadcastReq(1, gsession.User, buffer.NewNocopyBuffer(message))

	seq, kind, message, err := protocol.DecodeBroadcastReq(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("seq: %v", seq)
	t.Logf("kind: %v", kind)
	t.Logf("message: %v", string(message))
}

func TestEncodeBroadcastRes(t *testing.T) {
	buf := protocol.EncodeBroadcastRes(1, codes.OK, 20)

	t.Log(buf.Bytes())
}

func TestDecodeBroadcastRes(t *testing.T) {
	buf := protocol.EncodeBroadcastRes(1, codes.OK, 20)

	code, total, err := protocol.DecodeBroadcastRes(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("code: %v", code)
	t.Logf("total: %v", total)
}

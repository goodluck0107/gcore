package protocol_test

import (
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gutils/guuid"
	"github.com/goodluck0107/gcore/internal/transporter/internal/codes"
	"github.com/goodluck0107/gcore/internal/transporter/internal/protocol"
	"testing"
)

func TestEncodeHandshakeReq(t *testing.T) {
	buffer := protocol.EncodeHandshakeReq(1, gcluster.Gate, guuid.UUID())

	t.Log(buffer.Bytes())
}

func TestDecodeHandshakeReq(t *testing.T) {
	buffer := protocol.EncodeHandshakeReq(1, gcluster.Gate, guuid.UUID())

	seq, insKind, insID, err := protocol.DecodeHandshakeReq(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("seq: %v", seq)
	t.Logf("kind: %v", insKind)
	t.Logf("id: %v", insID)
}

func TestEncodeHandshakeRes(t *testing.T) {
	buffer := protocol.EncodeHandshakeRes(1, codes.OK)

	t.Log(buffer.Bytes())
}

func TestDecodeHandshakeRes(t *testing.T) {
	buffer := protocol.EncodeHandshakeRes(1, codes.OK)

	code, err := protocol.DecodeHandshakeRes(buffer.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("code: %v", code)
}

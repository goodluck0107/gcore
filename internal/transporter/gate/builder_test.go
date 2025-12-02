package gate_test

import (
	"context"
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gsession"
	"github.com/goodluck0107/gcore/gutils/guuid"
	"github.com/goodluck0107/gcore/internal/transporter/gate"
	"testing"
)

func TestBuilder(t *testing.T) {
	builder := gate.NewBuilder(&gate.Options{
		InsID:   guuid.UUID(),
		InsKind: gcluster.Node,
	})

	client, err := builder.Build("127.0.0.1:49899")
	if err != nil {
		t.Fatal(err)
	}

	ip, miss, err := client.GetIP(context.Background(), gsession.User, 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("miss: %v ip: %v", miss, ip)

	ip, miss, err = client.GetIP(context.Background(), gsession.User, 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("miss: %v ip: %v", miss, ip)
}

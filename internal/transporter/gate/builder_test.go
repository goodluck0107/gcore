package gate_test

import (
	"context"
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/gsession"
	"gitee.com/monobytes/gcore/gutils/guuid"
	"gitee.com/monobytes/gcore/internal/transporter/gate"
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

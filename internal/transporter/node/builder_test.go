package node_test

import (
	"context"
	"github.com/goodluck0107/gcore/gcluster"
	"github.com/goodluck0107/gcore/gutils/guuid"
	"github.com/goodluck0107/gcore/internal/transporter/node"
	"testing"
)

func TestBuilder(t *testing.T) {
	builder := node.NewBuilder(&node.Options{
		InsID:   guuid.UUID(),
		InsKind: gcluster.Gate,
	})

	client, err := builder.Build("127.0.0.1:49898")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Deliver(context.Background(), 1, 2, []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
}

package endpoint_test

import (
	"github.com/goodluck0107/gcore/gwrap/endpoint"
	"testing"
)

func TestNewEndpoint(t *testing.T) {
	e := endpoint.NewEndpoint("grpc", "127.0.0.1:8080", false)

	t.Log(e.String())

	ee, err := endpoint.ParseEndpoint(e.String())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ee.Address())
}

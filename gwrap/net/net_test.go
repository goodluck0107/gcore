package net_test

import (
	"github.com/goodluck0107/gcore/gwrap/net"
	"testing"
)

func TestParseAddr(t *testing.T) {
	listenAddr, exposeAddr, err := net.ParseAddr(":8889")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(listenAddr, exposeAddr)
}

func TestInternalIP(t *testing.T) {
	ip, err := net.InternalIP()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ip)
}

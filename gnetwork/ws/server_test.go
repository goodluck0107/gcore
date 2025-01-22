package ws_test

import (
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gnetwork"
	"gitee.com/monobytes/gcore/gnetwork/ws"
	"gitee.com/monobytes/gcore/gpacket"
	"gitee.com/monobytes/gcore/gutils/gcall"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	server := ws.NewServer()
	server.OnStart(func() {
		t.Logf("server is started")
	})
	server.OnConnect(func(conn gnetwork.Conn) {
		t.Logf("connection is opened, connection id: %d", conn.ID())
	})
	server.OnDisconnect(func(conn gnetwork.Conn) {
		t.Logf("connection is closed, connection id: %d", conn.ID())
	})
	server.OnReceive(func(conn gnetwork.Conn, msg []byte) {
		message, err := gpacket.UnpackMessage(msg)
		if err != nil {
			t.Error(err)
			return
		}

		t.Logf("receive msg from client, connection id: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))

		msg, err = gpacket.PackMessage(&gpacket.Message{
			Seq:    1,
			Route:  1,
			Buffer: []byte("I'm fine~~"),
		})
		if err != nil {
			t.Fatal(err)
		}

		if err = conn.Push(msg); err != nil {
			t.Error(err)
		}
	})
	server.OnUpgrade(func(w http.ResponseWriter, r *http.Request) (allowed bool) {
		return true
	})

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	gcall.Go(func() {
		err := http.ListenAndServe(":8089", nil)
		if err != nil {
			glog.Errorf("mpprof server start failed: %v", err)
		}
	})

	select {}
}

func TestServer_Benchmark(t *testing.T) {
	server := ws.NewServer()
	server.OnStart(func() {
		t.Logf("server is started")
	})
	server.OnReceive(func(conn gnetwork.Conn, msg []byte) {
		_, err := gpacket.UnpackMessage(msg)
		if err != nil {
			t.Error(err)
			return
		}

		msg, err = gpacket.PackMessage(&gpacket.Message{
			Seq:    1,
			Route:  1,
			Buffer: []byte("I'm fine~~"),
		})
		if err != nil {
			t.Fatal(err)
		}

		if err = conn.Push(msg); err != nil {
			t.Error(err)
		}
	})
	server.OnUpgrade(func(w http.ResponseWriter, r *http.Request) (allowed bool) {
		return true
	})

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	gcall.Go(func() {
		err := http.ListenAndServe(":8089", nil)
		if err != nil {
			glog.Errorf("mpprof server start failed: %v", err)
		}
	})

	select {}
}

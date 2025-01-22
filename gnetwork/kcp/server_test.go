package kcp_test

import (
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gnetwork"
	"gitee.com/monobytes/gcore/gnetwork/kcp"
	"gitee.com/monobytes/gcore/gpacket"
	"net/http"
	_ "net/http/pprof"
	"testing"
)

func TestServer_Simple(t *testing.T) {
	server := kcp.NewServer()

	server.OnStart(func() {
		glog.Info("server is started")
	})

	server.OnStop(func() {
		glog.Info("server is stopped")
	})

	server.OnConnect(func(conn gnetwork.Conn) {
		glog.Infof("connection is opened, connection id: %d", conn.ID())
	})

	server.OnDisconnect(func(conn gnetwork.Conn) {
		glog.Infof("connection is closed, connection id: %d", conn.ID())
	})

	server.OnReceive(func(conn gnetwork.Conn, msg []byte) {
		message, err := gpacket.UnpackMessage(msg)
		if err != nil {
			glog.Errorf("unpack message failed: %v", err)
			return
		}

		glog.Infof("receive message from client, cid: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))

		msg, err = gpacket.PackMessage(&gpacket.Message{
			Seq:    1,
			Route:  1,
			Buffer: []byte("I'm fine~~"),
		})
		if err != nil {
			glog.Errorf("pack message failed: %v", err)
			return
		}

		if err = conn.Push(msg); err != nil {
			glog.Errorf("push message failed: %v", err)
		}
	})

	if err := server.Start(); err != nil {
		glog.Fatalf("start server failed: %v", err)
	}

	select {}
}

func TestServer_Benchmark(t *testing.T) {
	server := kcp.NewServer(
		kcp.WithServerHeartbeatInterval(0),
	)

	server.OnStart(func() {
		glog.Info("server is started")
	})

	server.OnReceive(func(conn gnetwork.Conn, msg []byte) {
		message, err := gpacket.UnpackMessage(msg)
		if err != nil {
			glog.Errorf("unpack message failed: %v", err)
			return
		}

		data, err := gpacket.PackMessage(&gpacket.Message{
			Seq:    message.Seq,
			Route:  message.Route,
			Buffer: message.Buffer,
		})
		if err != nil {
			glog.Errorf("pack message failed: %v", err)
			return
		}

		if err = conn.Send(data); err != nil {
			glog.Errorf("push message failed: %v", err)
			return
		}
	})

	if err := server.Start(); err != nil {
		glog.Fatalf("start server failed: %v", err)
	}

	go func() {
		err := http.ListenAndServe(":8089", nil)
		if err != nil {
			glog.Errorf("mpprof server start failed: %v", err)
		}
	}()

	select {}
}

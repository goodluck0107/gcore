package ws_test

import (
	"fmt"
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/gnetwork"
	"github.com/goodluck0107/gcore/gnetwork/ws"
	"github.com/goodluck0107/gcore/gpacket"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClient_Dial(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)

		go func() {
			client := ws.NewClient()
			client.OnConnect(func(conn gnetwork.Conn) {
				t.Log("connection is opened")
			})
			client.OnDisconnect(func(conn gnetwork.Conn) {
				t.Log("connection is closed")
			})
			client.OnReceive(func(conn gnetwork.Conn, msg []byte) {
				message, err := gpacket.UnpackMessage(msg)
				if err != nil {
					t.Error(err)
					return
				}

				t.Logf("receive msg from server, connection id: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))
			})

			defer wg.Done()

			conn, err := client.Dial()
			if err != nil {
				t.Fatal(err)
			}

			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			defer conn.Close()

			times := 0
			msg, _ := gpacket.PackMessage(&gpacket.Message{
				Seq:    1,
				Route:  1,
				Buffer: []byte("hello server~~"),
			})

			for {
				select {
				case <-ticker.C:
					if err = conn.Push(msg); err != nil {
						t.Error(err)
						return
					}

					times++

					if times >= 5 {
						return
					}
				}
			}
		}()
	}

	wg.Wait()
}

func TestNewClient(t *testing.T) {
	client := ws.NewClient()

	client.OnConnect(func(conn gnetwork.Conn) {
		glog.Info("connection is opened")
	})
	client.OnDisconnect(func(conn gnetwork.Conn) {
		glog.Info("connection is closed")
	})
	client.OnReceive(func(conn gnetwork.Conn, msg []byte) {
		message, err := gpacket.UnpackMessage(msg)
		if err != nil {
			t.Error(err)
			return
		}

		t.Logf("receive msg from server, connection id: %d, seq: %d, route: %d, msg: %s", conn.ID(), message.Seq, message.Route, string(message.Buffer))
	})

	conn, err := client.Dial()
	if err != nil {
		glog.Fatalf("dial failed: %v", err)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer conn.Close()

	times := 0
	data, _ := gpacket.PackMessage(&gpacket.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte("hello server~~"),
	})

	for {
		select {
		case <-ticker.C:
			if err = conn.Push(data); err != nil {
				glog.Errorf("push message failed: %v", err)
				return
			}

			times++

			if times >= 5 {
				return
			}
		}
	}
}

func TestClient_Benchmark(t *testing.T) {
	// 并发数
	concurrency := 1000
	// 消息量
	total := 1000000
	// 总共发送的消息条数
	totalSent := int64(0)
	// 总共接收的消息条数
	totalRecv := int64(0)

	// 准备消息
	msg, err := gpacket.PackMessage(&gpacket.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte("hello server~~"),
	})
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	client := ws.NewClient()
	client.OnReceive(func(conn gnetwork.Conn, msg []byte) {
		atomic.AddInt64(&totalRecv, 1)

		wg.Done()
	})

	wg.Add(total)

	chMsg := make(chan struct{}, total)

	// 准备连接
	conns := make([]gnetwork.Conn, concurrency)
	for i := 0; i < concurrency; i++ {
		conn, err := client.Dial()
		if err != nil {
			fmt.Println("connect failed", i, err)
			i--
			continue
		}

		conns[i] = conn
		time.Sleep(time.Millisecond)
	}

	// 发送消息
	for _, conn := range conns {
		go func(conn gnetwork.Conn) {
			defer conn.Close(true)

			for {
				select {
				case _, ok := <-chMsg:
					if !ok {
						return
					}

					if err = conn.Push(msg); err != nil {
						t.Error(err)
						return
					}

					atomic.AddInt64(&totalSent, 1)
				}
			}
		}(conn)
	}

	startTime := time.Now().UnixNano()

	for i := 0; i < total; i++ {
		chMsg <- struct{}{}
	}

	wg.Wait()
	close(chMsg)

	totalTime := float64(time.Now().UnixNano()-startTime) / float64(time.Second)

	fmt.Printf("server               : %s\n", "websocket")
	fmt.Printf("concurrency          : %d\n", concurrency)
	fmt.Printf("latency              : %fs\n", totalTime)
	fmt.Printf("sent     requests    : %d\n", totalSent)
	fmt.Printf("received requests    : %d\n", totalRecv)
	fmt.Printf("throughput  (TPS)    : %d\n", int64(float64(totalRecv)/totalTime))
}

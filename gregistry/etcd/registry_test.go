package etcd_test

import (
	"context"
	"fmt"
	"gitee.com/monobytes/gcore/gcluster"
	"gitee.com/monobytes/gcore/gregistry"
	"gitee.com/monobytes/gcore/gregistry/etcd"
	"gitee.com/monobytes/gcore/gwrap/net"
	"testing"
	"time"
)

const (
	port        = 3553
	serviceName = "node"
)

var reg = etcd.NewRegistry()

func TestRegistry_Register1(t *testing.T) {
	host, err := net.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	ins := &gregistry.ServiceInstance{
		ID:       "test-1",
		Name:     serviceName,
		Kind:     gcluster.Node.String(),
		Alias:    "login-server",
		State:    gcluster.Work.String(),
		Endpoint: fmt.Sprintf("grpc://%s:%d", host, port),
	}

	rctx, rcancel := context.WithTimeout(ctx, 2*time.Second)
	err = reg.Register(rctx, ins)
	rcancel()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	ins.State = gcluster.Busy.String()
	rctx, rcancel = context.WithTimeout(ctx, 2*time.Second)
	err = reg.Register(rctx, ins)
	rcancel()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("register")
	}

	time.Sleep(20 * time.Second)

	if err = reg.Deregister(ctx, ins); err != nil {
		t.Fatal(err)
	} else {
		t.Log("deregister")
	}

	time.Sleep(40 * time.Second)
}

func TestRegistry_Register2(t *testing.T) {
	host, err := net.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}

	if err = reg.Register(context.Background(), &gregistry.ServiceInstance{
		ID:       "test-2",
		Name:     serviceName,
		Kind:     gcluster.Node.String(),
		State:    gcluster.Work.String(),
		Endpoint: fmt.Sprintf("grpc://%s:%d", host, port),
	}); err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(5 * time.Second)
		reg.Close()
	}()

	time.Sleep(30 * time.Second)
}

func TestRegistry_Services(t *testing.T) {
	services, err := reg.Services(context.Background(), serviceName)
	if err != nil {
		t.Fatal(err)
	}

	for _, service := range services {
		t.Logf("%+v", service)
	}
}

func TestRegistry_Watch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	watcher1, err := reg.Watch(ctx, serviceName)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	watcher2, err := reg.Watch(ctx, serviceName)
	cancel()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		//time.Sleep(5 * time.Second)
		//watcher1.Close()
		//time.Sleep(5 * time.Second)
		//watcher2.Close()
		//time.Sleep(5 * time.Second)
		//reg.Close()
	}()

	go func() {
		for {
			services, err := watcher1.Next()
			if err != nil {
				t.Errorf("goroutine 1: %v", err)
				return
			}

			fmt.Println("goroutine 1: new event entity")

			for _, service := range services {
				t.Logf("goroutine 1: %+v", service)
			}
		}
	}()

	go func() {
		for {
			services, err := watcher2.Next()
			if err != nil {
				t.Errorf("goroutine 2: %v", err)
				return
			}

			fmt.Println("goroutine 2: new event entity")

			for _, service := range services {
				t.Logf("goroutine 2: %+v", service)
			}
		}
	}()

	time.Sleep(60 * time.Second)
}

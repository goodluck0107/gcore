package geventbus_test

import (
	"context"
	"github.com/goodluck0107/gcore/geventbus"
	"log"
	"testing"
	"time"
)

const (
	loginTopic = "login"
	paidTopic  = "paid"
)

var eb = geventbus.NewEventbus()

func loginEventHandler(event *geventbus.Event) {
	log.Printf("%+v\n", event)
}

func paidEventHandler(event *geventbus.Event) {
	log.Printf("%+v\n", event)
}

func TestEventbus_Subscribe(t *testing.T) {
	var (
		err error
		ctx = context.Background()
	)

	err = eb.Subscribe(ctx, loginTopic, loginEventHandler)
	if err != nil {
		t.Fatal(err)
	}

	err = eb.Subscribe(ctx, paidTopic, paidEventHandler)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("subscribe success")

	err = eb.Publish(ctx, loginTopic, "login")
	if err != nil {
		t.Fatal(err)
	}

	err = eb.Publish(ctx, paidTopic, "paid")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("publish success")

	time.Sleep(30 * time.Second)
}

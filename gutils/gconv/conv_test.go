package gconv_test

import (
	"github.com/goodluck0107/gcore/gutils/gconv"
	"github.com/goodluck0107/gcore/gutils/gtime"
	"math"
	"math/cmplx"
	"testing"
)

func TestInt64(t *testing.T) {
	a := cmplx.Exp(1i*math.Pi) + 20
	t.Log(gconv.Int64(gtime.Now()))
	t.Log(gconv.Int64(func() {}))
	t.Log(gconv.Int64(&a))
}

func TestString(t *testing.T) {
	t.Log(gconv.String(1))
	t.Log(gconv.String(int8(1)))

	var a = int64(1)
	var b = 1.1
	var c = &b

	t.Log(gconv.String(&a))
	t.Log(*c)
	t.Log(gconv.String(&b))

	slice := []string{"1"}
	fun := func() {}
	t.Log(gconv.String(&slice))
	t.Log(gconv.String(fun))
}

func TestBool(t *testing.T) {
	a := float32(0)
	t.Log(gconv.Bool(a))
	t.Log(gconv.Bool(&a))
}

func TestDuration(t *testing.T) {
	t.Log(gconv.Duration("3d5m4h0.4d"))
}

func TestStrings(t *testing.T) {
	a := []int64{1, 2, 3, 4}
	t.Log(gconv.Strings(a))
}

func TestBytes(t *testing.T) {
	t.Log(gconv.Bytes("1"))
	t.Log(gconv.Int(gconv.String(gconv.Bytes("1"))))
	t.Log(gconv.Bytes(1))
	t.Log(gconv.Int(gconv.Bytes(1)))
	t.Log(gconv.Bytes(uint8(255)))
	t.Log(gconv.Int(gconv.Bytes(uint8(255))))
	t.Log(gconv.Bytes(255))
	t.Log(gconv.Int(gconv.Bytes(255)))
}

func TestAnys(t *testing.T) {
	a := []int64{1, 2, 3, 4}
	t.Log(gconv.Anys(a))
}

func TestJson(t *testing.T) {
	t.Log(gconv.Json("{}"))
	t.Log(gconv.Json(`{"id":1,"name":"fuxiao"}`))
	t.Log(gconv.Json("[]"))
	t.Log(gconv.Json(`[{"id":1,"name":"fuxiao"}]`))
	t.Log(gconv.Json(map[string]interface{}{
		"id":   1,
		"name": "fuxiao",
	}))
	t.Log(gconv.Json([]map[string]interface{}{{
		"id":   1,
		"name": "fuxiao",
	}}))
	t.Log(gconv.Json(struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{
		ID:   1,
		Name: "fuxiao",
	}))
	t.Log(gconv.Json([]struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{
		{
			ID:   1,
			Name: "fuxiao",
		},
	}))
}

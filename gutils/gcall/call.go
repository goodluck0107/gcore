package gcall

import (
	"context"
	"fmt"
	"github.com/goodluck0107/gcore/glog"
	"reflect"
	"strings"
	"time"
)

func Go(fn ...func()) {
	f := WrapperFuncs(fn...)
	go f()
}

// GoWithTimeout 执行多个协程（附带超时时间）
func GoWithTimeout(timeout time.Duration, fns ...func()) {
	NewGoroutines().Add(fns...).Run(context.Background(), timeout)
}

// GoWithDeadline 执行多个协程（附带最后期限）
func GoWithDeadline(deadline time.Time, fns ...func()) {
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	NewGoroutines().Add(fns...).Run(ctx)
}

func Call(fn ...func()) {
	f := WrapperFuncs(fn...)
	f()
}

func CallWithTimeout(duration time.Duration, fn ...func()) (timeout bool) {
	ch := make(chan struct{}, 1)

	go func(ch chan struct{}, fn []func()) {
		defer close(ch)
		WrapperFuncs(fn...)()
	}(ch, fn)

	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ch:
	case <-timer.C:
		timeout = true
		glog.Error("call:", fn, "timeout:", timeout)
	}
	return timeout
}

func Goes(fn ...func()) {
	for _, f := range fn {
		go WrapperFuncs(f)()
	}
}

func Method(v any, methodName string, in ...any) ([]any, error) {
	esValue := reflect.ValueOf(v).Elem()
	method := esValue.MethodByName("methodName")
	if !method.IsValid() {
		return nil, fmt.Errorf("method: %s not found", methodName)
	}
	args := make([]reflect.Value, 0, len(in))
	for _, i := range in {
		args = append(args, reflect.ValueOf(i))
	}
	got := method.Call(args)
	out := make([]any, 0, len(got))
	for _, i := range got {
		out = append(out, i.Interface())
	}
	return out, nil
}

func ThroughMethod(structs ...any) {
	actualContextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	actualStringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	for _, v := range structs {
		vt := reflect.TypeOf(v)
		vv := reflect.ValueOf(v)
		for i, l := 0, vt.NumMethod(); i < l; i++ {
			methodType := vt.Method(i)
			methodValue := vv.Method(i)
			numIn := methodType.Type.NumIn()
			if methodType.Type.IsVariadic() {
				numIn--
			}
			in := make([]reflect.Value, 0, numIn)
			//in = append(in, vv)
			for j := 1; j < numIn; j++ {
				if methodType.Type.In(j).Implements(actualContextType) {
					in = append(in, reflect.ValueOf(context.TODO()))
				} else {
					in = append(in, reflect.New(methodType.Type.In(j).Elem()))
				}
			}
			out := methodValue.Call(in)
			arguments := make([]string, len(in))
			for idx, val := range in {
				if val.Type().Implements(actualStringerType) {
					arguments[idx] = fmt.Sprintf("%s{%v}", val.Type().String(), val)
				} else if val.Type().Kind() == reflect.Ptr {
					arguments[idx] = fmt.Sprintf("%#v", val.Interface())
				} else {
					arguments[idx] = fmt.Sprintf("%v", val.Interface())
				}
			}
			returns := make([]any, len(out))
			for idx, val := range out {
				returns[idx] = val.Interface()
			}
			glog.Debugf("ThroughMethod: call %s.%s(%v), returns:%v", vt.String(), methodType.Name, strings.Join(arguments, ", "), returns)
		}
	}
}

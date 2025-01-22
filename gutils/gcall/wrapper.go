package gcall

import (
	"bytes"
	"fmt"
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/gutils/gconv"
	"runtime"
)

func AddPanicCallback(cb func(x any, stack string)) {
	panicCallbacks = append(panicCallbacks, cb)
}

var panicCallbacks []func(x any, stack string)

func WrapperFuncs(fn ...func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case runtime.Error:
					glog.Panic(err)
				default:
					glog.Panicf("panic error: %v", err)
				}
				s := []byte("/src/runtime/panic.go")
				e := []byte("\ngoroutine ")
				line := []byte("\n")
				stack := make([]byte, 4096) //4KB
				length := runtime.Stack(stack, true)
				start := bytes.Index(stack, s)
				stack = stack[start:length]
				start = bytes.Index(stack, line) + 1
				stack = stack[start:]
				end := bytes.LastIndex(stack, line)
				if end != -1 {
					stack = stack[:end]
				}
				end = bytes.Index(stack, e)
				if end != -1 {
					stack = stack[:end]
				}
				stackDetail := gconv.BytesToString(bytes.TrimRight(stack, "\n"))
				for _, cb := range panicCallbacks {
					cb(err, stackDetail)
				}
			}
		}()
		for _, f := range fn {
			f()
		}
	}
}

func Recover(args ...interface{}) any {
	x := recover()
	var normalHandlers []func()
	for i := 0; i < len(args); i++ {
		if ph, ok := args[i].(func()); ok {
			normalHandlers = append(normalHandlers, ph)
			args = append(args[:i], args[i+1:]...)
			i--
		}
	}
	if x != nil {
		s := []byte("/src/runtime/panic.go")
		e := []byte("\ngoroutine ")
		line := []byte("\n")
		stack := make([]byte, 4096) //4KB
		length := runtime.Stack(stack, true)
		start := bytes.Index(stack, s)
		stack = stack[start:length]
		start = bytes.Index(stack, line) + 1
		stack = stack[start:]
		end := bytes.LastIndex(stack, line)
		if end != -1 {
			stack = stack[:end]
		}
		end = bytes.Index(stack, e)
		if end != -1 {
			stack = stack[:end]
		}
		stack = bytes.TrimRight(stack, "\n")
		var panicHandlers []func(any)
		for i := 0; i < len(args); i++ {
			if ph, ok := args[i].(func(any)); ok {
				panicHandlers = append(panicHandlers, ph)
				args = append(args[:i], args[i+1:]...)
				i--
			}
		}
		additionalInfo := make(map[string]interface{})
		if len(args)%2 == 0 {
			for i, v := range args {
				if i%2 == 0 && i+1 < len(args) {
					if column, ok := v.(string); ok {
						additionalInfo[column] = args[i+1]
					} else {
						additionalInfo[fmt.Sprint(v)] = args[i+1]
					}
				}
			}
		} else {
			glog.Error("xruntime.Recover(...) args count error")
		}
		var buf bytes.Buffer
		buf.WriteString("\n")
		for k, v := range additionalInfo {
			buf.WriteString(fmt.Sprintf("%s:%v\n", k, v))
		}
		buf.WriteString(fmt.Sprintf("\n%v\n%s\n", x, string(stack)))
		panicInfo := buf.String()

		glog.Error(panicInfo)
		for _, h := range panicCallbacks {
			h(x, panicInfo)
		}
		for _, h := range panicHandlers {
			h(x)
		}
	} else {
		for _, h := range normalHandlers {
			h()
		}
	}
	return x
}

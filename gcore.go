package gcore

import (
	"gitee.com/monobytes/gcore/gengine"
)

var engine gengine.Engine

func Engine() gengine.Engine {
	return engine
}

func SetEngine(e gengine.Engine) {
	engine = e
}

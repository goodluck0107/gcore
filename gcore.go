package gcore

import (
	"github.com/goodluck0107/gcore/gengine"
)

var engine gengine.Engine

func Engine() gengine.Engine {
	return engine
}

func SetEngine(e gengine.Engine) {
	engine = e
}

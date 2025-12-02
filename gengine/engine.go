package gengine

import "github.com/goodluck0107/gcore/gmodules"

type Engine interface {
	Injection(mods ...gmodules.Module)
	Up()
}

package gengine

import "gitee.com/monobytes/gcore/gmodules"

type Engine interface {
	Injection(mods ...gmodules.Module)
	Up()
}

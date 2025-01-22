package info

import (
	"fmt"
	"gitee.com/monobytes/gcore/gmode"
	"strings"
	"syscall"
	"unicode/utf8"
)

const logo = `
		 _____  _____  _____  _____  _____ 
		/   __\/     \/  _  \/  _  \/   __\
		|  |_ ||  |--||  |  ||  _  <|   __|
		\_____/\_____/\_____/\__|\_/\_____/
`

const (
	boxWidth          = 56
	verticalBorder    = "|"
	horizontalBorder  = "─"
	leftTopBorder     = "┌"
	rightTopBorder    = "┐"
	leftBottomBorder  = "└"
	rightBottomBorder = "┘"
	website           = "https://gitee.com/monobytes/gcore"
	version           = "v0.0.1"
	global            = "Global"
)

func PrintFrameworkInfo() {
	fmt.Println(strings.TrimSuffix(strings.TrimPrefix(logo, "\n"), "\n"))
	PrintBoxInfo("",
		fmt.Sprintf("[Website] %s", website),
		fmt.Sprintf("[Version] %s", version),
	)
}

func PrintGlobalInfo() {
	PrintBoxInfo("Global",
		fmt.Sprintf("PID: %d", syscall.Getpid()),
		fmt.Sprintf("Mode: %s", gmode.GetMode()),
	)
}

func PrintBoxInfo(name string, infos ...string) {
	fmt.Println(buildTopBorder(name))
	for _, info := range infos {
		fmt.Println(buildRowInfo(info))
	}
	fmt.Println(buildBottomBorder())
}

func buildRowInfo(info string) string {
	str := fmt.Sprintf("%s %s", verticalBorder, info)
	str += strings.Repeat(" ", boxWidth-utf8.RuneCountInString(str)-1)
	str += verticalBorder
	return str
}

func buildTopBorder(name ...string) string {
	full := boxWidth - strLen(leftTopBorder) - strLen(rightTopBorder) - strLen(name...)
	half := full / 2
	str := leftTopBorder
	str += strings.Repeat(horizontalBorder, half)
	if len(name) > 0 {
		str += name[0]
	}
	str += strings.Repeat(horizontalBorder, full-half)
	str += rightTopBorder
	return str
}

func buildBottomBorder() string {
	full := boxWidth - strLen(leftBottomBorder) - strLen(rightBottomBorder)
	str := leftBottomBorder
	str += strings.Repeat(horizontalBorder, full)
	str += rightBottomBorder
	return str
}

func strLen(str ...string) int {
	if len(str) > 0 {
		return utf8.RuneCountInString(str[0])
	} else {
		return 0
	}
}

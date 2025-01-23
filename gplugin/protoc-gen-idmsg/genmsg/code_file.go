package genmsg

import (
	"fmt"
	"strings"
)

type CodeFile struct {
	lines []string
}

func NewCodeFile() *CodeFile {
	ret := &CodeFile{}
	return ret
}

func (f *CodeFile) P(line string) {
	f.lines = append(f.lines, line)
}

func (f *CodeFile) F(format string, args ...interface{}) {
	f.lines = append(f.lines, fmt.Sprintf(format, args...))
}

func (f *CodeFile) String() string {
	return strings.Join(f.lines, "\n")
}

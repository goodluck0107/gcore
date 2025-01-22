package hook

import (
	"gitee.com/monobytes/gcore/glog"
	"gitee.com/monobytes/gcore/glog/logrus/internal/define"
	"gitee.com/monobytes/gcore/gwrap/stack"
	"github.com/sirupsen/logrus"
)

type StackHook struct {
	stackLevel glog.Level
	callerSkip int
}

func NewStackHook(stackLevel glog.Level, callerSkip int) *StackHook {
	return &StackHook{stackLevel: stackLevel, callerSkip: callerSkip}
}

func (h *StackHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *StackHook) Fire(entry *logrus.Entry) error {
	var level glog.Level
	switch entry.Level {
	case logrus.DebugLevel:
		level = glog.DebugLevel
	case logrus.InfoLevel:
		level = glog.InfoLevel
	case logrus.WarnLevel:
		level = glog.WarnLevel
	case logrus.ErrorLevel:
		level = glog.ErrorLevel
	case logrus.FatalLevel:
		level = glog.FatalLevel
	case logrus.PanicLevel:
		level = glog.PanicLevel
	}

	depth := stack.First
	if _, ok := entry.Data[define.StackOutFlagField]; ok {
		if h.stackLevel != glog.NoneLevel && level >= h.stackLevel {
			depth = stack.Full
		} else {
			delete(entry.Data, define.StackOutFlagField)
		}
	}

	st := stack.Callers(8+h.callerSkip, depth)
	defer st.Free()
	entry.Data[define.StackFramesFlagField] = st.Frames()

	return nil
}

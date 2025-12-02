package logger

import (
	"github.com/goodluck0107/gcore/glog"
	rpcxlog "github.com/smallnest/rpcx/log"
	"sync"
)

var once sync.Once

func InitLogger() {
	once.Do(func() {
		rpcxlog.SetLogger(&logger{
			level:  glog.ErrorLevel,
			logger: glog.GetLogger(),
		})
	})
}

type logger struct {
	level  glog.Level
	logger glog.Logger
}

func (l *logger) Debug(v ...interface{}) {
	if l.level <= glog.DebugLevel {
		l.logger.Print(glog.DebugLevel, v...)
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.level <= glog.DebugLevel {
		l.logger.Printf(glog.DebugLevel, format, v...)
	}
}

func (l *logger) Info(v ...interface{}) {
	if l.level <= glog.InfoLevel {
		l.logger.Print(glog.InfoLevel, v...)
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.level <= glog.InfoLevel {
		l.logger.Printf(glog.InfoLevel, format, v...)
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.level <= glog.WarnLevel {
		l.logger.Print(glog.WarnLevel, v...)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level <= glog.WarnLevel {
		l.logger.Printf(glog.WarnLevel, format, v...)
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.level <= glog.ErrorLevel {
		l.logger.Print(glog.ErrorLevel, v...)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if l.level <= glog.ErrorLevel {
		l.logger.Printf(glog.ErrorLevel, format, v...)
	}
}

func (l *logger) Fatal(v ...interface{}) {
	if l.level <= glog.FatalLevel {
		l.logger.Print(glog.FatalLevel, v...)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	if l.level <= glog.FatalLevel {
		l.logger.Printf(glog.FatalLevel, format, v...)
	}
}

func (l *logger) Panic(v ...interface{}) {
	if l.level <= glog.PanicLevel {
		l.logger.Print(glog.PanicLevel, v...)
	}
}

func (l *logger) Panicf(format string, v ...interface{}) {
	if l.level <= glog.PanicLevel {
		l.logger.Printf(glog.PanicLevel, format, v...)
	}
}

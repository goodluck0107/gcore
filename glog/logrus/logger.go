package logrus

import (
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/glog/logrus/internal/define"
	"github.com/goodluck0107/gcore/glog/logrus/internal/formatter"
	"github.com/goodluck0107/gcore/glog/logrus/internal/hook"
	"github.com/goodluck0107/gcore/gmode"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var _ glog.Logger = NewLogger()

type closer interface {
	Close() error
}

type Logger struct {
	opts    *options
	logger  *logrus.Logger
	writers []io.Writer
}

func NewLogger(opts ...Option) *Logger {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	l := &Logger{opts: o, logger: logrus.New(), writers: make([]io.Writer, 0, 6)}

	switch o.level {
	case glog.DebugLevel:
		l.logger.SetLevel(logrus.DebugLevel)
	case glog.InfoLevel:
		l.logger.SetLevel(logrus.InfoLevel)
	case glog.WarnLevel:
		l.logger.SetLevel(logrus.WarnLevel)
	case glog.ErrorLevel:
		l.logger.SetLevel(logrus.ErrorLevel)
	case glog.FatalLevel:
		l.logger.SetLevel(logrus.FatalLevel)
	case glog.PanicLevel:
		l.logger.SetLevel(logrus.PanicLevel)
	}

	switch o.format {
	case glog.JsonFormat:
		l.logger.SetFormatter(&formatter.JsonFormatter{
			TimeFormat:     o.timeFormat,
			CallerFullPath: o.callerFullPath,
		})
	default:
		l.logger.SetFormatter(&formatter.TextFormatter{
			TimeFormat:     o.timeFormat,
			CallerFullPath: o.callerFullPath,
		})
	}

	l.logger.AddHook(hook.NewStackHook(o.stackLevel, o.callerSkip))

	if o.file != "" {
		if o.classifiedStorage {
			writers := hook.WriterMap{
				logrus.DebugLevel: l.buildWriter(glog.DebugLevel),
				logrus.InfoLevel:  l.buildWriter(glog.InfoLevel),
				logrus.WarnLevel:  l.buildWriter(glog.WarnLevel),
				logrus.ErrorLevel: l.buildWriter(glog.ErrorLevel),
				logrus.FatalLevel: l.buildWriter(glog.FatalLevel),
				logrus.PanicLevel: l.buildWriter(glog.PanicLevel),
			}

			for key := range writers {
				l.writers = append(l.writers, writers[key])
			}

			l.logger.AddHook(hook.NewWriterHook(writers))
		} else {
			writer := l.buildWriter(glog.NoneLevel)
			l.writers = append(l.writers, writer)
			l.logger.AddHook(hook.NewWriterHook(writer))
		}
	}

	if gmode.IsDebugMode() && o.stdout {
		l.logger.SetOutput(os.Stdout)
	}

	return l
}

func (l *Logger) buildWriter(level glog.Level) io.Writer {
	writer, err := glog.NewWriter(glog.WriterOptions{
		Path:    l.opts.file,
		Level:   level,
		MaxAge:  l.opts.fileMaxAge,
		MaxSize: l.opts.fileMaxSize * 1024 * 1024,
		CutRule: l.opts.fileCutRule,
	})
	if err != nil {
		panic(err)
	}

	return writer
}

// Print 打印日志，不含堆栈信息
func (l *Logger) Print(level glog.Level, a ...interface{}) {
	switch level {
	case glog.DebugLevel:
		l.logger.Debug(a...)
	case glog.InfoLevel:
		l.logger.Info(a...)
	case glog.WarnLevel:
		l.logger.Warn(a...)
	case glog.ErrorLevel:
		l.logger.Error(a...)
	case glog.FatalLevel:
		l.logger.Fatal(a...)
	case glog.PanicLevel:
		l.logger.Panic(a...)
	}
}

// Printf 打印模板日志，不含堆栈信息
func (l *Logger) Printf(level glog.Level, format string, a ...interface{}) {
	switch level {
	case glog.DebugLevel:
		l.logger.Debugf(format, a...)
	case glog.InfoLevel:
		l.logger.Infof(format, a...)
	case glog.WarnLevel:
		l.logger.Warnf(format, a...)
	case glog.ErrorLevel:
		l.logger.Errorf(format, a...)
	case glog.FatalLevel:
		l.logger.Fatalf(format, a...)
	case glog.PanicLevel:
		l.logger.Panicf(format, a...)
	}
}

// Debug 打印调试日志
func (l *Logger) Debug(a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Debug(a...)
}

// Debugf 打印调试模板日志
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Debugf(format, a...)
}

// Info 打印信息日志
func (l *Logger) Info(a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Info(a...)
}

// Infof 打印信息模板日志
func (l *Logger) Infof(format string, a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Infof(format, a...)
}

// Warn 打印警告日志
func (l *Logger) Warn(a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Warn(a...)
}

// Warnf 打印警告模板日志
func (l *Logger) Warnf(format string, a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Warnf(format, a...)
}

// Error 打印错误日志
func (l *Logger) Error(a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Error(a...)
}

// Errorf 打印错误模板日志
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Errorf(format, a...)
}

// Fatal 打印致命错误日志
func (l *Logger) Fatal(a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Fatal(a...)
}

// Fatalf 打印致命错误模板日志
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Fatalf(format, a...)
}

// Panic 打印Panic日志
func (l *Logger) Panic(a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Fatal(a...)
}

// Panicf 打印Panic模板日志
func (l *Logger) Panicf(format string, a ...interface{}) {
	l.logger.WithField(define.StackOutFlagField, true).Fatalf(format, a...)
}

// Close 关闭日志
func (l *Logger) Close() (err error) {
	for _, writer := range l.writers {
		w, ok := writer.(interface{ Close() error })
		if !ok {
			continue
		}

		if e := w.Close(); e != nil {
			err = e
		}
	}

	return
}

package zap

import (
	"fmt"
	"github.com/goodluck0107/gcore/glog"
	"github.com/goodluck0107/gcore/glog/zap/internal/encoder"
	"github.com/goodluck0107/gcore/gmode"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var _ glog.Logger = NewLogger()

var levelMap map[zapcore.Level]glog.Level

func init() {
	levelMap = map[zapcore.Level]glog.Level{
		zap.DebugLevel:  glog.DebugLevel,
		zap.InfoLevel:   glog.InfoLevel,
		zap.WarnLevel:   glog.WarnLevel,
		zap.ErrorLevel:  glog.ErrorLevel,
		zap.FatalLevel:  glog.FatalLevel,
		zap.DPanicLevel: glog.PanicLevel,
	}
}

type Logger struct {
	logger *zap.SugaredLogger
	opts   *options
}

func NewLogger(opts ...Option) *Logger {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	var (
		fileEncoder     zapcore.Encoder
		terminalEncoder zapcore.Encoder
	)
	switch o.format {
	case glog.JsonFormat:
		fileEncoder = encoder.NewJsonEncoder(o.timeFormat, o.callerFullPath)
		terminalEncoder = fileEncoder
	default:
		fileEncoder = encoder.NewTextEncoder(o.timeFormat, o.callerFullPath, false)
		terminalEncoder = encoder.NewTextEncoder(o.timeFormat, o.callerFullPath, true)
	}

	options := make([]zap.Option, 0, 3)
	options = append(options, zap.AddCaller())
	switch o.stackLevel {
	case glog.DebugLevel:
		options = append(options, zap.AddStacktrace(zapcore.DebugLevel), zap.AddCallerSkip(1+o.callerSkip))
	case glog.InfoLevel:
		options = append(options, zap.AddStacktrace(zapcore.InfoLevel), zap.AddCallerSkip(1+o.callerSkip))
	case glog.WarnLevel:
		options = append(options, zap.AddStacktrace(zapcore.WarnLevel), zap.AddCallerSkip(1+o.callerSkip))
	case glog.ErrorLevel:
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1+o.callerSkip))
	case glog.FatalLevel:
		options = append(options, zap.AddStacktrace(zapcore.FatalLevel), zap.AddCallerSkip(1+o.callerSkip))
	case glog.PanicLevel:
		options = append(options, zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1+o.callerSkip))
	}

	l := &Logger{opts: o}

	var cores []zapcore.Core
	if o.file != "" {
		if o.classifiedStorage {
			cores = append(cores,
				zapcore.NewCore(fileEncoder, l.buildWriteSyncer(glog.DebugLevel), l.buildLevelEnabler(glog.DebugLevel)),
				zapcore.NewCore(fileEncoder, l.buildWriteSyncer(glog.InfoLevel), l.buildLevelEnabler(glog.InfoLevel)),
				zapcore.NewCore(fileEncoder, l.buildWriteSyncer(glog.WarnLevel), l.buildLevelEnabler(glog.WarnLevel)),
				zapcore.NewCore(fileEncoder, l.buildWriteSyncer(glog.ErrorLevel), l.buildLevelEnabler(glog.ErrorLevel)),
				zapcore.NewCore(fileEncoder, l.buildWriteSyncer(glog.FatalLevel), l.buildLevelEnabler(glog.FatalLevel)),
				zapcore.NewCore(fileEncoder, l.buildWriteSyncer(glog.PanicLevel), l.buildLevelEnabler(glog.PanicLevel)),
			)
		} else {
			cores = append(cores, zapcore.NewCore(fileEncoder, l.buildWriteSyncer(glog.NoneLevel), l.buildLevelEnabler(glog.NoneLevel)))
		}
	}

	if gmode.IsDebugMode() && o.stdout {
		cores = append(cores, zapcore.NewCore(terminalEncoder, zapcore.AddSync(os.Stdout), l.buildLevelEnabler(glog.NoneLevel)))
	}

	if len(cores) >= 0 {
		l.logger = zap.New(zapcore.NewTee(cores...), options...).Sugar()
	}

	return l
}

func (l *Logger) buildWriteSyncer(level glog.Level) zapcore.WriteSyncer {
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

	return zapcore.AddSync(writer)
}

func (l *Logger) buildLevelEnabler(level glog.Level) zapcore.LevelEnabler {
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if v := levelMap[lvl]; l.opts.level != glog.NoneLevel {
			return v >= l.opts.level && (level == glog.NoneLevel || (level >= l.opts.level && v >= level))
		} else {
			return level == glog.NoneLevel || v >= level
		}
	})
}

// 打印日志
func (l *Logger) print(level glog.Level, stack bool, a ...interface{}) {
	if l.logger == nil {
		return
	}

	var msg string
	if len(a) == 1 {
		if str, ok := a[0].(string); ok {
			msg = str
		} else {
			msg = fmt.Sprint(a...)
		}
	} else {
		msg = fmt.Sprint(a...)
	}

	switch level {
	case glog.DebugLevel:
		l.logger.Debugw(msg, encoder.StackFlag, stack)
	case glog.InfoLevel:
		l.logger.Infow(msg, encoder.StackFlag, stack)
	case glog.WarnLevel:
		l.logger.Warnw(msg, encoder.StackFlag, stack)
	case glog.ErrorLevel:
		l.logger.Errorw(msg, encoder.StackFlag, stack)
	case glog.FatalLevel:
		l.logger.Fatalw(msg, encoder.StackFlag, stack)
	case glog.PanicLevel:
		l.logger.DPanicw(msg, encoder.StackFlag, stack)
	}
}

// Print 打印日志，不打印堆栈信息
func (l *Logger) Print(level glog.Level, a ...interface{}) {
	l.print(level, false, a...)
}

// Printf 打印模板日志，不打印堆栈信息
func (l *Logger) Printf(level glog.Level, format string, a ...interface{}) {
	l.print(level, false, fmt.Sprintf(format, a...))
}

// Debug 打印调试日志
func (l *Logger) Debug(a ...interface{}) {
	l.print(glog.DebugLevel, true, a...)
}

// Debugf 打印调试模板日志
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.print(glog.DebugLevel, true, fmt.Sprintf(format, a...))
}

// Info 打印信息日志
func (l *Logger) Info(a ...interface{}) {
	l.print(glog.InfoLevel, true, a...)
}

// Infof 打印信息模板日志
func (l *Logger) Infof(format string, a ...interface{}) {
	l.print(glog.InfoLevel, true, fmt.Sprintf(format, a...))
}

// Warn 打印警告日志
func (l *Logger) Warn(a ...interface{}) {
	l.print(glog.WarnLevel, true, a...)
}

// Warnf 打印警告模板日志
func (l *Logger) Warnf(format string, a ...interface{}) {
	l.print(glog.WarnLevel, true, fmt.Sprintf(format, a...))
}

// Error 打印错误日志
func (l *Logger) Error(a ...interface{}) {
	l.print(glog.ErrorLevel, true, a...)
}

// Errorf 打印错误模板日志
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.print(glog.ErrorLevel, true, fmt.Sprintf(format, a...))
}

// Fatal 打印致命错误日志
func (l *Logger) Fatal(a ...interface{}) {
	l.print(glog.FatalLevel, true, a...)
}

// Fatalf 打印致命错误模板日志
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.print(glog.FatalLevel, true, fmt.Sprintf(format, a...))
}

// Panic 打印Panic日志
func (l *Logger) Panic(a ...interface{}) {
	l.print(glog.PanicLevel, true, a...)
}

// Panicf 打印Panic模板日志
func (l *Logger) Panicf(format string, a ...interface{}) {
	l.print(glog.PanicLevel, true, fmt.Sprintf(format, a...))
}

// Sync 同步缓存中的日志
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// Close 关闭日志
func (l *Logger) Close() error {
	return l.logger.Sync()
}

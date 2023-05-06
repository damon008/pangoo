package util

import (
	"context"
	"fmt"
	//"github.com/bytedance/gopkg/util/logger"
	//"gopkg.in/natefinch/lumberjack.v2"

	"github.com/natefinch/lumberjack"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/sirupsen/logrus"
)

type KLogger struct {
}

type Logger struct {
	l *logrus.Logger
}

var logger = logrus.New()

func NewKLogger(logFilePath string) *Logger {
	if err := os.MkdirAll(logFilePath, 0o777); err != nil {
		log.Println(err.Error())
	}
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
		}
	}
	//logger := logrus.New()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    20,   // A file can be up to 20M.
		MaxBackups: 5,    // Save up to 5 files at the same time.
		MaxAge:     10,   // A file can exist for a maximum of 10 days.
		Compress:   true, // Compress with gzip.
	}
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetOutput(lumberjackLogger)
	logger.SetLevel(logrus.DebugLevel)
	return nil
	//return &Logger{l: logger}
}

func NewHLogger(logFilePath string) *Logger {
	if err := os.MkdirAll(logFilePath, 0o777); err != nil {
		log.Println(err.Error())
	}
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
		}
	}
	//logger := logrus.New()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    20,   // A file can be up to 20M.
		MaxBackups: 5,    // Save up to 5 files at the same time.
		MaxAge:     10,   // A file can exist for a maximum of 10 days.
		Compress:   true, // Compress with gzip.
	}

	//logger.SetReportCaller(true)
	//logger.AddHook(NewCallerHook())

	logger.SetOutput(lumberjackLogger)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	/*logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		//CallerPrettyfier:  caller(),
		FieldMap: logrus.FieldMap {
			logrus.FieldKeyFile: "caller",
		},
	}*/
	return &Logger{l: logger}
}

type CallerHook struct {
}

func NewCallerHook() *CallerHook {
	return &CallerHook{}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		//这里可以调整文件名具体多长。
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// Levels get levels
func (h *CallerHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.InfoLevel}
}

// Fire logrus hook fire
func (h *CallerHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}
	// 这里可以调整下。
	entry.Data["file"] = fileInfo(2)
	return nil
}

func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

/*func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}*/

func NewLogger(logFile string) *lumberjack.Logger {
	if err := os.MkdirAll(logFile, 0o777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFile, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//logger := logrus.New()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,     // A file can be up to 20M.
		MaxBackups: 5,       // Save up to 5 files at the same time.
		MaxAge:     365 * 2, // A file can exist for a maximum of 10 days.
		Compress:   true,    // Compress with gzip.
	}
	return lumberjackLogger
}

func (l *Logger) Trace(v ...interface{}) {
	l.l.Trace(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.l.Debug(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.l.Info(v...)
}

func (l *Logger) Notice(v ...interface{}) {
	l.l.Warn(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.l.Warn(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.l.Error(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.l.Fatal(v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.l.Tracef(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.l.Debugf(format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.l.Infof(format, v...)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.l.Warnf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.l.Warnf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.l.Errorf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.l.Fatalf(format, v...)
}

func (l *Logger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	l.l.WithContext(ctx).Tracef(format, v...)
}

func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	l.l.WithContext(ctx).Debugf(format, v...)
}

func (l *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	l.l.WithContext(ctx).Infof(format, v...)
}

func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	l.l.WithContext(ctx).Warnf(format, v...)
}

func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	l.l.WithContext(ctx).Warnf(format, v...)
}

func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	l.l.WithContext(ctx).Errorf(format, v...)
}

func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	l.l.WithContext(ctx).Fatalf(format, v...)
}

func (l *Logger) SetLevel(level hlog.Level) {
	var lv logrus.Level
	switch level {
	case hlog.LevelTrace:
		lv = logrus.TraceLevel
	case hlog.LevelDebug:
		lv = logrus.DebugLevel
	case hlog.LevelInfo:
		lv = logrus.InfoLevel
	case hlog.LevelWarn, hlog.LevelNotice:
		lv = logrus.WarnLevel
	case hlog.LevelError:
		lv = logrus.ErrorLevel
	case hlog.LevelFatal:
		lv = logrus.FatalLevel
	default:
		lv = logrus.WarnLevel
	}
	l.l.SetLevel(lv)
}

func (l *Logger) SetOutput(writer io.Writer) {
	l.l.SetOutput(writer)
}

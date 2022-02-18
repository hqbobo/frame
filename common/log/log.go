package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/hqbobo/frame/common/sirupsen/logrus"
)

// Level describes the log severity level.
type Level uint8

const (
	// PanicLevel level
	PanicLevel Level = iota
	// FatalLevel level
	FatalLevel
	// ErrorLevel level
	ErrorLevel
	// WarnLevel level
	WarnLevel
	// InfoLevel level
	InfoLevel
	// DebugLevel level
	DebugLevel
	// TraceLevel level
	TraceLevel
)

// Logger is an interface that describes logging.
type Logger interface {
	SetLevel(level Level)
	SetOut(out io.Writer)

	Traceln(...interface{})
	Trace(...interface{})
	Tracef(string, ...interface{})

	Debugln(...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})

	Infoln(...interface{})
	Info(...interface{})
	Infof(string, ...interface{})

	Warnln(...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})

	Errorln(...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})

	Fatalln(...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})

	Panicln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

type logger struct {
	entry *logrus.Entry
}

// With attaches a key-value pair to a logger.
func (l logger) With(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

// WithError attaches an error to a logger.
func (l logger) WithError(err error) Logger {
	return logger{l.entry.WithError(err)}
}

// SetLevel sets the level of a logger.
func (l logger) SetLevel(level Level) {
	l.entry.Logger.Level = logrus.Level(level)
}

// SetOut sets the output destination for a logger.
func (l logger) SetOut(out io.Writer) {
	l.entry.Logger.Out = out
}

// Trace logs a message at level Debug on the standard logger.
func (l logger) Traceln(args ...interface{}) {
	l.sourced().Traceln(args...)
}

// Trace logs a message at level Debug on the standard logger.
func (l logger) Trace(args ...interface{}) {
	l.sourced().Trace(args...)
}

// Tracef logs a message at level Debug on the standard logger.
func (l logger) Tracef(format string, args ...interface{}) {
	l.sourced().Tracef(format, args...)
}

// Trace logs a message at level Debug on the standard logger.
func (l logger) Debugln(args ...interface{}) {
	l.sourced().Debugln(args...)
}

// Debug logs a message at level Debug on the standard logger.
func (l logger) Debug(args ...interface{}) {
	l.sourced().Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l logger) Debugf(format string, args ...interface{}) {
	l.sourced().Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Infoln(args ...interface{}) {
	l.sourced().Infoln(args...)
}

// Info logs a message at level Info on the standard logger.
func (l logger) Info(args ...interface{}) {
	l.sourced().Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func (l logger) Infof(format string, args ...interface{}) {
	l.sourced().Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warnln(args ...interface{}) {
	l.sourced().Warnln(args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l logger) Warn(args ...interface{}) {
	l.sourced().Warn(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func (l logger) Warnf(format string, args ...interface{}) {
	l.sourced().Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Errorln(args ...interface{}) {
	l.sourced().Errorln(args...)
}

// Error logs a message at level Error on the standard logger.
func (l logger) Error(args ...interface{}) {
	l.sourced().Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (l logger) Errorf(format string, args ...interface{}) {
	l.sourced().Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatalln(args ...interface{}) {
	l.sourced().Fatalln(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l logger) Fatal(args ...interface{}) {
	l.sourced().Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (l logger) Fatalf(format string, args ...interface{}) {
	l.sourced().Fatalf(format, args...)
}

// Panicln logs a message at level Fatal on the standard logger.
func (l logger) Panicln(args ...interface{}) {
	l.sourced().Panicln(args...)
}

// Panic logs a message at level Fatal on the standard logger.
func (l logger) Panic(args ...interface{}) {
	l.sourced().Panic(args...)
}

// Panic logs a message at level Fatal on the standard logger.
func (l logger) Panicf(format string, args ...interface{}) {
	l.sourced().Panicf(format, args...)
}

const (
	filetextlen = 30
)

// sourced adds a source field to the logger that contains
// the file name and line where the logging happened.
func (l logger) sourced() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	fn := "(unknown)"
	if !ok {
		file = "<???>"
		line = 1
	} else {
		function := runtime.FuncForPC(pc).Name()
		index := strings.LastIndex(function, ".")
		slash := strings.LastIndex(file, "go/src/")
		if len(file) > filetextlen {
			file = file[len(file)-filetextlen:]
		} else {
			file = file[slash+7:]
		}
		fn = function[index+1:]
	}
	return l.entry.WithField("PATH", fmt.Sprintf("%s() %0.30s:%d", fn, file, line))
}

var origLogger = logrus.New()
var baseLogger = logger{entry: logrus.NewEntry(origLogger)}

// init
func init() {
	baseLogger.entry.Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "06-01-02 15:04:05.0000",
		ForceColors:     true,
	})
	SetLevel(TraceLevel)
}

// New returns a new logger.
func New() Logger {
	return logger{entry: logrus.NewEntry(origLogger)}
}

// Base returns the base logger.
func Base() Logger {
	return baseLogger
}

// SetLevel sets the Level of the base logger
func SetLevel(level Level) {
	baseLogger.entry.Logger.Level = logrus.Level(level)
}

// SetOut sets the output destination base logger
func SetOut(out io.Writer) {
	baseLogger.entry.Logger.Out = out
}

// Withf attaches a key,value pair to a logger.
func Withf(key string, value interface{}) Logger {
	return baseLogger.With(key, value)
}

// WithError returns a Logger that will print an error along with the next message.
func WithError(err error) Logger {
	return logger{entry: baseLogger.sourced().WithError(err)}
}

// Traceln logs a message at level Debug on the standard logger.
func Traceln(args ...interface{}) {
	baseLogger.sourced().Traceln(args...)
}

// Traceln logs a message at level Debug on the standard logger.
func Trace(args ...interface{}) {
	baseLogger.sourced().Trace(args...)
}

// Tracef logs a message at level Debug on the standard logger.
func Tracef(format string, args ...interface{}) {
	baseLogger.sourced().Tracef(format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	baseLogger.sourced().Debugln(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	baseLogger.sourced().Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	baseLogger.sourced().Debugf(format, args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	baseLogger.sourced().Infoln(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	baseLogger.sourced().Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	baseLogger.sourced().Infof(format, args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	baseLogger.sourced().Warnln(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	baseLogger.sourced().Warn(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	baseLogger.sourced().Warnf(format, args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	baseLogger.sourced().Errorln(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	baseLogger.sourced().Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	baseLogger.sourced().Errorf(format, args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func Fatalln(args ...interface{}) {
	baseLogger.sourced().Fatalln(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	baseLogger.sourced().Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func Fatalf(format string, args ...interface{}) {
	baseLogger.sourced().Fatalf(format, args...)
}

// Panicln logs a message at level Fatal on the standard logger.
func Panicln(args ...interface{}) {
	baseLogger.sourced().Panicln(args...)
}

// Panic logs a message at level Fatal on the standard logger.
func Panic(args ...interface{}) {
	baseLogger.sourced().Panic(args...)
}

// Panicf logs a message at level Fatal on the standard logger.
func Panicf(format string, args ...interface{}) {
	baseLogger.sourced().Panicf(format, args...)
}

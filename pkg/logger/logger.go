package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	constFields "github.com/shahbaz275817/prismo/constants/fields"
)

var logger *logrus.Logger
var accessLogFile lumberjack.Logger

const logPath = "logs/"
const RFC3339Milli = "2006-01-02T15:04:05.000Z"

// Fields represents key-value pairs and can be used to
// provide additional context in logs
type Fields map[string]interface{}

// SetupLogger initialize logger with given config.
func SetupLogger(config Config) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("Error while parsing Application Log Level %s", err.Error())
	}

	logger = getNewLoggerInstance(level, logPath, "application.log")

	// setup access log file
	accessLogFile = lumberjack.Logger{
		Filename:   logPath + "access.log",
		MaxSize:    200, // megabytes
		MaxBackups: 20,
		Compress:   true,  // disabled by default
		LocalTime:  false, // optional
		MaxAge:     7,
	}

}

func getNewLoggerInstance(level logrus.Level, logPath string, filename string) *logrus.Logger {
	instance := &logrus.Logger{
		Out:   os.Stdout,
		Hooks: make(logrus.LevelHooks),
		Level: level,
		Formatter: &nested.Formatter{
			TimestampFormat: RFC3339Milli,
			CustomCallerFormatter: func(frame *runtime.Frame) string {
				return ""
			},
		},
		ReportCaller: true,
	}

	lumberjackHook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   logPath + filename,
			MaxSize:    200,
			MaxBackups: 20,
			Compress:   true,
			LocalTime:  false,
			MaxAge:     7,
		},
		level,
		&nested.Formatter{
			NoColors:        true,
			TimestampFormat: RFC3339Milli,
			CustomCallerFormatter: func(frame *runtime.Frame) string {
				return ""
			},
		},
		&lumberjackrus.LogFileOpts{},
	)
	instance.Hooks.Add(&LineNumberHook{})
	if err == nil {
		instance.Hooks.Add(lumberjackHook)
	}
	return instance
}

func GetAccessLogFile() io.WriteCloser {
	return &accessLogFile
}

// WithField captures additional information for the logger
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// WithRequest capture information from http request
func withRequest(r *http.Request) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"method": r.Method,
		"host":   r.Host,
		"path":   r.URL.Path,
	})
}

// Debug logs message at debug log level
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf logs formated message at debug log level
func Debugf(format string, args ...interface{}) {
	fields := getSourceInfoFields(3)
	fields["root"] = true
	WithFields(fields).Debugf(format, args...)
}

// Info logs message at info log level
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof logs formated message at info log level
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn logs message at warn log level
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf logs formated message at warn log level
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Error logs message at error log level
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf logs formated message at error log level
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatalf logs formated message at fatal log level
// This method also exit the program
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// AddHook adds a new hook to logger
func AddHook(hook logrus.Hook) {
	logger.Hooks.Add(hook)
}

// For making a default logger, this makes testing easier
func init() {
	SetupLogger(Config{
		LogLevel: "info",
		Format:   "text",
	})
}

type LgrWrapper struct {
	ctx    context.Context
	fields Fields
	logger *logrus.Entry
}

func WithContext(ctx context.Context) LgrWrapper {
	fields := getSourceInfoFields(3)
	fields["root"] = true

	addFieldsFromContext(ctx, fields)

	return WithFields(fields)
}

func WithFields(fields map[string]interface{}) LgrWrapper {
	srcInfo := getSourceInfoFields(3)
	if len(fields) != 0 {
		for key, val := range fields {
			srcInfo[key] = val
		}
	}
	withFields := logger.WithFields(srcInfo)

	return LgrWrapper{
		ctx:    nil,
		fields: fields,
		logger: withFields,
	}
}

func WithRequest(r *http.Request) LgrWrapper {
	return LgrWrapper{
		ctx:    r.Context(),
		fields: nil,
		logger: withRequest(r),
	}
}

func (l LgrWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l LgrWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l LgrWrapper) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l LgrWrapper) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l LgrWrapper) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l LgrWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l LgrWrapper) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l LgrWrapper) WithContext(ctx context.Context) LgrWrapper {
	fields := getSourceInfoFields(3)
	fields["root"] = true

	addFieldsFromContext(ctx, fields)

	return WithFields(fields)
}

func getSourceInfoFields(subtractStackLevels int) map[string]interface{} {
	file, line := getFileInfo(subtractStackLevels)
	m := map[string]interface{}{
		"f": fmt.Sprintf("%s:%d", file, line),
	}
	return m
}

func getFileInfo(subtractStackLevels int) (string, int) {
	_, file, line, _ := runtime.Caller(subtractStackLevels)
	return chopPath(file), line
}

// return the source filename after the last slash
func chopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i != -1 {
		return original[i+1:]
	}
	return original
}

func addFieldsFromContext(ctx context.Context, fields map[string]interface{}) {
	if ctx != nil {
		requestID := ctx.Value(constFields.RequestIDKey)
		requestURI := ctx.Value(constFields.RequestURI)
		requestMethod := ctx.Value(constFields.RequestMethod)

		if requestID != nil {
			fields["request_id"] = requestID
		}
		if requestURI != nil {
			fields["request_uri"] = requestURI
		}
		if requestMethod != nil {
			fields["request_method"] = requestMethod
		}
	}
}

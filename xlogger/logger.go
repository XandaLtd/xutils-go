package xlogger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
)

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type DefaultLogger struct {
	log *zap.Logger
}

// NoOpLogger logs nothing
type NoOpLogger struct{}

type Config struct {
	LogOutputTo []string
	LoggErrsTo  []string
}

func NewLogger(level Level, config Config) (Logger, error) {
	if config.LogOutputTo == nil || len(config.LogOutputTo) == 0 {
		config.LogOutputTo = []string{"stdout"}
	}

	if config.LoggErrsTo == nil || len(config.LoggErrsTo) == 0 {
		config.LoggErrsTo = []string{"stderr"}
	}

	logConfig := zap.Config{
		OutputPaths:      config.LogOutputTo,
		ErrorOutputPaths: config.LoggErrsTo,
		Level:            zap.NewAtomicLevelAt(getLevel(level)),
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger := &DefaultLogger{}

	internalLogger, err := logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("error setting up default logger - %w", err)
	}

	logger.log = internalLogger
	return logger, nil
}

func getLevel(level Level) zapcore.Level {
	switch string(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func (l *DefaultLogger) Print(v ...interface{}) {
	l.Info(fmt.Sprintf("%v", v))
}

func (l *DefaultLogger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		l.Info(format)
	} else {
		l.Info(fmt.Sprintf(format, v...))
	}
}

// Debug logs are typically voluminous, and are usually disabled in production
func (l *DefaultLogger) Debug(msg string, tags ...zap.Field) {
	l.log.Debug(msg, tags...)
}

// Info is the default logging priority.
func (l *DefaultLogger) Info(msg string, tags ...zap.Field) {
	l.log.Info(msg, tags...)
}

// Warning logs are more important than Info, but don't need individual
// human review.
func (l *DefaultLogger) Warning(msg string, tags ...zap.Field) {
	l.log.Debug(msg, tags...)
}

// Error logs are high-priority. If an application is running smoothly,
// it shouldn't generate any error-level logs.
func (l *DefaultLogger) Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Error(msg, tags...)
}

// Panic logs a message, then panics.
func (l *DefaultLogger) Panic(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Panic(msg, tags...)
}

// Fatal logs a message, then calls os.Exit(1).
func (l *DefaultLogger) Fatal(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.log.Fatal(msg, tags...)
}

func (n NoOpLogger) Print(v ...interface{}) {}

func (n NoOpLogger) Printf(format string, v ...interface{}) {}

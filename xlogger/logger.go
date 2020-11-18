package xlogger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = string

const (
	DebugLevel   Level = "debug"
	InfoLevel    Level = "info"
	WarningLevel Level = "warning"
	ErrorLevel   Level = "error"
)

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Debug(msg string, tags ...zap.Field)
	Info(msg string, tags ...zap.Field)
	Warning(msg string, tags ...zap.Field)
	Error(msg string, err error, tags ...zap.Field)
	Panic(msg string, err error, tags ...zap.Field)
	Fatal(msg string, err error, tags ...zap.Field)
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

// NewLogger creates a new logger which logs to the sink(s) specified in config
// and provides structured logging for the metadata
func NewLogger(level Level, config Config, withMeta map[string]interface{}) (Logger, error) {
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
		InitialFields: withMeta,
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
	switch level {
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

// Print prints an info level log
func (l *DefaultLogger) Print(v ...interface{}) {
	l.Info(fmt.Sprintf("%v", v))
}

// Printf prints a formatted info level log
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

// ChildLoggerWithFields returns a new logger with specified structured fields
// The child logger has no impact on the parent nor the parent on the child
func (l *DefaultLogger) ChildLoggerWithFields(fields map[string]interface{}) Logger {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &DefaultLogger{log: l.log.With(zapFields...)}
}

// Unwrap returns the underlying zap logger
func (l *DefaultLogger) Unwrap() *zap.Logger {
	return l.log
}

// Print NoOp
func (n NoOpLogger) Print(v ...interface{}) {}

// Printf NoOp
func (n NoOpLogger) Printf(format string, v ...interface{}) {}

// Debug NoOp
func (n NoOpLogger) Debug(msg string, tags ...zap.Field) {}

// Info NoOp
func (n NoOpLogger) Info(msg string, tags ...zap.Field) {}

// Warning NoOp
func (n NoOpLogger) Warning(msg string, tags ...zap.Field) {}

// Error NoOp
func (n NoOpLogger) Error(msg string, err error, tags ...zap.Field) {}

// Panic NoOp
func (n NoOpLogger) Panic(msg string, err error, tags ...zap.Field) {}

// Fatal NoOp
func (n NoOpLogger) Fatal(msg string, err error, tags ...zap.Field) {}

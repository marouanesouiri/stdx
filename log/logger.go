package log

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
)

// Logger defines the universal logging interface for the stdx library.
// It supports structured logging, leveled logging, and context-aware logging.
type Logger interface {
	// Info logs a message at Info level.
	Info(msg string)
	// Infof logs a formatted message at Info level.
	Infof(format string, args ...any)
	// InfoContext logs a message at Info level with context.
	InfoContext(ctx context.Context, msg string)
	// InfoContextf logs a formatted message at Info level with context.
	InfoContextf(ctx context.Context, format string, args ...any)

	// Debug logs a message at Debug level.
	Debug(msg string)
	// Debugf logs a formatted message at Debug level.
	Debugf(format string, args ...any)
	// DebugContext logs a message at Debug level with context.
	DebugContext(ctx context.Context, msg string)
	// DebugContextf logs a formatted message at Debug level with context.
	DebugContextf(ctx context.Context, format string, args ...any)

	// Warn logs a message at Warn level.
	Warn(msg string)
	// Warnf logs a formatted message at Warn level.
	Warnf(format string, args ...any)
	// WarnContext logs a message at Warn level with context.
	WarnContext(ctx context.Context, msg string)
	// WarnContextf logs a formatted message at Warn level with context.
	WarnContextf(ctx context.Context, format string, args ...any)

	// Error logs a message at Error level.
	Error(msg string)
	// Errorf logs a formatted message at Error level.
	Errorf(format string, args ...any)
	// ErrorContext logs a message at Error level with context.
	ErrorContext(ctx context.Context, msg string)
	// ErrorContextf logs a formatted message at Error level with context.
	ErrorContextf(ctx context.Context, format string, args ...any)

	// Fatal logs a message at Fatal level and then exits the application with status 1.
	Fatal(msg string)
	// Fatalf logs a formatted message at Fatal level and then exits.
	Fatalf(format string, args ...any)
	// FatalContext logs a message at Fatal level with context and then exits.
	FatalContext(ctx context.Context, msg string)
	// FatalContextf logs a formatted message at Fatal level with context and then exits.
	FatalContextf(ctx context.Context, format string, args ...any)

	// WithField adds a single field to the logger context.
	// It returns a new Logger instance with the field added.
	WithField(key string, value any) Logger
	// WithFields adds multiple fields to the logger context.
	// It returns a new Logger instance with the fields added.
	WithFields(fields map[string]any) Logger

	// SetLevel sets the minimum logging level.
	SetLevel(level LogLevel)
	// GetLevel returns the current logging level.
	GetLevel() LogLevel
}

// LogLevel defines the severity level of a log message.
type LogLevel int32

const (
	LogLevelDebugLevel LogLevel = iota
	LogLevelInfoLevel
	LogLevelWarnLevel
	LogLevelErrorLevel
	LogLevelFatalLevel
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelDebugLevel:
		return "debug"
	case LogLevelInfoLevel:
		return "info"
	case LogLevelWarnLevel:
		return "warn"
	case LogLevelErrorLevel:
		return "error"
	case LogLevelFatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

type logCore struct {
	mu    sync.Mutex
	out   io.Writer
	level atomic.Int32
}

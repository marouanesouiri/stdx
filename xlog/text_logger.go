package xlog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"os"
	"sort"
	"strings"
)

// ANSI Color Codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

// TextLogger implements the Logger interface and outputs logs in a human-readable text format.
// It supports ANSI colors for different log levels.
type TextLogger struct {
	core   *logCore
	fields map[string]any
}

var _ Logger = (*TextLogger)(nil)

// NewTextLogger creates a new TextLogger writing to the provided io.Writer at the specified level.
// If out is nil, it defaults to os.Stdout.
func NewTextLogger(out io.Writer, level LogLevel) TextLogger {
	if out == nil {
		out = os.Stdout
	}
	l := TextLogger{
		core: &logCore{
			out: out,
		},
		fields: nil,
	}
	l.core.level.Store(int32(level))
	return l
}

func (l TextLogger) CheckLevel(level LogLevel) bool {
	return LogLevel(l.core.level.Load()) <= level
}

func (l TextLogger) SetLevel(level LogLevel) {
	l.core.level.Store(int32(level))
}

func (l TextLogger) GetLevel() LogLevel {
	return LogLevel(l.core.level.Load())
}

func (l TextLogger) WithField(key string, value any) Logger {
	return l.WithFields(map[string]any{key: value})
}

func (l TextLogger) WithFields(fields map[string]any) Logger {
	if len(fields) == 0 {
		return l
	}

	newFields := make(map[string]any, len(l.fields)+len(fields))
	maps.Copy(newFields, l.fields)
	maps.Copy(newFields, fields)

	return TextLogger{
		core:   l.core,
		fields: newFields,
	}
}

func (l TextLogger) log(ctx context.Context, level LogLevel, msg string) {
	if !l.CheckLevel(level) {
		return
	}

	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return
		}
	}

	var buf bytes.Buffer

	var color string
	switch level {
	case LogLevelDebugLevel:
		color = colorCyan
	case LogLevelInfoLevel:
		color = colorYellow
	case LogLevelWarnLevel:
		color = colorPurple
	case LogLevelErrorLevel:
		color = colorRed
	case LogLevelFatalLevel:
		color = colorRed
	}

	buf.WriteString(color)
	buf.WriteString(strings.ToUpper(level.String()))
	buf.WriteString(colorReset)
	buf.WriteByte(' ')

	buf.WriteString(msg)

	if len(l.fields) > 0 {
		keys := make([]string, 0, len(l.fields))
		for k := range l.fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := l.fields[k]
			buf.WriteByte(' ')
			buf.WriteString(k)
			buf.WriteByte('=')
			l.appendValue(&buf, v)
		}
	}

	buf.WriteByte('\n')

	l.core.mu.Lock()
	l.core.out.Write(buf.Bytes())
	l.core.mu.Unlock()

	if level == LogLevelFatalLevel {
		os.Exit(1)
	}
}

func (l TextLogger) appendValue(buf *bytes.Buffer, v any) {
	switch val := v.(type) {
	case string:
		if strings.ContainsAny(val, " \t\n\r\"=") {
			buf.WriteString(fmt.Sprintf("%q", val))
		} else {
			buf.WriteString(val)
		}
	case error:
		buf.WriteString(fmt.Sprintf("%q", val.Error()))
	default:
		s := fmt.Sprint(v)
		if strings.ContainsAny(s, " \t\n\r\"=") {
			buf.WriteString(fmt.Sprintf("%q", s))
		} else {
			buf.WriteString(s)
		}
	}
}

func (l TextLogger) Info(msg string) {
	l.log(context.Background(), LogLevelInfoLevel, msg)
}

func (l TextLogger) Infof(format string, args ...any) {
	l.log(context.Background(), LogLevelInfoLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) InfoContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelInfoLevel, msg)
}

func (l TextLogger) InfoContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelInfoLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) Debug(msg string) {
	l.log(context.Background(), LogLevelDebugLevel, msg)
}

func (l TextLogger) Debugf(format string, args ...any) {
	l.log(context.Background(), LogLevelDebugLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) DebugContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelDebugLevel, msg)
}

func (l TextLogger) DebugContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelDebugLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) Warn(msg string) {
	l.log(context.Background(), LogLevelWarnLevel, msg)
}

func (l TextLogger) Warnf(format string, args ...any) {
	l.log(context.Background(), LogLevelWarnLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) WarnContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelWarnLevel, msg)
}

func (l TextLogger) WarnContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelWarnLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) Error(msg string) {
	l.log(context.Background(), LogLevelErrorLevel, msg)
}

func (l TextLogger) Errorf(format string, args ...any) {
	l.log(context.Background(), LogLevelErrorLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) ErrorContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelErrorLevel, msg)
}

func (l TextLogger) ErrorContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelErrorLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) Fatal(msg string) {
	l.log(context.Background(), LogLevelFatalLevel, msg)
}

func (l TextLogger) Fatalf(format string, args ...any) {
	l.log(context.Background(), LogLevelFatalLevel, fmt.Sprintf(format, args...))
}

func (l TextLogger) FatalContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelFatalLevel, msg)
}

func (l TextLogger) FatalContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelFatalLevel, fmt.Sprintf(format, args...))
}

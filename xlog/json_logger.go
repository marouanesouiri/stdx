package xlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

// JSONLogger implements the Logger interface and outputs logs in JSON format.
type JSONLogger struct {
	core       *logCore
	jsonFields []byte
}

var _ Logger = (*JSONLogger)(nil)

// NewJSONLogger creates a new JSONLogger writing to the provided io.Writer at the specified level.
// If out is nil, it defaults to os.Stdout.
func NewJSONLogger(out io.Writer, level LogLevel) JSONLogger {
	if out == nil {
		out = os.Stdout
	}
	l := JSONLogger{
		core: &logCore{
			out: out,
		},
		jsonFields: nil,
	}
	l.core.level.Store(int32(level))
	return l
}

func (l JSONLogger) CheckLevel(level LogLevel) bool {
	return LogLevel(l.core.level.Load()) <= level
}

func (l JSONLogger) SetLevel(level LogLevel) {
	l.core.level.Store(int32(level))
}

func (l JSONLogger) GetLevel() LogLevel {
	return LogLevel(l.core.level.Load())
}

func (l JSONLogger) WithField(key string, value any) Logger {
	return l.WithFields(map[string]any{key: value})
}

func (l JSONLogger) WithFields(fields map[string]any) Logger {
	if len(fields) == 0 {
		return l
	}

	b, err := json.Marshal(fields)
	if err != nil {
		return l
	}

	if len(b) > 2 {
		b = b[1 : len(b)-1]
	} else {
		return l
	}

	var newJsonFields []byte
	if len(l.jsonFields) == 0 {
		newJsonFields = b
	} else {
		newJsonFields = make([]byte, len(l.jsonFields)+1+len(b))
		copy(newJsonFields, l.jsonFields)
		newJsonFields[len(l.jsonFields)] = ','
		copy(newJsonFields[len(l.jsonFields)+1:], b)
	}

	return JSONLogger{
		core:       l.core,
		jsonFields: newJsonFields,
	}
}

func (l JSONLogger) log(ctx context.Context, level LogLevel, msg string) {
	if !l.CheckLevel(level) {
		return
	}

	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return
		}
	}

	var buf bytes.Buffer

	buf.WriteString(`{"time":"`)
	buf.WriteString(time.Now().Format(time.RFC3339))
	buf.WriteString(`","level":"`)
	buf.WriteString(level.String())
	buf.WriteString(`","msg":`)
	msgBytes, _ := json.Marshal(msg)
	buf.Write(msgBytes)

	_, file, line, ok := runtime.Caller(2)
	if ok {
		buf.WriteString(`,"source":"`)
		buf.WriteString(file)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(line))
		buf.WriteByte('"')
	}

	if len(l.jsonFields) > 0 {
		buf.WriteByte(',')
		buf.Write(l.jsonFields)
	}

	buf.WriteString("}\n")

	l.core.mu.Lock()
	l.core.out.Write(buf.Bytes())
	l.core.mu.Unlock()

	if level == LogLevelFatalLevel {
		os.Exit(1)
	}
}

func (l JSONLogger) Info(msg string) {
	l.log(context.Background(), LogLevelInfoLevel, msg)
}

func (l JSONLogger) Infof(format string, args ...any) {
	l.log(context.Background(), LogLevelInfoLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) InfoContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelInfoLevel, msg)
}

func (l JSONLogger) InfoContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelInfoLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) Debug(msg string) {
	l.log(context.Background(), LogLevelDebugLevel, msg)
}

func (l JSONLogger) Debugf(format string, args ...any) {
	l.log(context.Background(), LogLevelDebugLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) DebugContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelDebugLevel, msg)
}

func (l JSONLogger) DebugContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelDebugLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) Warn(msg string) {
	l.log(context.Background(), LogLevelWarnLevel, msg)
}

func (l JSONLogger) Warnf(format string, args ...any) {
	l.log(context.Background(), LogLevelWarnLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) WarnContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelWarnLevel, msg)
}

func (l JSONLogger) WarnContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelWarnLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) Error(msg string) {
	l.log(context.Background(), LogLevelErrorLevel, msg)
}

func (l JSONLogger) Errorf(format string, args ...any) {
	l.log(context.Background(), LogLevelErrorLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) ErrorContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelErrorLevel, msg)
}

func (l JSONLogger) ErrorContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelErrorLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) Fatal(msg string) {
	l.log(context.Background(), LogLevelFatalLevel, msg)
}

func (l JSONLogger) Fatalf(format string, args ...any) {
	l.log(context.Background(), LogLevelFatalLevel, fmt.Sprintf(format, args...))
}

func (l JSONLogger) FatalContext(ctx context.Context, msg string) {
	l.log(ctx, LogLevelFatalLevel, msg)
}

func (l JSONLogger) FatalContextf(ctx context.Context, format string, args ...any) {
	l.log(ctx, LogLevelFatalLevel, fmt.Sprintf(format, args...))
}

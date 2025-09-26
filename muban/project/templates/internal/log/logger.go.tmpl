package log

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Logger 统一的日志记录接口
type Logger interface {
	Debug(msg string, attrs ...slog.Attr)
	Info(msg string, attrs ...slog.Attr)
	Warn(msg string, attrs ...slog.Attr)
	Error(msg string, attrs ...slog.Attr)
	Fatal(msg string, attrs ...slog.Attr)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	With(attrs ...slog.Attr) Logger
	WithGroup(name string) Logger
}

// Sink 日志输出目标抽象
type Sink interface {
	Write(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) error
	Close() error
}

// MultiSink 多输出目标
type MultiSink struct {
	sinks []Sink
}

func NewMultiSink(sinks ...Sink) *MultiSink {
	return &MultiSink{sinks: sinks}
}

func (m *MultiSink) Write(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) error {
	for _, sink := range m.sinks {
		if err := sink.Write(ctx, level, msg, attrs); err != nil {
			// 记录错误但不中断其他sink
			continue
		}
	}
	return nil
}

func (m *MultiSink) Close() error {
	for _, sink := range m.sinks {
		sink.Close()
	}
	return nil
}

// DefaultLogger 默认日志记录器实现
type DefaultLogger struct {
	slog *slog.Logger
	sink Sink
	mu   sync.RWMutex
}

func NewDefaultLogger(sink Sink, level slog.Level) *DefaultLogger {
	handler := &SinkHandler{sink: sink, level: level}
	return &DefaultLogger{
		slog: slog.New(handler),
		sink: sink,
	}
}

func (l *DefaultLogger) Debug(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

func (l *DefaultLogger) Info(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

func (l *DefaultLogger) Warn(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

func (l *DefaultLogger) Error(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

func (l *DefaultLogger) Fatal(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.Background(), slog.LevelError+1, msg, attrs...)
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	l.slog.LogAttrs(context.Background(), slog.LevelDebug, format, slog.String("args", fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	l.slog.LogAttrs(context.Background(), slog.LevelInfo, format, slog.String("args", fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	l.slog.LogAttrs(context.Background(), slog.LevelWarn, format, slog.String("args", fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	l.slog.LogAttrs(context.Background(), slog.LevelError, format, slog.String("args", fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) Fatalf(format string, args ...interface{}) {
	l.slog.LogAttrs(context.Background(), slog.LevelError+1, format, slog.String("args", fmt.Sprintf(format, args...)))
}

func (l *DefaultLogger) With(attrs ...slog.Attr) Logger {
	// 将slog.Attr转换为interface{}
	args := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value.Any())
	}
	return &DefaultLogger{
		slog: l.slog.With(args...),
		sink: l.sink,
	}
}

func (l *DefaultLogger) WithGroup(name string) Logger {
	return &DefaultLogger{
		slog: l.slog.WithGroup(name),
		sink: l.sink,
	}
}

// SinkHandler slog.Handler 实现
type SinkHandler struct {
	sink  Sink
	level slog.Level
}

func (h *SinkHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *SinkHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	return h.sink.Write(ctx, r.Level, r.Message, attrs)
}

func (h *SinkHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SinkHandler{sink: h.sink, level: h.level}
}

func (h *SinkHandler) WithGroup(name string) slog.Handler {
	return &SinkHandler{sink: h.sink, level: h.level}
}

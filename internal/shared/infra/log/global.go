package log

import (
	"log/slog"
	"sync"
)

var (
	globalLogger Logger
	mu           sync.RWMutex
)

// SetGlobalLogger 设置全局日志记录器
func SetGlobalLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = logger
}

// GetGlobalLogger 获取全局日志记录器
func GetGlobalLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	return globalLogger
}

// 全局日志函数
func Debug(msg string, attrs ...slog.Attr) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Debug(msg, attrs...)
	}
}

func Info(msg string, attrs ...slog.Attr) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Info(msg, attrs...)
	}
}

func Warn(msg string, attrs ...slog.Attr) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Warn(msg, attrs...)
	}
}

func Error(msg string, attrs ...slog.Attr) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Error(msg, attrs...)
	}
}

func Fatal(msg string, attrs ...slog.Attr) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Fatal(msg, attrs...)
	}
}

func Debugf(format string, args ...interface{}) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Fatalf(format, args...)
	}
}

func With(attrs ...slog.Attr) Logger {
	if logger := GetGlobalLogger(); logger != nil {
		return logger.With(attrs...)
	}
	return nil
}

func WithGroup(name string) Logger {
	if logger := GetGlobalLogger(); logger != nil {
		return logger.WithGroup(name)
	}
	return nil
}

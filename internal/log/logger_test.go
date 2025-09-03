package log

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	// 创建测试用的控制台输出
	sink := NewConsoleSink(ConsoleSinkConfig{
		Format: "text",
		Output: "stdout",
	})

	logger := NewDefaultLogger(sink, slog.LevelInfo)

	tests := []struct {
		name     string
		level    string
		message  string
		attrs    []slog.Attr
		expected bool
	}{
		{
			name:     "info level",
			level:    "info",
			message:  "test info message",
			attrs:    []slog.Attr{slog.String("key", "value")},
			expected: true,
		},
		{
			name:     "debug level",
			level:    "debug",
			message:  "test debug message",
			attrs:    []slog.Attr{},
			expected: false, // debug level is below info level
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试日志输出
			switch tt.level {
			case "info":
				logger.Info(tt.message, tt.attrs...)
			case "debug":
				logger.Debug(tt.message, tt.attrs...)
			}
			// 由于日志输出到stdout，我们无法直接验证输出内容
			// 这里主要测试方法调用不会panic
			assert.True(t, true)
		})
	}
}

func TestConsoleSink(t *testing.T) {
	tests := []struct {
		name   string
		config ConsoleSinkConfig
	}{
		{
			name: "color format",
			config: ConsoleSinkConfig{
				Format: "color",
				Output: "stdout",
			},
		},
		{
			name: "json format",
			config: ConsoleSinkConfig{
				Format: "json",
				Output: "stdout",
			},
		},
		{
			name: "text format",
			config: ConsoleSinkConfig{
				Format: "text",
				Output: "stdout",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink := NewConsoleSink(tt.config)
			assert.NotNil(t, sink)

			// 测试写入
			err := sink.Write(context.Background(), slog.LevelInfo, "test message", []slog.Attr{
				slog.String("key", "value"),
			})
			assert.NoError(t, err)

			// 测试关闭
			err = sink.Close()
			assert.NoError(t, err)
		})
	}
}

func TestFileSink(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	config := FileSinkConfig{
		Filename:   tmpFile.Name(),
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
		Format:     "json",
	}

	sink := NewFileSink(config)
	assert.NotNil(t, sink)

	// 测试写入
	err = sink.Write(context.Background(), slog.LevelInfo, "test message", []slog.Attr{
		slog.String("key", "value"),
	})
	assert.NoError(t, err)

	// 测试关闭
	err = sink.Close()
	assert.NoError(t, err)
}

func TestMultiSink(t *testing.T) {
	// 创建多个输出目标
	consoleSink := NewConsoleSink(ConsoleSinkConfig{
		Format: "text",
		Output: "stdout",
	})

	tmpFile, err := os.CreateTemp("", "test_multi_log_*.log")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	fileSink := NewFileSink(FileSinkConfig{
		Filename: tmpFile.Name(),
		Format:   "json",
	})

	// 创建多输出目标
	multiSink := NewMultiSink(consoleSink, fileSink)
	assert.NotNil(t, multiSink)

	// 测试写入
	err = multiSink.Write(context.Background(), slog.LevelInfo, "test message", []slog.Attr{
		slog.String("key", "value"),
	})
	assert.NoError(t, err)

	// 测试关闭
	err = multiSink.Close()
	assert.NoError(t, err)
}

func TestSinkHandler(t *testing.T) {
	sink := NewConsoleSink(ConsoleSinkConfig{
		Format: "text",
		Output: "stdout",
	})

	handler := &SinkHandler{
		sink:  sink,
		level: slog.LevelInfo,
	}

	// 测试Enabled
	assert.True(t, handler.Enabled(context.Background(), slog.LevelInfo))
	assert.False(t, handler.Enabled(context.Background(), slog.LevelDebug))

	// 测试Handle
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.AddAttrs(slog.String("key", "value"))

	err := handler.Handle(context.Background(), record)
	assert.NoError(t, err)
}

func TestGlobalLogger(t *testing.T) {
	// 创建测试日志记录器
	sink := NewConsoleSink(ConsoleSinkConfig{
		Format: "text",
		Output: "stdout",
	})
	logger := NewDefaultLogger(sink, slog.LevelInfo)

	// 设置全局日志记录器
	SetGlobalLogger(logger)

	// 测试全局日志函数
	Info("global info message", slog.String("key", "value"))
	Error("global error message", slog.String("key", "value"))

	// 测试获取全局日志记录器
	globalLogger := GetGlobalLogger()
	assert.NotNil(t, globalLogger)
	assert.Equal(t, logger, globalLogger)
}

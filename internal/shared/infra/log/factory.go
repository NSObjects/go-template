package log

import (
	"log/slog"
	"strings"

	"github.com/NSObjects/go-template/internal/configs"
)

// LogConfig 日志配置
type LogConfig struct {
	Level  string `json:"level" yaml:"level" toml:"level"`
	Format string `json:"format" yaml:"format" toml:"format"`

	Console ConsoleSinkConfig `json:"console" yaml:"console" toml:"console"`
	File    FileSinkConfig    `json:"file" yaml:"file" toml:"file"`

	Elasticsearch ElasticsearchSinkConfig `json:"elasticsearch" yaml:"elasticsearch" toml:"elasticsearch"`
	Loki          LokiSinkConfig          `json:"loki" yaml:"loki" toml:"loki"`
}

// NewLogger 根据配置创建日志记录器
func NewLogger(cfg configs.Config) Logger {
	logCfg := cfg.Log

	// 解析日志级别
	level := parseLevel(logCfg.Level)

	// 创建输出目标
	var sinks []Sink

	// 根据环境调整日志配置
	env := cfg.System.Env
	if env == "" {
		env = "dev" // 默认开发环境
	}

	// 控制台输出（默认启用）
	if logCfg.Console.Format != "" || logCfg.Console.Output != "" {
		sinks = append(sinks, NewConsoleSink(ConsoleSinkConfig{
			Format: logCfg.Console.Format,
			Output: logCfg.Console.Output,
		}))
	} else {
		// 根据环境设置默认格式
		format := "color"
		if env == "prod" {
			format = "json"
		}
		sinks = append(sinks, NewConsoleSink(ConsoleSinkConfig{
			Format: format,
			Output: "stdout",
		}))
	}

	// 文件输出 - 生产环境默认启用
	if logCfg.File.Filename != "" && (env == "prod" || env == "test") {
		sinks = append(sinks, NewFileSink(FileSinkConfig{
			Filename:   logCfg.File.Filename,
			MaxSize:    logCfg.File.MaxSize,
			MaxBackups: logCfg.File.MaxBackups,
			MaxAge:     logCfg.File.MaxAge,
			Compress:   logCfg.File.Compress,
			Format:     logCfg.File.Format,
		}))
	}

	// Elasticsearch输出
	if logCfg.Elasticsearch.URL != "" {
		sinks = append(sinks, NewElasticsearchSink(ElasticsearchSinkConfig{
			URL:     logCfg.Elasticsearch.URL,
			Index:   logCfg.Elasticsearch.Index,
			Timeout: logCfg.Elasticsearch.Timeout,
		}))
	}

	// Loki输出
	if logCfg.Loki.URL != "" {
		sinks = append(sinks, NewLokiSink(LokiSinkConfig{
			URL:     logCfg.Loki.URL,
			Labels:  logCfg.Loki.Labels,
			Timeout: logCfg.Loki.Timeout,
		}))
	}

	// 创建多输出目标
	var sink Sink
	if len(sinks) == 1 {
		sink = sinks[0]
	} else {
		sink = NewMultiSink(sinks...)
	}

	// 创建日志记录器
	logger := NewDefaultLogger(sink, level)

	// 设置全局日志记录器
	SetGlobalLogger(logger)

	return logger
}

// parseLevel 解析日志级别
func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

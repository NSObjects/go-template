package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// FileSink 文件输出
type FileSink struct {
	writer *lumberjack.Logger
	format string // "json" | "text"
}

// FileSinkConfig 文件输出配置
type FileSinkConfig struct {
	Filename   string `json:"filename" yaml:"filename" toml:"filename"`
	MaxSize    int    `json:"max_size" yaml:"max_size" toml:"max_size"` // MB
	MaxBackups int    `json:"max_backups" yaml:"max_backups" toml:"max_backups"`
	MaxAge     int    `json:"max_age" yaml:"max_age" toml:"max_age"` // days
	Compress   bool   `json:"compress" yaml:"compress" toml:"compress"`
	Format     string `json:"format" yaml:"format" toml:"format"` // json, text
}

func NewFileSink(cfg FileSinkConfig) *FileSink {
	format := cfg.Format
	if format == "" {
		format = "json"
	}

	// 确保目录存在
	dir := filepath.Dir(cfg.Filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create log directory: %v", err))
	}

	writer := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	return &FileSink{
		writer: writer,
		format: format,
	}
}

func (f *FileSink) Write(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) error {
	switch f.format {
	case "json":
		return f.writeJSON(level, msg, attrs)
	case "text":
		return f.writeText(level, msg, attrs)
	default:
		return f.writeJSON(level, msg, attrs)
	}
}

func (f *FileSink) writeJSON(level slog.Level, msg string, attrs []slog.Attr) error {
	entry := map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": level.String(),
		"msg":   msg,
	}

	for _, attr := range attrs {
		entry[attr.Key] = attr.Value.Any()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = f.writer.Write(append(data, '\n'))
	return err
}

func (f *FileSink) writeText(level slog.Level, msg string, attrs []slog.Attr) error {
	text := fmt.Sprintf("%s %s %s",
		time.Now().Format("2006-01-02 15:04:05"),
		level.String(),
		msg)

	for _, attr := range attrs {
		text += fmt.Sprintf(" %s=%v", attr.Key, attr.Value.Any())
	}
	text += "\n"

	_, err := f.writer.Write([]byte(text))
	return err
}

func (f *FileSink) Close() error {
	return f.writer.Close()
}

package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

// ConsoleSink 控制台输出
type ConsoleSink struct {
	writer io.Writer
	format string // "json" | "text" | "color"
}

// ConsoleSinkConfig 控制台输出配置
type ConsoleSinkConfig struct {
	Format string `json:"format" yaml:"format" toml:"format"` // json, text, color
	Output string `json:"output" yaml:"output" toml:"output"` // stdout, stderr
}

func NewConsoleSink(cfg ConsoleSinkConfig) *ConsoleSink {
	format := cfg.Format
	if format == "" {
		format = "color"
	}

	writer := os.Stdout
	if cfg.Output == "stderr" {
		writer = os.Stderr
	}

	return &ConsoleSink{
		writer: writer,
		format: format,
	}
}

func (c *ConsoleSink) Write(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) error {
	switch c.format {
	case "json":
		return c.writeJSON(level, msg, attrs)
	case "text":
		return c.writeText(level, msg, attrs)
	case "color":
		return c.writeColor(level, msg, attrs)
	default:
		return c.writeColor(level, msg, attrs)
	}
}

func (c *ConsoleSink) writeJSON(level slog.Level, msg string, attrs []slog.Attr) error {
	// 简单的JSON格式输出
	json := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s"`,
		time.Now().Format(time.RFC3339),
		level.String(),
		msg)

	for _, attr := range attrs {
		json += fmt.Sprintf(`,"%s":"%v"`, attr.Key, attr.Value.Any())
	}
	json += "}\n"

	_, err := c.writer.Write([]byte(json))
	return err
}

func (c *ConsoleSink) writeText(level slog.Level, msg string, attrs []slog.Attr) error {
	// 简单的文本格式输出
	text := fmt.Sprintf("%s %s %s",
		time.Now().Format("2006-01-02 15:04:05"),
		strings.ToUpper(level.String()),
		msg)

	for _, attr := range attrs {
		text += fmt.Sprintf(" %s=%v", attr.Key, attr.Value.Any())
	}
	text += "\n"

	_, err := c.writer.Write([]byte(text))
	return err
}

func (c *ConsoleSink) writeColor(level slog.Level, msg string, attrs []slog.Attr) error {
	// 使用tint进行彩色输出
	handler := tint.NewHandler(c.writer, &tint.Options{
		AddSource:  false,
		TimeFormat: time.DateTime,
		Level:      slog.LevelDebug,
	})

	record := slog.NewRecord(time.Now(), level, msg, 0)
	for _, attr := range attrs {
		record.AddAttrs(attr)
	}

	return handler.Handle(context.Background(), record)
}

func (c *ConsoleSink) Close() error {
	return nil
}

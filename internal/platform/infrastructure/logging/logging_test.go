package logging

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/NSObjects/go-template/internal/platform/configs"
)

func TestFromAppConfigDerivesDefaultsFromSystemLevel(t *testing.T) {
	debugCfg := FromAppConfig(configs.Config{
		System: configs.SystemConfig{Level: configs.DebugLevel},
	})
	if debugCfg.Level != zerolog.DebugLevel {
		t.Fatalf("debug Level = %v, want %v", debugCfg.Level, zerolog.DebugLevel)
	}
	if debugCfg.Format != configs.LogFormatConsole {
		t.Fatalf("debug Format = %q, want %q", debugCfg.Format, configs.LogFormatConsole)
	}

	onlineCfg := FromAppConfig(configs.Config{})
	if onlineCfg.Level != zerolog.InfoLevel {
		t.Fatalf("online Level = %v, want %v", onlineCfg.Level, zerolog.InfoLevel)
	}
	if onlineCfg.Format != configs.LogFormatJSON {
		t.Fatalf("online Format = %q, want %q", onlineCfg.Format, configs.LogFormatJSON)
	}
	if onlineCfg.AppName != configs.DefaultAppName {
		t.Fatalf("online AppName = %q, want %q", onlineCfg.AppName, configs.DefaultAppName)
	}
}

func TestNewLoggerWritesStructuredJSON(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewLogger(&buf, Config{
		Level:      zerolog.InfoLevel,
		Format:     configs.LogFormatJSON,
		Output:     configs.LogOutputStdout,
		AppName:    "billing-api",
		AppVersion: "2026.06.17",
	})
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.Info().Str("request_id", "req-1").Msg("request handled")

	output := buf.String()
	if !strings.Contains(output, `"message":"request handled"`) {
		t.Fatalf("log output = %q, want JSON message", output)
	}
	if !strings.Contains(output, `"request_id":"req-1"`) {
		t.Fatalf("log output = %q, want request_id field", output)
	}
	if !strings.Contains(output, `"app":"billing-api"`) {
		t.Fatalf("log output = %q, want app field", output)
	}
	if !strings.Contains(output, `"version":"2026.06.17"`) {
		t.Fatalf("log output = %q, want version field", output)
	}
}

func TestNewLoggerWritesConsoleOutput(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewLogger(&buf, Config{
		Level:  zerolog.InfoLevel,
		Format: configs.LogFormatConsole,
		Output: configs.LogOutputStdout,
	})
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.Info().Str("request_id", "req-1").Msg("request handled")

	output := buf.String()
	if !strings.Contains(output, "request handled") {
		t.Fatalf("log output = %q, want console message", output)
	}
	if !strings.Contains(output, "request_id=req-1") {
		t.Fatalf("log output = %q, want request_id field", output)
	}
}

func TestNewLoggerUsesZeroValueDefaults(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewLogger(&buf, Config{})
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	logger.Info().Msg("default logger")

	output := buf.String()
	if !strings.Contains(output, `"message":"default logger"`) {
		t.Fatalf("log output = %q, want default JSON message", output)
	}
}

func TestNewLoggerRejectsInvalidFormat(t *testing.T) {
	if _, err := NewLogger(&bytes.Buffer{}, Config{Format: "color"}); err == nil {
		t.Fatal("NewLogger() error = nil, want invalid format error")
	}
}

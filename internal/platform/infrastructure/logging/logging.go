// Package logging installs the process-level zerolog runtime.
package logging

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/NSObjects/go-template/internal/platform/configs"
)

const defaultTimeFormat = time.RFC3339Nano

// Config contains process-level logging infrastructure settings.
type Config struct {
	Level      zerolog.Level
	Format     string
	Output     string
	Caller     bool
	AppName    string
	AppVersion string
}

// FromAppConfig derives logging settings from application config.
func FromAppConfig(cfg configs.Config) Config {
	cfg = configs.Normalize(cfg)

	level := zerolog.InfoLevel
	if cfg.System.Level == configs.DebugLevel {
		level = zerolog.DebugLevel
	}

	return Config{
		Level:      level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		Caller:     cfg.Log.Caller,
		AppName:    cfg.App.Name,
		AppVersion: cfg.App.Version,
	}
}

// Install configures zerolog's process-wide logger and standard-library log output.
func Install(cfg Config) error {
	writer, err := outputWriter(normalizeConfig(cfg).Output)
	if err != nil {
		return err
	}
	logger, err := NewLogger(writer, cfg)
	if err != nil {
		return err
	}

	zlog.Logger = logger
	zerolog.DefaultContextLogger = &zlog.Logger
	zerolog.SetGlobalLevel(normalizeConfig(cfg).Level)
	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog.Logger)
	return nil
}

// NewLogger creates a zerolog logger for the provided writer.
func NewLogger(writer io.Writer, cfg Config) (zerolog.Logger, error) {
	if writer == nil {
		return zerolog.Logger{}, fmt.Errorf("log writer is nil")
	}
	cfg = normalizeConfig(cfg)

	output, err := formatWriter(writer, cfg)
	if err != nil {
		return zerolog.Logger{}, err
	}

	builder := zerolog.New(output).Level(cfg.Level).With().Timestamp()
	if cfg.AppName != "" {
		builder = builder.Str("app", cfg.AppName)
	}
	if cfg.AppVersion != "" {
		builder = builder.Str("version", cfg.AppVersion)
	}
	if cfg.Caller {
		builder = builder.Caller()
	}
	return builder.Logger(), nil
}

// FromContext returns the request-scoped zerolog logger.
func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

func normalizeConfig(cfg Config) Config {
	if cfg == (Config{}) || cfg.Level == zerolog.NoLevel {
		cfg.Level = zerolog.InfoLevel
	}
	if cfg.Format == "" {
		cfg.Format = configs.LogFormatJSON
	}
	if cfg.Format == "text" {
		cfg.Format = configs.LogFormatConsole
	}
	if cfg.Output == "" {
		cfg.Output = configs.LogOutputStdout
	}
	return cfg
}

func formatWriter(writer io.Writer, cfg Config) (io.Writer, error) {
	switch cfg.Format {
	case configs.LogFormatConsole:
		return zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: defaultTimeFormat,
		}, nil
	case configs.LogFormatJSON:
		return writer, nil
	default:
		return nil, fmt.Errorf("unsupported log format %q", cfg.Format)
	}
}

func outputWriter(output string) (io.Writer, error) {
	switch output {
	case configs.LogOutputStdout:
		return os.Stdout, nil
	case configs.LogOutputStderr:
		return os.Stderr, nil
	default:
		return nil, fmt.Errorf("unsupported log output %q", output)
	}
}

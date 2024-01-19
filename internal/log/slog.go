/*
 *
 * slog.go
 * log
 *
 * Created by lintao on 2023/12/5 09:57
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/lmittmann/tint"
)

type log struct {
	logger *slog.Logger
	level  slog.Level
}

func New(cfg configs.Config) log {

	l := log{
		level:  slog.Level(cfg.Log.Level),
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})),
	}

	return l
}

var logger *slog.Logger

func defaultLog() *slog.Logger {
	if logger != nil {
		return logger
	}

	w := os.Stdout
	logger = slog.New(tint.NewHandler(w, &tint.Options{
		AddSource:  true,
		TimeFormat: time.DateTime,
		Level:      slog.LevelDebug,
	}))

	return logger
}

func Info(format string, args ...slog.Attr) {
	defaultLog().LogAttrs(context.Background(), slog.LevelInfo, format, args...)
}

func Error(err error, args ...slog.Attr) {
	defaultLog().LogAttrs(context.Background(), slog.LevelError, fmt.Sprintf("%+v", err), args...)
}

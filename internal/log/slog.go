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
	"fmt"
	"log/slog"

	"github.com/NSObjects/echo-admin/internal/configs"
)

// New 创建日志记录器（兼容旧接口）
func New(cfg configs.Config) Logger {
	return NewLogger(cfg)
}

// InfoCompat 兼容旧接口
func InfoCompat(format string, args ...slog.Attr) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Info(format, args...)
	}
}

// ErrorCompat 兼容旧接口
func ErrorCompat(err error, args ...slog.Attr) {
	if logger := GetGlobalLogger(); logger != nil {
		logger.Error(fmt.Sprintf("%+v", err), args...)
	}
}

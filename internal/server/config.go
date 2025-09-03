/*
 * Server Configuration
 * 服务器配置管理
 *
 * Created by lintao on 2024/1/4
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package server

import (
	"time"

	"github.com/NSObjects/echo-admin/internal/configs"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	// 端口
	Port string
	// 读取超时
	ReadTimeout time.Duration
	// 写入超时
	WriteTimeout time.Duration
	// 空闲超时
	IdleTimeout time.Duration
	// 关闭超时
	ShutdownTimeout time.Duration
	// 是否隐藏Banner
	HideBanner bool
	// 是否启用调试模式
	Debug bool
}

// DefaultServerConfig 默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:            ":8080",
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		HideBanner:      true,
		Debug:           false,
	}
}

// FromAppConfig 从应用配置创建服务器配置
func FromAppConfig(cfg configs.Config) *ServerConfig {
	return &ServerConfig{
		Port:            cfg.System.Port,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		IdleTimeout:     120 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		HideBanner:      true,
		Debug:           cfg.System.Level == 1, // 1 = debug, 2 = online
	}
}

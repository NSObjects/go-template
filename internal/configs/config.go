/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package configs

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota + 1
	// OnlineLevel is the default production priority.
	OnlineLevel
)

type Config struct {
	System SystemConfig `mapstructure:"system"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type SystemConfig struct {
	Port  string `mapstructure:"port"`
	Level Level  `mapstructure:"level"`
	Env   string `mapstructure:"env"`
}

type JWTConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Secret    string   `mapstructure:"secret"`
	Expire    int      `mapstructure:"expire"`
	SkipPaths []string `mapstructure:"skip_paths"`
}

package configs

import (
	configs "github.com/NSObjects/go-kit/config"
)

// AppConfig 应用全局配置
type AppConfig struct {
	configs.Config `mapstructure:",squash"`
	Business       BusinessConfig `mapstructure:"business"`
}

// BusinessConfig 业务专用配置
type BusinessConfig struct {
	// 在这里添加业务配置项
	// Example:
	// FeatureXEnabled bool `mapstructure:"feature_x_enabled"`
	Hello string `mapstructure:"hello"`
}

// Package configs loads and validates static application configuration.
package configs

// Level identifies the runtime logging and debug mode.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota + 1
	// OnlineLevel is the default production priority.
	OnlineLevel
)

const (
	// DefaultPort is the HTTP port used when config omits system.port.
	DefaultPort = ":9322"

	// DefaultEnv is the runtime environment name used when config omits
	// system.env.
	DefaultEnv = "dev"
)

// Config is the complete application configuration loaded at startup.
type Config struct {
	System SystemConfig `mapstructure:"system"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

// SystemConfig controls process-level runtime settings.
type SystemConfig struct {
	Port  string `mapstructure:"port"`
	Level Level  `mapstructure:"level"`
	Env   string `mapstructure:"env"`
}

// JWTConfig controls optional server-level JWT verification.
type JWTConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Secret    string   `mapstructure:"secret"`
	SkipPaths []string `mapstructure:"skip_paths"`
}

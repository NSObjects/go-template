package configs

import (
	"bytes"
	"errors"
	"strings"

	"github.com/spf13/viper"
)

func decodeConfigWithEnv(data []byte, format string, useEnv bool) (Config, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType(configType(format))
	if useEnv {
		v.SetEnvPrefix("GO_TEMPLATE")
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer("::", "_", ".", "_"))
	}
	if err := v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	cfg = Normalize(cfg)
	if err := Validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Normalize applies application config defaults.
func Normalize(cfg Config) Config {
	if cfg.System.Port == "" {
		cfg.System.Port = DefaultPort
	}
	if cfg.System.Level == 0 {
		cfg.System.Level = OnlineLevel
	}
	if cfg.System.Env == "" {
		cfg.System.Env = DefaultEnv
	}
	if cfg.JWT.SkipPaths == nil {
		cfg.JWT.SkipPaths = []string{"/api/health", "/api/info"}
	}
	return cfg
}

// Validate checks cross-field configuration rules that must fail before the
// HTTP runtime starts.
func Validate(cfg Config) error {
	if cfg.JWT.Enabled && strings.TrimSpace(cfg.JWT.Secret) == "" {
		return errors.New("jwt secret is required when jwt is enabled")
	}
	return nil
}

func configType(format string) string {
	switch strings.ToLower(format) {
	case "json":
		return "json"
	case "yaml", "yml":
		return "yaml"
	default:
		return "toml"
	}
}

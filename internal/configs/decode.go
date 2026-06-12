package configs

import (
	"bytes"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func decodeConfig(data []byte, format string) (Config, error) {
	return decodeConfigWithEnv(data, format, false)
}

func decodeConfigWithEnv(data []byte, format string, useEnv bool) (Config, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType(configType(format))
	if useEnv {
		v.SetEnvPrefix("ECHOADMIN")
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
	if useEnv {
		applyCapabilityEnvOverrides(&cfg)
	}
	return cfg, nil
}

func applyCapabilityEnvOverrides(cfg *Config) {
	if provider := strings.TrimSpace(os.Getenv("ECHOADMIN_USER_STORAGE_PROVIDER")); provider != "" {
		if cfg.Capabilities.Providers == nil {
			cfg.Capabilities.Providers = make(map[string]string, 1)
		}
		cfg.Capabilities.Providers["user.storage"] = provider
	}
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

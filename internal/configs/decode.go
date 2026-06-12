package configs

import (
	"bytes"
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
	cfg = Normalize(cfg)
	if useEnv {
		cfg = Normalize(cfg)
	}
	return cfg, nil
}

// Normalize applies platform-level config defaults.
func Normalize(cfg Config) Config {
	return cfg
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

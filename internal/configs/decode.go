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
	cfg = Normalize(cfg)
	if useEnv {
		applyCapabilityEnvOverrides(&cfg)
		cfg = Normalize(cfg)
	}
	return cfg, nil
}

// Normalize derives platform selections from business-facing configuration.
func Normalize(cfg Config) Config {
	normalizeCapabilityProviderSelections(&cfg)
	return cfg
}

func normalizeCapabilityProviderSelections(cfg *Config) {
	if provider := strings.TrimSpace(cfg.User.Storage.Provider); provider != "" {
		setCapabilityProvider(cfg, "user.storage", provider)
	}
}

func applyCapabilityEnvOverrides(cfg *Config) {
	if provider := strings.TrimSpace(os.Getenv("ECHOADMIN_USER_STORAGE_PROVIDER")); provider != "" {
		cfg.User.Storage.Provider = provider
		setCapabilityProvider(cfg, "user.storage", provider)
	}
}

func setCapabilityProvider(cfg *Config, capability, provider string) {
	if cfg.Capabilities.Providers == nil {
		cfg.Capabilities.Providers = make(map[string]string, 1)
	}
	cfg.Capabilities.Providers[capability] = provider
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

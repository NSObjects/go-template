package configs

import (
	"os"
	"strings"
)

// Load reads a static application config file and applies supported environment
// variable overrides.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return decodeConfigWithEnv(data, formatFromPath(path), true)
}

func formatFromPath(path string) string {
	if dot := strings.LastIndex(path, "."); dot >= 0 {
		return strings.ToLower(path[dot+1:])
	}
	return "toml"
}

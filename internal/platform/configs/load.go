package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Load reads a static application config file and applies supported environment
// variable overrides.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	format, err := formatFromPath(path)
	if err != nil {
		return Config{}, err
	}
	return decodeConfigWithEnv(data, format, true)
}

func formatFromPath(path string) (string, error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return configFormatTOML, nil
	}

	format := strings.ToLower(strings.TrimPrefix(ext, "."))
	switch format {
	case configFormatTOML, configFormatJSON, configFormatYAML, configFormatYML:
		return format, nil
	default:
		return "", fmt.Errorf("unsupported config file extension %q", format)
	}
}

package configs

import (
	"fmt"
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
	format, err := formatFromPath(path)
	if err != nil {
		return Config{}, err
	}
	return decodeConfigWithEnv(data, format, true)
}

func formatFromPath(path string) (string, error) {
	if dot := strings.LastIndex(path, "."); dot >= 0 {
		format := strings.ToLower(path[dot+1:])
		switch format {
		case configFormatTOML, configFormatJSON, configFormatYAML, configFormatYML:
			return format, nil
		default:
			return "", fmt.Errorf("unsupported config file extension %q", format)
		}
	}
	return configFormatTOML, nil
}

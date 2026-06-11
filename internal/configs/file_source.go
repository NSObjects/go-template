package configs

import (
	"bytes"
	"context"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// FileSource 本地文件来源，支持 json/yaml/toml（通过文件后缀识别）。
type FileSource struct{ Path string }

func (f FileSource) Load(ctx context.Context) (Config, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return Config{}, err
	}
	return decodeConfigWithEnv(data, formatFromPath(f.Path), true)
}

// Watch 支持本地文件热更新，变更后回调新的 Config
func (f FileSource) Watch(ctx context.Context, onChange func(Config)) error {
	if f.Path == "" {
		return nil
	}
	viper.SetConfigFile(f.Path)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		data, err := os.ReadFile(e.Name)
		if err != nil {
			return
		}
		cfg, err := decodeConfig(data, formatFromPath(e.Name))
		if err == nil {
			onChange(cfg)
		}
	})
	return nil
}

// viperInit 支持根据文件后缀自动设置配置类型
func viperInit(configPath string) (err error) {
	viper.SetOptions(viper.KeyDelimiter("::"))
	if configPath != "" {
		viper.SetConfigType(configType(formatFromPath(configPath)))
		content, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		return viper.ReadConfig(bytes.NewBuffer(content))
	}
	return nil
}

func formatFromPath(path string) string {
	if dot := strings.LastIndex(path, "."); dot >= 0 {
		return strings.ToLower(path[dot+1:])
	}
	return "toml"
}

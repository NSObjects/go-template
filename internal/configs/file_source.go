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
	if err := viperInit(f.Path); err != nil {
		return Config{}, err
	}
	viper.SetEnvPrefix("ECHOADMIN")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// Watch 支持本地文件热更新，变更后回调新的 Config
func (f FileSource) Watch(ctx context.Context, onChange func(Config)) error {
	if f.Path == "" {
		return nil
	}
	viper.SetConfigFile(f.Path)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		var c Config
		if err := viper.Unmarshal(&c); err == nil {
			onChange(c)
		}
	})
	return nil
}

// viperInit 支持根据文件后缀自动设置配置类型
func viperInit(configPath string) (err error) {
	if configPath != "" {
		// 根据文件后缀自动设置类型
		ext := ""
		if dot := strings.LastIndex(configPath, "."); dot >= 0 {
			ext = strings.ToLower(configPath[dot+1:])
		}
		switch ext {
		case "json":
			viper.SetConfigType("json")
		case "yaml", "yml":
			viper.SetConfigType("yaml")
		default:
			viper.SetConfigType("toml")
		}
		content, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		return viper.ReadConfig(bytes.NewBuffer(content))
	}
	return nil
}

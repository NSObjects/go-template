package configs

import (
	"bytes"
	"context"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// ConsulSource 从 Consul KV 读取配置，支持 json/yaml/toml 三种格式。
type ConsulSource struct {
	Address string
	Token   string
	Key     string
	Format  string // json|yaml|toml
}

func (c ConsulSource) Load(ctx context.Context) (Config, error) {
	cli, err := api.NewClient(&api.Config{Address: c.Address, Token: c.Token})
	if err != nil {
		return Config{}, err
	}
	pair, _, err := cli.KV().Get(c.Key, nil)
	if err != nil || pair == nil {
		return Config{}, err
	}
	switch c.Format {
	case "json":
		viper.SetConfigType("json")
	case "yaml", "yml":
		viper.SetConfigType("yaml")
	default:
		viper.SetConfigType("toml")
	}
	if err := viper.ReadConfig(bytes.NewBuffer(pair.Value)); err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Watch 通过阻塞查询实现简单热更新
func (c ConsulSource) Watch(ctx context.Context, onChange func(Config)) error {
	cli, err := api.NewClient(&api.Config{Address: c.Address, Token: c.Token})
	if err != nil {
		return err
	}
	go func() {
		defer func() {}()
		var index uint64
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			q := &api.QueryOptions{WaitIndex: index, WaitTime: 5 * time.Minute}
			pair, meta, err := cli.KV().Get(c.Key, q)
			if err != nil || pair == nil {
				continue
			}
			if meta.LastIndex == index {
				continue
			}
			index = meta.LastIndex
			switch c.Format {
			case "json":
				viper.SetConfigType("json")
			case "yaml", "yml":
				viper.SetConfigType("yaml")
			default:
				viper.SetConfigType("toml")
			}
			if err := viper.ReadConfig(bytes.NewBuffer(pair.Value)); err == nil {
				var cfg Config
				if err := viper.Unmarshal(&cfg); err == nil {
					onChange(cfg)
				}
			}
		}
	}()
	return nil
}

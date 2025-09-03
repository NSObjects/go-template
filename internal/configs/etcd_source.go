package configs

import (
	"bytes"
	"context"
	"time"

	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdSource 从 etcd 读取配置，支持 json/yaml/toml 三种格式。
// 使用示例：
//
//	cfg := configs.NewCfgFrom(configs.EtcdSource{Endpoints: []string{"127.0.0.1:2379"}, Key: "/echo-admin/config", Format: "toml"})
type EtcdSource struct {
	Endpoints          []string
	Key                string
	Format             string // json|yaml|toml（默认 toml）
	Username           string
	Password           string
	DialTimeoutSeconds int
}

func (e EtcdSource) Load(ctx context.Context) (Config, error) {
	dialTimeout := 5 * time.Second
	if e.DialTimeoutSeconds > 0 {
		dialTimeout = time.Duration(e.DialTimeoutSeconds) * time.Second
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		Username:    e.Username,
		Password:    e.Password,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return Config{}, err
	}
	defer func() { _ = cli.Close() }()

	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := cli.Get(cctx, e.Key)
	if err != nil {
		return Config{}, err
	}
	if len(resp.Kvs) == 0 {
		return Config{}, nil
	}
	data := resp.Kvs[0].Value

	// 选择解析格式
	switch e.Format {
	case "json":
		viper.SetConfigType("json")
	case "yaml", "yml":
		viper.SetConfigType("yaml")
	default:
		viper.SetConfigType("toml")
	}

	if err := viper.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return Config{}, err
	}
	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// Watch 支持 etcd 热更新：watch 指定 Key，变更后回调新的 Config
func (e EtcdSource) Watch(ctx context.Context, onChange func(Config)) error {
	dialTimeout := 5 * time.Second
	if e.DialTimeoutSeconds > 0 {
		dialTimeout = time.Duration(e.DialTimeoutSeconds) * time.Second
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		Username:    e.Username,
		Password:    e.Password,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return err
	}

	go func() {
		defer func() { _ = cli.Close() }()
		rch := cli.Watch(ctx, e.Key)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Kv == nil {
					continue
				}
				switch e.Format {
				case "json":
					viper.SetConfigType("json")
				case "yaml", "yml":
					viper.SetConfigType("yaml")
				default:
					viper.SetConfigType("toml")
				}
				if err := viper.ReadConfig(bytes.NewBuffer(ev.Kv.Value)); err == nil {
					var c Config
					if err := viper.Unmarshal(&c); err == nil {
						onChange(c)
					}
				}
			}
		}
	}()
	return nil
}

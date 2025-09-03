/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Environment string

const (
	ENVIRONMENT Environment = "RUN_ENVIRONMENT"
)

type EnvironmentType int

const (
	Dev EnvironmentType = iota
	Docker
	Test
)

var runContext = map[string]EnvironmentType{
	"":       Dev,
	"docker": Docker,
	"test":   Test,
}

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota + 1
	// ProsecutionLevel InfoLevel is the default logging priority.
	ProsecutionLevel
)

var (
	Mysql  MysqlConfig
	System SystemConfig
	Log    LogConfig
	Mgo    Mongodb
	JWT    JWTConfig
)

type Config struct {
	Mysql   MysqlConfig        `mapstructure:"mysql"`
	System  SystemConfig       `mapstructure:"system"`
	Log     LogConfig          `mapstructure:"log"`
	Mongodb Mongodb            `mapstructure:"mongodb"`
	Redis   RedisConfig        `mapstructure:"redis"`
	JWT     JWTConfig          `mapstructure:"jwt"`
	CORS    CORSConfig         `mapstructure:"cors"`
	Casbin  CasbinConfig       `mapstructure:"casbin"`
	Kafka   KafkaConfig        `mapstructure:"kafka"`
	Etcd    EtcdClientConfig   `mapstructure:"etcd"`
	Consul  ConsulClientConfig `mapstructure:"consul"`
}

type SystemConfig struct {
	Port  string `mapstructure:"port"`
	Level Level  `mapstructure:"level"`
	Env   string `mapstructure:"env"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`

	Console ConsoleSinkConfig `mapstructure:"console"`
	File    FileSinkConfig    `mapstructure:"file"`

	Elasticsearch ElasticsearchSinkConfig `mapstructure:"elasticsearch"`
	Loki          LokiSinkConfig          `mapstructure:"loki"`
}

type ConsoleSinkConfig struct {
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type FileSinkConfig struct {
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	Format     string `mapstructure:"format"`
}

type ElasticsearchSinkConfig struct {
	URL     string        `mapstructure:"url"`
	Index   string        `mapstructure:"index"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type LokiSinkConfig struct {
	URL     string            `mapstructure:"url"`
	Labels  map[string]string `mapstructure:"labels"`
	Timeout time.Duration     `mapstructure:"timeout"`
}

type MysqlConfig struct {
	DockerHost   string `mapstructure:"docker_host"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	Database     string `mapstructure:"database"`
}

type JWTConfig struct {
	Secret    string   `mapstructure:"secret"`
	Expire    int      `mapstructure:"expire"`
	SkipPaths []string `mapstructure:"skip_paths"`
}

type Mongodb struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DataBase string `mapstructure:"database"`
}

type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type CasbinConfig struct {
	Model     string `mapstructure:"model"`
	ModelFile string `mapstructure:"model_file"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	ClientID string   `mapstructure:"client_id"`
	Topic    string   `mapstructure:"topic"`
}

type EtcdClientConfig struct {
	Endpoints          []string `mapstructure:"endpoints"`
	Key                string   `mapstructure:"key"`
	Format             string   `mapstructure:"format"`
	Username           string   `mapstructure:"username"`
	Password           string   `mapstructure:"password"`
	DialTimeoutSeconds int      `mapstructure:"dial_timeout_seconds"`
}

type ConsulClientConfig struct {
	Address string `mapstructure:"address"`
	Token   string `mapstructure:"token"`
	Key     string `mapstructure:"key"`
	Format  string `mapstructure:"format"`
}

func InitConfig(configPath string) (err error) {

	if err = viperInit(configPath); err != nil {
		return
	}

	var c Config
	if err = viper.Unmarshal(&c); err != nil {
		fmt.Println(err)
		return
	}

	Mysql = c.Mysql
	System = c.System
	Log = c.Log
	Mgo = c.Mongodb
	JWT = c.JWT
	return
}

func NewCfg(p string) Config {
	return NewCfgFrom(FileSource{Path: p})
}

// Source 抽象：配置来源（本地文件、远程配置中心等）。
// 实现方直接返回完整的 Config，便于从任意介质填充（文件/etcd/http 等）。
type Source interface {
	Load(ctx context.Context) (Config, error)
}

// WatchableSource 可选：支持热更新的配置源
// Watch 应启动监听并在变更时回调返回新的 Config
type WatchableSource interface {
	Source
	Watch(ctx context.Context, onChange func(Config)) error
}

// NewCfgFrom 通过自定义 Source 加载 Config。
func NewCfgFrom(src Source) Config {
	c, err := src.Load(context.Background())
	if err != nil {
		panic(err)
	}
	return c
}

func RunEnvironment() EnvironmentType {
	return runContext[os.Getenv(string(ENVIRONMENT))]
}

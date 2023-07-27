/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package configs

import (
	"bytes"
	"fmt"
	"os"

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
)

type Config struct {
	Mysql   MysqlConfig  `mapstructure:"mysql"`
	System  SystemConfig `mapstructure:"system"`
	Log     LogConfig    `mapstructure:"log"`
	Mongodb Mongodb      `mapstructure:"mongodb"`
	Redis   RedisConfig  `mapstructure:"redis"`
}

type SystemConfig struct {
	Port  string `mapstructure:"port"`
	Level Level  `mapstructure:"level"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type LogConfig struct {
	Path       string `mapstructure:"path"`
	Level      int    `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
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

type Mongodb struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DataBase string `mapstructure:"database"`
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
	return
}

func NewCfg(p string) Config {
	if err := viperInit(p); err != nil {
		panic(err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	return c
}

func viperInit(configPath string) (err error) {
	viper.SetConfigType("toml")
	if configPath != "" {
		content, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		err = viper.ReadConfig(bytes.NewBuffer(content))
		if err != nil {
			return err
		}
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("../../")
		viper.AddConfigPath("../")
		//viper.AddConfigPath("/etc/xxx/config")
		if err = viper.ReadInConfig(); err != nil {
			return
		}
	}
	return
}

func RunEnvironment() EnvironmentType {
	return runContext[os.Getenv(string(ENVIRONMENT))]
}

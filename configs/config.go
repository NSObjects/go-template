/*
 *
 * config.go
 * configs
 *
 * Created by lin on 2018/12/10 4:31 PM
 * Copyright © 2017-2018 PYL. All rights reserved.
 *
 */

package configs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

type Environment string

const (
	RUNENVIRONMENT Environment = "RUN_ENVIRONMENT"
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
	// InfoLevel is the default logging priority.
	PrdocutionLevel
)

var (
	Mysql  MysqlConfig
	System SystemConfig
	Log    LogConfig
)

type Config struct {
	Mysql  MysqlConfig  `mapstructure:"mysql"`
	System SystemConfig `mapstructure:"system"`
	Log    LogConfig    `mapstructure:"log"`
}

type SystemConfig struct {
	Prot  string `mapstructure:"prot"`
	Level Level  `mapstructure:"level"`
}

type LogConfig struct {
	Path       string `mapstructure:"path"`
	Level      int    `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MysqlConfig struct {
	DockerHost   string   `mapstructure:"docker_host"`
	Host         string   `mapstructure:"host"`
	Port         string   `mapstructure:"prot"`
	User         string   `mapstructure:"user"`
	Password     string   `mapstructure:"password"`
	MaxOpenConns int      `mapstructure:"max_open_conns"`
	MaxIdleConns int      `mapstructure:"max_idle_conns"`
	Database     []string `mapstructure:"database"`
}

func InitConfig(configPath string, configType string) (err error) {

	if err = viperInit(configPath, configType); err != nil {
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

	return
}

func viperInit(configPath string, configType string) (err error) {
	viper.SetConfigType(configType)
	if configPath != "" {
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			return err
		}
		err = viper.ReadConfig(bytes.NewBuffer(content))
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.AddConfigPath("../")
		//viper.AddConfigPath("/etc/xxx/config")
		if err = viper.ReadInConfig(); err != nil {
			return
		}
	}
	return
}

func RunEnvironment() EnvironmentType {
	return runContext[os.Getenv(string(RUNENVIRONMENT))]
}

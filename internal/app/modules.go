// Package app declares the modules included in this application.
package app

import (
	"github.com/NSObjects/go-template/internal/capabilities/kafka"
	"github.com/NSObjects/go-template/internal/capabilities/mongodb"
	"github.com/NSObjects/go-template/internal/capabilities/mysql"
	"github.com/NSObjects/go-template/internal/capabilities/redis"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// Modules returns the capability and business modules explicitly included in the app.
func Modules(cfg configs.Config) []module.Module {
	return []module.Module{
		mysql.New(cfg.Mysql),
		redis.New(cfg.Redis),
		mongodb.New(cfg.Mongodb),
		kafka.New(cfg.Kafka),
	}
}

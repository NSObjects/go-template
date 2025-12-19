/*
 * Created by lintao on 2023/7/27 上午10:04
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package cmd

import (
	"context"
	"log/slog"

	kitconfig "github.com/NSObjects/go-kit/config"
	"github.com/NSObjects/go-kit/log"
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/api/service"
	configs "github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/pkg/casbin"
	"github.com/NSObjects/go-template/internal/server"

	"go.uber.org/fx"
)

func Run(cfg string) {
	fx.New(
		fx.Module("config", fx.Provide(func() (configs.AppConfig, *kitconfig.Store[configs.AppConfig]) {
			merged, store := kitconfig.Bootstrap[configs.AppConfig](cfg)
			return merged, store
		}), fx.Provide(func(cfg configs.AppConfig) kitconfig.Config {
			return cfg.Config
		})),
		fx.Module("log", fx.Provide(func(cfg configs.AppConfig) log.Logger {
			// Create a simple console logger
			level := slog.LevelInfo
			if cfg.System.Level == "debug" {
				level = slog.LevelDebug
			}
			sink := log.NewConsoleSink(log.ConsoleSinkConfig{
				Format: "color",
				Output: "stdout",
			})
			logger := log.NewDefaultLogger(sink, level)
			log.SetGlobalLogger(logger)
			return logger
		})),
		fx.Module("db", db.Model),
		fx.Module("casbin", casbin.CasbinModule),
		fx.Module("biz", biz.Model),
		fx.Module("data", data.Model),
		fx.Module("service", service.Model),
		fx.Module("server", fx.Provide(server.NewEchoServer)),
		fx.Invoke(func(lifecycle fx.Lifecycle, s *server.EchoServer, cfg configs.AppConfig, logger log.Logger) {
			// 测试日志输出
			logger.Info("Application starting", slog.String("port", cfg.System.Port))

			lifecycle.Append(
				fx.Hook{
					OnStart: func(context.Context) error {
						logger.Info("Server starting", slog.String("port", cfg.System.Port))
						go s.Run(cfg.System.Port)
						return nil
					},
					OnStop: func(context.Context) error {
						logger.Info("Server stopping")
						return nil
					},
				})
		}),
	).Run()
}

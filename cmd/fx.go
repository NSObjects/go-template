/*
 * Created by lintao on 2023/7/27 上午10:04
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package cmd

import (
	"context"

	"log/slog"

	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/api/service"
	"github.com/NSObjects/go-template/internal/application"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/log"
	"github.com/NSObjects/go-template/internal/server"
	"github.com/NSObjects/go-template/internal/utils"

	"go.uber.org/fx"
)

func Run(cfg string) {
	fx.New(
		fx.Module("config", fx.Provide(func() (configs.Config, *configs.Store) {
			merged, store := configs.Bootstrap(cfg)
			return merged, store
		})),
		fx.Module("log", fx.Provide(func(cfg configs.Config) log.Logger {
			return log.NewLogger(cfg)
		})),
		fx.Module("data", db.Model, utils.CasbinModule),
		fx.Module("application", application.Module),
		fx.Module("repos", data.Model),
		fx.Module("service", service.Model),
		fx.Module("server", fx.Provide(server.NewEchoServer)),
		fx.Invoke(func(lifecycle fx.Lifecycle, s *server.EchoServer, cfg configs.Config, logger log.Logger) {
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

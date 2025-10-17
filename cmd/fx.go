/*
 * Created by lintao on 2023/7/27 上午10:04
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package cmd

import (
	"context"
	"log/slog"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/infra"
	"github.com/NSObjects/go-template/internal/infra/persistence"
	"github.com/NSObjects/go-template/internal/pkg/utils"
	"github.com/NSObjects/go-template/internal/shared/infra/log"
	"github.com/NSObjects/go-template/internal/shared/infra/server"
	"github.com/NSObjects/go-template/internal/user"

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
		fx.Module("data", persistence.Model, utils.CasbinModule),

		// DDD架构模块
		fx.Module("infra", infra.InfraModule),
		fx.Module("user", user.UserModule),
		fx.Module("server",
			fx.Provide(server.NewEchoServer),
			fx.Provide(server.NewRateLimiter),
		),
		fx.Invoke(func(lifecycle fx.Lifecycle, s *server.EchoServer, cfg configs.Config, logger log.Logger) {
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

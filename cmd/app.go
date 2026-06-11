/*
 * Created by lintao on 2023/7/27 上午10:04
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/api/service"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/log"
	"github.com/NSObjects/go-template/internal/server"
	"github.com/samber/do/v2"
)

func Run(cfg string) {
	i := do.New(
		registerConfig(cfg),
		registerLogger,
		db.Register,
		data.Register,
		biz.Register,
		service.Register,
		registerServer,
	)

	logger := do.MustInvoke[log.Logger](i)
	cfgValue := do.MustInvoke[configs.Config](i)
	dataManager := do.MustInvoke[*db.DataManager](i)
	echoServer := do.MustInvoke[*server.EchoServer](i)

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		report := i.ShutdownWithContext(shutdownCtx)
		if !report.Succeed {
			logger.Error("Dependency shutdown failed", slog.Any("errors", report.Errors))
		}
	}()

	logger.Info("Application starting", slog.String("port", cfgValue.System.Port))
	if err := dataManager.Start(context.Background()); err != nil {
		logger.Fatal("Data layer startup failed", slog.Any("error", err))
		return
	}

	logger.Info("Server starting", slog.String("port", cfgValue.System.Port))
	echoServer.Run(cfgValue.System.Port)
}

func registerConfig(path string) func(do.Injector) {
	return func(i do.Injector) {
		merged, store := configs.Bootstrap(path)
		do.ProvideValue(i, merged)
		do.ProvideValue(i, store)
	}
}

func registerLogger(i do.Injector) {
	do.Provide[log.Logger](i, func(i do.Injector) (log.Logger, error) {
		cfg, err := do.Invoke[configs.Config](i)
		if err != nil {
			return nil, err
		}
		return log.NewLogger(cfg), nil
	})
}

func registerServer(i do.Injector) {
	do.Provide(i, func(i do.Injector) (*server.EchoServer, error) {
		routes, err := do.Invoke[[]service.RegisterRouter](i)
		if err != nil {
			return nil, err
		}
		cfg, err := do.Invoke[configs.Config](i)
		if err != nil {
			return nil, err
		}
		store, err := do.Invoke[*configs.Store](i)
		if err != nil {
			return nil, err
		}
		return server.NewEchoServer(routes, nil, cfg, store), nil
	})
}

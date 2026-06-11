/*
 * Created by lintao on 2023/7/27 上午10:04
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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

type configBundle struct {
	cfg   configs.Config
	store *configs.Store
}

func Run(cfg string) {
	if err := run(cfg); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "application startup failed:", err)
		os.Exit(1)
	}
}

func run(cfg string) error {
	i := do.New(
		registerConfig(cfg),
		registerLogger,
		db.Register,
		data.Register,
		biz.Register,
		service.Register,
		registerServer,
	)

	cfgValue, err := do.Invoke[configs.Config](i)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	logger, err := do.Invoke[log.Logger](i)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		report := i.ShutdownWithContext(shutdownCtx)
		if !report.Succeed {
			logger.Error("Dependency shutdown failed", slog.Any("errors", report.Errors))
		}
	}()

	dataManager, err := do.Invoke[*db.DataManager](i)
	if err != nil {
		return fmt.Errorf("initialize data manager: %w", err)
	}
	echoServer, err := do.Invoke[*server.EchoServer](i)
	if err != nil {
		return fmt.Errorf("initialize echo server: %w", err)
	}

	logger.Info("Application starting", slog.String("port", cfgValue.System.Port))
	if err := dataManager.Start(context.Background()); err != nil {
		logger.Fatal("Data layer startup failed", slog.Any("error", err))
		return fmt.Errorf("start data layer: %w", err)
	}

	logger.Info("Server starting", slog.String("port", cfgValue.System.Port))
	echoServer.Run(cfgValue.System.Port)
	return nil
}

func registerConfig(path string) func(do.Injector) {
	return func(i do.Injector) {
		do.Provide(i, func(i do.Injector) (*configBundle, error) {
			merged, store, err := configs.BootstrapE(path)
			if err != nil {
				return nil, err
			}
			return &configBundle{
				cfg:   merged,
				store: store,
			}, nil
		})
		do.Provide(i, func(i do.Injector) (configs.Config, error) {
			bundle, err := do.Invoke[*configBundle](i)
			if err != nil {
				return configs.Config{}, err
			}
			return bundle.cfg, nil
		})
		do.Provide(i, func(i do.Injector) (*configs.Store, error) {
			bundle, err := do.Invoke[*configBundle](i)
			if err != nil {
				return nil, err
			}
			return bundle.store, nil
		})
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
		return server.NewEchoServer(routes, cfg, store), nil
	})
}

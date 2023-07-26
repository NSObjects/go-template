/*
 * Created by lintao on 2023/7/18 下午4:27
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */
package main

import (
	"os"

	"github.com/NSObjects/go-template/internal/api/server"

	"github.com/NSObjects/go-template/internal/api/data/db"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: newFlag(),
		Action: func(cCtx *cli.Context) error {
			log.Init()
			if err := configs.InitConfig(cCtx.String("conf")); err != nil {
				return err
			}

			api, err := NewHttpServer(configs.Config{}, db.NewDataSource(configs.Config{}))
			if err != nil {
				panic(err)
			}
			api.Run(configs.System.Port)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "conf",
			Aliases:     []string{"f"},
			DefaultText: "configs",
			Value:       "configs/config.toml",
			Usage:       "配置文件(.json,.yaml,.toml)",
			Required:    false,
		},
	}
}

func NewHttpServer(cfg configs.Config, dataSource *db.DataSource) (*server.EchoServer, error) {
	return server.NewEchoServer(dataSource, cfg), nil
}

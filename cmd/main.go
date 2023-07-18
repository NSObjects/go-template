/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */
package main

import (
	"github.com/NSObjects/go-template/internal/api/service"
	"github.com/NSObjects/go-template/tools/configs"
	"github.com/NSObjects/go-template/tools/db"
	"github.com/NSObjects/go-template/tools/log"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-template"
	app.Commands = []cli.Command{
		newWebCmd(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}

}

func newWebCmd() cli.Command {
	return cli.Command{
		Name:  "web",
		Usage: "run api server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Usage:    "配置文件(.json,.yaml,.toml)",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			if err := configs.InitConfig(c.String("conf")); err != nil {
				return err
			}
			log.Init()
			api, err := InitializeEchoServer()
			if err != nil {
				panic(err)
			}
			api.Run(configs.System.Prot)
			return nil
		},
	}
}

func InitializeEchoServer() (*service.EchoServer, error) {
	dataSource, err := db.NewDataSource()
	if err != nil {
		return nil, err
	}
	echoServer := service.NewEchoServer(dataSource)
	return echoServer, nil
}

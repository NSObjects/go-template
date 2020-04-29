/*
 *
 * main.go
 * pyl_erp
 *
 * Created by lin on 2018/12/10 3:05 PM
 * Copyright © 2017-2018 PYL. All rights reserved.
 *
 */
package main

import (
	"go-template/delivery"
	"go-template/tools/configs"
	"go-template/tools/db"
	"go-template/tools/log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-template"
	app.Commands = []*cli.Command{
		newWebCmd(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}

}

func newWebCmd() *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "run api server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
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

func InitializeEchoServer() (*delivery.EchoServer, error) {
	dataSource, err := db.NewDataSource()
	if err != nil {
		return nil, err
	}
	echoServer := delivery.NewEchoServer(dataSource)
	return echoServer, nil
}

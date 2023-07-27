/*
 * Created by lintao on 2023/7/18 下午4:27
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */
package main

import (
	"os"

	"github.com/NSObjects/go-template/internal/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: newFlag(),
		Action: func(cCtx *cli.Context) error {
			Run(cCtx.String("conf"))
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

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
	"flag"
	"fmt"
	"go-template/delivery/server"
	"go-template/tools/configs"
	"go-template/tools/db"
	"go-template/tools/log"
)

func main() {

	//database := db.NewDataSource()
	//err := database.Engine.Sync2(new(models.User))
	//if err != nil {
	//	panic(err)
	//}
	database, err := db.NewDataSource()
	if err != nil {
		panic(err)
	}
	api := server.NewEchoServer(database)

	api.LoadMiddleware()
	api.RegisterRouter()

	api.Run(configs.System.Prot)
}

func init() {
	initConfig()
	log.Init()
}

func initConfig() {

	configPath := flag.String("config", "", "config path")
	if flag.Parsed() == false {
		flag.Parse()
	}
	if err := configs.InitConfig(*configPath, "toml"); err != nil {
		fmt.Println(err)
	}
}
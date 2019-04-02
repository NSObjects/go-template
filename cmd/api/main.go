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
	"go-template/apis"
	"go-template/configs"
	_ "go-template/init"
	"go-template/tools/db"
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
	api := apis.NewEchoServer(database)

	api.LoadMiddleware()
	api.RegisterRouter()

	api.Run(configs.System.Prot)
}

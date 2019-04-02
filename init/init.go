/*
 *
 * init.go
 * init
 *
 * Created by lin on 2018/12/10 3:53 PM
 * Copyright © 2017-2018 PYL. All rights reserved.
 *
 */

package init

import (
	"flag"
	"fmt"
	"go-template/configs"
	"go-template/tools/log"
)

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

/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/NSObjects/go-template/internal/configs"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"xorm.io/xorm"
)

type DataBase int

const (
	DefaultDB DataBase = iota
)

// DataSource 在使用多个db的项目中在DataSource结构体中增加Engine即可
type DataSource struct {
	Engine *xorm.EngineGroup
}

var db *DataSource

func NewDataSource() (datasource *DataSource, err error) {
	if db != nil {
		return db, nil
	}

	db = new(DataSource)
	if db.Engine, err = createEngin(DefaultDB); err != nil {
		return nil, err
	}

	return db, nil
}

func NewTestDataSource(db2 *sql.DB) (datasource *DataSource, err error) {
	db = new(DataSource)
	if db.Engine, err = createEngin(DefaultDB, db2); err != nil {
		return nil, err
	}

	return db, nil
}

func createEngin(db DataBase, db2 ...*sql.DB) (engine *xorm.EngineGroup, err error) {
	conn, err := dataSource(db)
	if err != nil {
		return nil, err
	}
	if engine, err = xorm.NewEngineGroup("mysql", conn); err != nil {
		return nil, err
	}

	engine.SetMaxIdleConns(configs.Mysql.MaxIdleConns)
	engine.SetMaxOpenConns(configs.Mysql.MaxOpenConns)

	if configs.System.Level == 1 {
		engine.ShowSQL(true)
	}

	if len(db2) > 0 {
		engine.DB().DB = db2[0]
	}

	//retryConnect := 0
	//for {
	//	if err = engine.Ping(); err == nil {
	//		break
	//	}
	//	time.Sleep(3 * time.Second)
	//	if retryConnect == 5 {
	//		panic(err)
	//	}
	//	retryConnect++
	//}

	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}
	time.Local = cst
	engine.DatabaseTZ = cst

	return
}

func dataSource(db DataBase) (conn []string, err error) {

	host := ""

	switch configs.RunEnvironment() {
	case configs.Docker:
		host = configs.Mysql.DockerHost
	case configs.Dev:
		host = configs.Mysql.Host
	}

	if host == "" {
		return nil, errors.New("mysql host error, check config file")
	}

	if len(configs.Mysql.Database) <= 0 {
		return nil, errors.New("database not set, check config file")
	}

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4",
		configs.Mysql.User, configs.Mysql.Password, host, configs.Mysql.Port, configs.Mysql.Database[db])
	conn = append(conn, dbConn)
	return
}

/*
 * Created by lintao on 2023/7/18 下午4:00
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package service

import (
	"context"
	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"

	"github.com/NSObjects/go-template/internal/api/service/middlewares"
	"github.com/NSObjects/go-template/tools/db"
	"github.com/NSObjects/go-template/tools/log"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	server     *echo.Echo
	dataSource *db.DataSource
}

func (s *EchoServer) Server() *echo.Echo {
	return s.server
}

func NewEchoServer(db *db.DataSource) *EchoServer {
	s := &EchoServer{
		server:     echo.New(),
		dataSource: db,
	}
	s.loadMiddleware()
	s.registerRouter()
	return s
}

func (s *EchoServer) loadMiddleware() {
	s.server.Validator = &middlewares.Validator{Validator: validator.New()}
	s.server.Use(middleware.Gzip())
	s.server.Use(middleware.Recover())
	s.server.Use(middleware.Logger())
	s.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//todo 域名设置
		//AllowOrigins:     []string{"http://xxx:8080","https://xxxx:8080"},
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}))
}

func InitializeController(d *db.DataSource) *UserController {
	userDataSource := data.NewUserDataSource(d)
	userHandler := biz.NewUserHandler(userDataSource)
	userController := NewUserController(userHandler)
	return userController
}

func (s *EchoServer) registerRouter() {
	routers := []RegisterRouter{
		InitializeController(s.dataSource),
	}

	g := s.server.Group("api")
	for _, v := range routers {
		v.RegisterRouter(g)
	}
}

func (s *EchoServer) Run(port string) {
	go func() {
		if err := s.server.Start(port); err != nil {
			log.Panic(err)
		}
		log.Info("start")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error(err)
	}
}

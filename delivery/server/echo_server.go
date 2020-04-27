/*
 *
 * echo_server.go
 * apis
 *
 * Created by lintao on 2019-01-31 11:26
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package server

import (
	"context"

	"go-template/delivery"
	"go-template/delivery/middlewares"
	"go-template/repository"
	"go-template/tools/db"
	"go-template/tools/log"
	"go-template/usecase"
	"os"
	"os/signal"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

type RegisterRouter interface {
	RegisterRouter(s *echo.Group, middlewareFunc ...echo.MiddlewareFunc)
}

type EchoServer struct {
	server     *echo.Echo
	dataSource *db.DataSource
}

func NewEchoServer(db *db.DataSource) *EchoServer {
	return &EchoServer{
		server:     echo.New(),
		dataSource: db,
	}
}

func (s *EchoServer) LoadMiddleware() {
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

func (s *EchoServer) RegisterRouter() {
	dataSource, err := db.NewDataSource()
	if err != nil {
		panic(err)
	}

	routers := []RegisterRouter{
		delivery.NewUserController(usecase.NewUserHandler(repository.NewUserDataSource(dataSource))),
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

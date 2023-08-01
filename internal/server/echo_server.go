/*
 * Created by lintao on 2023/7/26 下午2:22
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package server

import (
	"context"
	"github.com/NSObjects/go-template/internal/log"
	"github.com/NSObjects/go-template/internal/resp"
	"net/http"

	"github.com/NSObjects/go-template/internal/api/service"
	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"

	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	server  *echo.Echo
	Routers []service.RegisterRouter `group:"routes"`
}

func (s *EchoServer) Server() *echo.Echo {
	return s.server
}

func NewEchoServer(routes []service.RegisterRouter) *EchoServer {
	s := &EchoServer{
		server:  echo.New(),
		Routers: routes,
	}
	s.loadMiddleware()
	s.registerRouter()
	return s
}

func errorHandler(err error, c echo.Context) {
	er := resp.APIError(err, c)
	if er != nil {
		log.Error(er)
	}
}

func (s *EchoServer) loadMiddleware() {
	s.server.Validator = &middlewares.Validator{Validator: validator.New()}
	s.server.Use(middleware.Gzip())
	s.server.HTTPErrorHandler = errorHandler
	//s.server.Use(middleware.Recover())
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

func (s *EchoServer) registerRouter() {
	g := s.server.Group("api")
	for _, v := range s.Routers {
		v.RegisterRouter(g)
	}
}

func (s *EchoServer) Run(port string) {
	go func() {
		if err := s.server.Start(port); err != nil && err != http.ErrServerClosed {
			s.server.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.server.Logger.Fatal(err)
	}
}

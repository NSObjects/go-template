/*
 * Created by lintao on 2023/7/26 下午2:22
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NSObjects/go-template/internal/api/data"
	"github.com/NSObjects/go-template/internal/api/service"
	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/log"
	"github.com/NSObjects/go-template/internal/resp"
	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marmotedu/errors"
)

type EchoServer struct {
	server  *echo.Echo
	Routers []service.RegisterRouter `group:"routes"`
	cfg     configs.Config
}

func (s *EchoServer) Server() *echo.Echo {
	return s.server
}

func NewEchoServer(routes []service.RegisterRouter, cfg configs.Config) *EchoServer {
	s := &EchoServer{
		server:  echo.New(),
		Routers: routes,
		cfg:     cfg,
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
	s.server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	s.server.Validator = &middlewares.Validator{Validator: validator.New()}
	s.server.Use(middleware.Gzip())
	s.server.HTTPErrorHandler = errorHandler
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(data.JwtCustomClaims)
		},
		SigningKey: []byte(s.cfg.JWT.Secret),
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/login/account"
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return errors.WrapC(err, code.ErrSignatureInvalid, "JWT签名无效")
		},
	}

	s.server.Use(echojwt.WithConfig(config))
	//s.server.Use(middleware.Recover())

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
	g.GET("/api", func(c echo.Context) error {
		return resp.ListDataResponse(s.server.Routes(), int64(len(s.server.Routes())), c)
	})
}

func (s *EchoServer) Run(port string) {
	go func() {
		if err := s.server.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.server.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //nolint:gomnd
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.server.Logger.Fatal(err)
	}
}

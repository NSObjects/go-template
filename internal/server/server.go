package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const (
	apiPrefix  = "/api"
	appName    = "go-template"
	appVersion = "1.0.0"
)

// Server owns the Echo HTTP server lifecycle and system routes.
type Server struct {
	echo      *echo.Echo
	api       *echo.Group
	config    *Config
	appConfig configs.Config
}

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type infoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Time    string `json:"time"`
}

// Echo returns the underlying Echo instance for HTTP adapter tests and
// framework-level integration.
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

// API returns the root API route group used by boot to register business routes.
func (s *Server) API() *echo.Group {
	return s.api
}

// New creates an Echo-backed HTTP server.
func New(cfg configs.Config) (*Server, error) {
	cfg = configs.Normalize(cfg)
	if err := configs.Validate(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	e := echo.New()
	s := &Server{
		echo:      e,
		api:       e.Group(apiPrefix),
		config:    FromAppConfig(cfg),
		appConfig: cfg,
	}

	s.configureEcho()
	if err := s.installMiddleware(); err != nil {
		return nil, err
	}
	s.registerSystemRoutes()

	return s, nil
}

func (s *Server) configureEcho() {
	s.echo.Validator = &middlewares.Validator{Validator: validator.New()}
	s.echo.HTTPErrorHandler = middlewares.ErrorHandler
	s.echo.HideBanner = s.config.HideBanner
	s.echo.Debug = s.config.Debug
	s.echo.Server.ReadTimeout = s.config.ReadTimeout
	s.echo.Server.WriteTimeout = s.config.WriteTimeout
	s.echo.Server.IdleTimeout = s.config.IdleTimeout
}

func (s *Server) installMiddleware() error {
	if err := middlewares.ApplyMiddlewares(s.echo, s.middlewareConfig()); err != nil {
		return fmt.Errorf("install middleware: %w", err)
	}
	return nil
}

func (s *Server) middlewareConfig() *middlewares.MiddlewareConfig {
	jwtConfig := middlewares.CreateJWTConfig(
		s.appConfig.JWT.Secret,
		s.appConfig.JWT.SkipPaths,
		s.appConfig.JWT.Enabled,
	)

	return &middlewares.MiddlewareConfig{
		EnableRecovery:       true,
		EnableRequestContext: true,
		EnableLogger:         true,
		EnableGzip:           true,
		EnableCORS:           false,
		EnableJWT:            jwtConfig.Enabled,
		JWT:                  jwtConfig,
	}
}

func (s *Server) registerSystemRoutes() {
	s.api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, healthResponse{
			Status: "ok",
			Time:   time.Now().Format(time.RFC3339),
		})
	})

	s.api.GET("/info", func(c echo.Context) error {
		return c.JSON(http.StatusOK, infoResponse{
			Name:    appName,
			Version: appVersion,
			Time:    time.Now().Format(time.RFC3339),
		})
	})
}

// Run starts the HTTP server and blocks until ctx is cancelled or startup fails.
func (s *Server) Run(ctx context.Context) error {
	if ctx == nil {
		return errors.New("server run: nil context")
	}

	errCh := make(chan error, 1)
	go s.serve(errCh)

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	return s.shutdown(errCh)
}

func (s *Server) serve(errCh chan<- error) {
	addr := s.config.Port
	if addr == "" {
		addr = defaultServerPort
	}

	s.echo.Logger.Infof("starting server on %s", addr)
	err := s.echo.Start(addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		errCh <- err
		return
	}
	errCh <- nil
}

func (s *Server) shutdown(errCh <-chan error) error {
	s.echo.Logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.echo.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-shutdownCtx.Done():
		return fmt.Errorf("wait for server shutdown: %w", shutdownCtx.Err())
	}

	s.echo.Logger.Info("server exited")
	return nil
}

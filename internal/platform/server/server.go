package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	zlog "github.com/rs/zerolog/log"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/server/middlewares"
)

const apiPrefix = "/api"

// Server owns the Echo HTTP server lifecycle and system routes.
type Server struct {
	echo           *echo.Echo
	api            *echo.Group
	config         *Config
	appConfig      configs.Config
	statusReporter StatusReporter
}

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

// CapabilityStatus is the server-owned JSON shape for capability routes.
type CapabilityStatus struct {
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	Available bool   `json:"available"`
	State     string `json:"state"`
	Message   string `json:"message,omitempty"`
}

// StatusReporter supplies readiness and capability status to system routes.
type StatusReporter interface {
	Status(context.Context) []CapabilityStatus
	Ready(context.Context) error
}

type infoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Time    string `json:"time"`
}

type readinessResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type capabilitiesResponse struct {
	Capabilities []CapabilityStatus `json:"capabilities"`
	Time         string             `json:"time"`
}

// Option customizes the HTTP server.
type Option func(*Server)

// WithStatusReporter installs the readiness and capability reporter.
func WithStatusReporter(reporter StatusReporter) Option {
	return func(s *Server) {
		s.statusReporter = reporter
	}
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
func New(cfg configs.Config, opts ...Option) (*Server, error) {
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
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
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
	httpConfig := s.appConfig.HTTP

	return &middlewares.MiddlewareConfig{
		EnableRecovery:       !httpConfig.RecoveryDisabled,
		EnableRequestContext: !httpConfig.RequestContextDisabled,
		EnableLogger:         !httpConfig.RequestLogDisabled,
		EnableTracing:        s.appConfig.Tracing.Enabled,
		TracingServiceName:   s.appConfig.App.Name,
		EnableGzip:           !httpConfig.GzipDisabled,
		EnableCORS:           httpConfig.CORS.Enabled,
		CORS:                 corsMiddlewareConfig(httpConfig.CORS),
		EnableJWT:            jwtConfig.Enabled,
		JWT:                  jwtConfig,
	}
}

func corsMiddlewareConfig(cfg configs.CORSConfig) middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		ExposeHeaders:    cfg.ExposeHeaders,
		MaxAge:           cfg.MaxAgeSeconds,
	}
}

func (s *Server) registerSystemRoutes() {
	s.api.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, healthResponse{
			Status: "ok",
			Time:   time.Now().Format(time.RFC3339),
		})
	})

	s.api.GET("/info", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, infoResponse{
			Name:    s.appConfig.App.Name,
			Version: s.appConfig.App.Version,
			Time:    time.Now().Format(time.RFC3339),
		})
	})

	s.api.GET("/ready", func(c *echo.Context) error {
		response := readinessResponse{
			Status: "ready",
			Time:   time.Now().Format(time.RFC3339),
		}
		if s.statusReporter == nil {
			return c.JSON(http.StatusOK, response)
		}
		if err := s.statusReporter.Ready(c.Request().Context()); err != nil {
			response.Status = "unavailable"
			return c.JSON(http.StatusServiceUnavailable, response)
		}
		return c.JSON(http.StatusOK, response)
	})

	s.api.GET("/capabilities", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, capabilitiesResponse{
			Capabilities: s.statuses(c.Request().Context()),
			Time:         time.Now().Format(time.RFC3339),
		})
	})
}

func (s *Server) statuses(ctx context.Context) []CapabilityStatus {
	if s.statusReporter == nil {
		return nil
	}
	statuses := s.statusReporter.Status(ctx)
	copied := make([]CapabilityStatus, len(statuses))
	copy(copied, statuses)
	return copied
}

// Run starts the HTTP server and blocks until ctx is canceled or startup fails.
func (s *Server) Run(ctx context.Context) error {
	if ctx == nil {
		return errors.New("server run: nil context")
	}

	addr := s.config.Port
	if addr == "" {
		addr = defaultServerPort
	}

	zlog.Info().Str("addr", addr).Msg("starting server")
	shutdownErrCh := make(chan error, 1)
	startConfig := s.startConfig(addr)
	startConfig.OnShutdownError = func(err error) {
		select {
		case shutdownErrCh <- err:
		default:
		}
	}

	if err := startConfig.Start(ctx, s.echo); err != nil {
		return err
	}
	select {
	case err := <-shutdownErrCh:
		return fmt.Errorf("shutdown server: %w", err)
	default:
	}

	zlog.Info().Msg("server exited")
	return nil
}

func (s *Server) startConfig(addr string) echo.StartConfig {
	return echo.StartConfig{
		Address:         addr,
		HideBanner:      s.config.HideBanner,
		HidePort:        true,
		GracefulTimeout: s.config.ShutdownTimeout,
		BeforeServeFunc: func(server *http.Server) error {
			server.ReadTimeout = s.config.ReadTimeout
			server.WriteTimeout = s.config.WriteTimeout
			server.IdleTimeout = s.config.IdleTimeout
			return nil
		},
	}
}

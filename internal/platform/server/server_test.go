package server

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"

	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/configs"
)

func TestServerEcho(t *testing.T) {
	server := mustNewServer(t, configs.Config{})

	assert.NotNil(t, server)
	assert.NotNil(t, server.Echo())
	assert.IsType(t, &echo.Echo{}, server.Echo())
}

func TestServerConfigureEcho(t *testing.T) {
	e := echo.New()
	server := &Server{
		echo:   e,
		config: DefaultConfig(),
	}

	server.configureEcho()

	assert.NotNil(t, server.echo.Validator)
	assert.NotNil(t, server.echo.HTTPErrorHandler)
}

func TestServerStartConfigAppliesRuntimeSettings(t *testing.T) {
	server := &Server{
		config: &Config{
			Port:            ":9323",
			ReadTimeout:     2 * time.Second,
			WriteTimeout:    3 * time.Second,
			IdleTimeout:     4 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			HideBanner:      true,
		},
	}

	startConfig := server.startConfig(server.config.Port)
	httpServer := &http.Server{}
	if err := startConfig.BeforeServeFunc(httpServer); err != nil {
		t.Fatalf("BeforeServeFunc() error = %v", err)
	}

	assert.Equal(t, ":9323", startConfig.Address)
	assert.True(t, startConfig.HideBanner)
	assert.True(t, startConfig.HidePort)
	assert.Equal(t, 5*time.Second, startConfig.GracefulTimeout)
	assert.Equal(t, 2*time.Second, httpServer.ReadTimeout)
	assert.Equal(t, 3*time.Second, httpServer.WriteTimeout)
	assert.Equal(t, 4*time.Second, httpServer.IdleTimeout)
}

func TestServerMiddlewareConfig(t *testing.T) {
	server := &Server{
		config:    DefaultConfig(),
		appConfig: configs.Config{},
	}

	config := server.middlewareConfig()

	assert.NotNil(t, config)
	assert.True(t, config.EnableRecovery)
	assert.True(t, config.EnableRequestContext)
	assert.True(t, config.EnableLogger)
	assert.True(t, config.EnableGzip)
	assert.False(t, config.EnableCORS)
	assert.False(t, config.EnableJWT)
	assert.NotNil(t, config.JWT)
}

func TestServerMiddlewareConfigEnablesJWTFromConfig(t *testing.T) {
	server := &Server{
		config: DefaultConfig(),
		appConfig: configs.Config{
			JWT: configs.JWTConfig{
				Enabled:   true,
				Secret:    "test-secret",
				SkipPaths: []string{"/api/health"},
			},
		},
	}

	config := server.middlewareConfig()

	assert.True(t, config.EnableJWT)
	assert.NotNil(t, config.JWT)
	assert.True(t, config.JWT.Enabled)
	assert.Equal(t, []byte("test-secret"), config.JWT.SigningKey)
	assert.Equal(t, []string{"/api/health"}, config.JWT.SkipPaths)
}

func TestServerMiddlewareConfigUsesHTTPConfig(t *testing.T) {
	server := &Server{
		config: DefaultConfig(),
		appConfig: configs.Config{
			HTTP: configs.HTTPConfig{
				RecoveryDisabled:       true,
				RequestContextDisabled: true,
				RequestLogDisabled:     true,
				GzipDisabled:           true,
				CORS: configs.CORSConfig{
					Enabled:          true,
					AllowOrigins:     []string{"https://app.example.com"},
					AllowMethods:     []string{"GET"},
					AllowHeaders:     []string{"Authorization"},
					AllowCredentials: true,
					ExposeHeaders:    []string{"X-Request-ID"},
					MaxAgeSeconds:    600,
				},
			},
		},
	}

	config := server.middlewareConfig()

	assert.False(t, config.EnableRecovery)
	assert.False(t, config.EnableRequestContext)
	assert.False(t, config.EnableLogger)
	assert.False(t, config.EnableGzip)
	assert.True(t, config.EnableCORS)
	assert.Equal(t, []string{"https://app.example.com"}, config.CORS.AllowOrigins)
	assert.Equal(t, []string{"GET"}, config.CORS.AllowMethods)
	assert.Equal(t, []string{"Authorization"}, config.CORS.AllowHeaders)
	assert.True(t, config.CORS.AllowCredentials)
	assert.Equal(t, []string{"X-Request-ID"}, config.CORS.ExposeHeaders)
	assert.Equal(t, 600, config.CORS.MaxAge)
}

func TestServerRegisterSystemRoutes(t *testing.T) {
	e := echo.New()
	server := &Server{
		echo:   e,
		api:    e.Group(apiPrefix),
		config: DefaultConfig(),
	}

	server.registerSystemRoutes()

	routes := server.echo.Router().Routes()
	assert.NotEmpty(t, routes)

	hasHealthRoute := false
	hasInfoRoute := false
	hasRoutesRoute := false

	for _, route := range routes {
		if route.Path == "/api/health" && route.Method == http.MethodGet {
			hasHealthRoute = true
		}
		if route.Path == "/api/routes" && route.Method == http.MethodGet {
			hasRoutesRoute = true
		}
		if route.Path == "/api/info" && route.Method == http.MethodGet {
			hasInfoRoute = true
		}
	}

	assert.True(t, hasHealthRoute, "Health route should be registered")
	assert.True(t, hasInfoRoute, "Info route should be registered")
	assert.False(t, hasRoutesRoute, "Routes introspection route should not be registered by default")
}

func TestServerRunReturnsStartupError(t *testing.T) {
	server := &Server{
		echo: echo.New(),
		config: &Config{
			Port:            "invalid-address",
			ReadTimeout:     1 * time.Second,
			WriteTimeout:    1 * time.Second,
			IdleTimeout:     1 * time.Second,
			ShutdownTimeout: 1 * time.Second,
		},
	}

	err := server.Run(context.Background())
	assert.Error(t, err)
}

func TestServerRunReturnsShutdownError(t *testing.T) {
	addr := freeLocalAddr(t)
	started := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseHandler := func() {
		releaseOnce.Do(func() {
			close(release)
		})
	}
	defer releaseHandler()

	server := newSlowShutdownServer(addr, started, release)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runErrCh := make(chan error, 1)
	go func() {
		runErrCh <- server.Run(ctx)
	}()

	clientErrCh := make(chan error, 1)
	go func() {
		clientErrCh <- requestUntilServerReady("http://"+addr+"/slow", 2*time.Second)
	}()

	waitForSlowRequestStarted(t, started, runErrCh, clientErrCh)
	cancel()

	err := waitForRunResult(t, runErrCh)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Run() error = %v, want context deadline exceeded", err)
	}
	assert.Contains(t, err.Error(), "shutdown server")

	releaseHandler()
	waitForClientRequest(t, clientErrCh)
}

func newSlowShutdownServer(addr string, started chan<- struct{}, release <-chan struct{}) *Server {
	e := echo.New()
	e.GET("/slow", func(c *echo.Context) error {
		close(started)
		<-release
		return c.NoContent(http.StatusNoContent)
	})

	return &Server{
		echo: e,
		config: &Config{
			Port:            addr,
			ReadTimeout:     time.Second,
			WriteTimeout:    time.Second,
			IdleTimeout:     time.Second,
			ShutdownTimeout: 50 * time.Millisecond,
			HideBanner:      true,
		},
	}
}

func waitForSlowRequestStarted(
	t *testing.T,
	started <-chan struct{},
	runErrCh <-chan error,
	clientErrCh <-chan error,
) {
	t.Helper()

	select {
	case <-started:
	case err := <-runErrCh:
		t.Fatalf("Run() returned before request started: %v", err)
	case err := <-clientErrCh:
		t.Fatalf("client request failed before reaching handler: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for slow request to start")
	}
}

func waitForRunResult(t *testing.T, runErrCh <-chan error) error {
	t.Helper()

	select {
	case err := <-runErrCh:
		return err
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Run() shutdown result")
		return nil
	}
}

func waitForClientRequest(t *testing.T, clientErrCh <-chan error) {
	t.Helper()

	select {
	case err := <-clientErrCh:
		assert.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for slow request to finish")
	}
}

func TestServerRunRejectsNilContext(t *testing.T) {
	server := mustNewServer(t, configs.Config{})

	var ctx context.Context
	err := server.Run(ctx)
	if err == nil {
		t.Fatal("Run(nil) error = nil, want nil context error")
	}
	assert.Contains(t, err.Error(), "nil context")
}

func TestServerAPIGroupRegistersBusinessRoutes(t *testing.T) {
	server := mustNewServer(t, configs.Config{})

	server.API().GET("/ping", func(c *echo.Context) error {
		return c.NoContent(204)
	})

	hasPingRoute := false
	for _, route := range server.Echo().Router().Routes() {
		if route.Method == http.MethodGet && route.Path == "/api/ping" {
			hasPingRoute = true
			break
		}
	}
	assert.True(t, hasPingRoute, "API group should register routes under /api")
}

func TestServerNew(t *testing.T) {
	cfg := configs.Config{
		System: configs.SystemConfig{
			Port:  ":9323",
			Level: 1,
		},
		JWT: configs.JWTConfig{
			Enabled:   true,
			Secret:    "test-secret",
			SkipPaths: []string{"/api/health"},
		},
	}

	server, err := New(cfg)

	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.NotNil(t, server.echo)
	assert.NotNil(t, server.config)
	assert.NotNil(t, server.api)
	assert.Equal(t, configs.Normalize(cfg), server.appConfig)
}

func TestServerNewReturnsConfigError(t *testing.T) {
	server, err := New(configs.Config{
		JWT: configs.JWTConfig{
			Enabled: true,
		},
	})

	assert.Nil(t, server)
	assert.Error(t, err)
}

func TestServerSystemRoutes(t *testing.T) {
	e := echo.New()
	server := &Server{
		echo:   e,
		api:    e.Group(apiPrefix),
		config: DefaultConfig(),
	}

	server.registerSystemRoutes()

	routes := server.echo.Router().Routes()
	assert.NotEmpty(t, routes)

	hasSystemRoutes := false
	hasRoutesRoute := false
	for _, route := range routes {
		if route.Path == "/api/health" || route.Path == "/api/info" {
			hasSystemRoutes = true
		}
		if route.Path == "/api/routes" {
			hasRoutesRoute = true
		}
	}
	assert.True(t, hasSystemRoutes, "System routes should be registered")
	assert.False(t, hasRoutesRoute, "Routes introspection route should not be registered by default")
}

func TestServerInfoRouteUsesConfiguredAppIdentity(t *testing.T) {
	server := mustNewServer(t, configs.Config{
		App: configs.AppConfig{
			Name:    "payments-api",
			Version: "2026.06.17",
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
	server.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body infoResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode info response: %v", err)
	}

	assert.Equal(t, "payments-api", body.Name)
	assert.Equal(t, "2026.06.17", body.Version)
}

func TestServerUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	server := mustNewServer(t, configs.Config{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/health", nil)
	server.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	assert.Equal(t, float64(apperr.ErrMethodNotAllowed), body["code"])
	assert.Equal(t, "Method not allowed", body["message"])
}

func TestServerReadinessReturnsUnavailableCapability(t *testing.T) {
	server := mustNewServer(t, configs.Config{}, WithStatusReporter(fakeStatusReporter{
		statuses: []CapabilityStatus{
			{Name: "redis", Enabled: true, Available: false, State: "unavailable", Message: "dial refused"},
		},
		readyErr: errors.New("redis unavailable"),
	}))

	readyRec := httptest.NewRecorder()
	readyReq := httptest.NewRequest(http.MethodGet, "/api/ready", nil)
	server.Echo().ServeHTTP(readyRec, readyReq)

	assert.Equal(t, http.StatusServiceUnavailable, readyRec.Code)

	var readyBody map[string]any
	if err := json.Unmarshal(readyRec.Body.Bytes(), &readyBody); err != nil {
		t.Fatalf("decode ready response: %v", err)
	}
	assert.Equal(t, "unavailable", readyBody["status"])
	if _, ok := readyBody["capabilities"]; ok {
		t.Fatal("ready response exposed capabilities, want generic readiness only")
	}

	healthRec := httptest.NewRecorder()
	healthReq := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	server.Echo().ServeHTTP(healthRec, healthReq)

	assert.Equal(t, http.StatusOK, healthRec.Code)
}

func TestServerCapabilitiesReturnsStatusList(t *testing.T) {
	server := mustNewServer(t, configs.Config{}, WithStatusReporter(fakeStatusReporter{
		statuses: []CapabilityStatus{
			{Name: "mysql", Enabled: true, Available: true, State: "available", Message: "ping ok"},
			{Name: "mongodb", Enabled: false, Available: false, State: "disabled"},
		},
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/capabilities", nil)
	server.Echo().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body struct {
		Capabilities []CapabilityStatus `json:"capabilities"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode capabilities response: %v", err)
	}
	if len(body.Capabilities) != 2 {
		t.Fatalf("capabilities len = %d, want 2", len(body.Capabilities))
	}
	if body.Capabilities[0].Name != "mysql" || !body.Capabilities[0].Available {
		t.Fatalf("first capability = %+v, want available mysql", body.Capabilities[0])
	}
	if body.Capabilities[1].Name != "mongodb" || body.Capabilities[1].Available {
		t.Fatalf("second capability = %+v, want disabled mongodb", body.Capabilities[1])
	}
}

type fakeStatusReporter struct {
	statuses []CapabilityStatus
	readyErr error
}

func (f fakeStatusReporter) Status(context.Context) []CapabilityStatus {
	return f.statuses
}

func (f fakeStatusReporter) Ready(context.Context) error {
	return f.readyErr
}

func TestServerConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, ":9322", config.Port)
	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.IdleTimeout)
	assert.Equal(t, 10*time.Second, config.ShutdownTimeout)
	assert.True(t, config.HideBanner)

	config.Port = ":9090"
	config.HideBanner = false

	assert.Equal(t, ":9090", config.Port)
	assert.False(t, config.HideBanner)
}

func freeLocalAddr(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on free local port: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("close free local port listener: %v", err)
	}
	return addr
}

func requestUntilServerReady(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: timeout}
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			return resp.Body.Close()
		}
		lastErr = err
		time.Sleep(5 * time.Millisecond)
	}
	return lastErr
}

func mustNewServer(t *testing.T, cfg configs.Config, opts ...Option) *Server {
	t.Helper()

	server, err := New(cfg, opts...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return server
}

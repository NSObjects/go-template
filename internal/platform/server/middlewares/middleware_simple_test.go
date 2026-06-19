package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/requestctx"
)

func TestDefaultMiddlewareConfig(t *testing.T) {
	config := DefaultMiddlewareConfig()

	assert.NotNil(t, config)
	assert.True(t, config.EnableRecovery)
	assert.True(t, config.EnableRequestContext)
	assert.True(t, config.EnableLogger)
	assert.True(t, config.EnableGzip)
	assert.False(t, config.EnableCORS)
	assert.False(t, config.EnableJWT)
	assert.NotNil(t, config.JWT)
}

func TestDefaultJWTConfig(t *testing.T) {
	config := DefaultJWTConfig()

	assert.NotNil(t, config)
	assert.False(t, config.Enabled)
	assert.NotNil(t, config.SkipPaths)
	assert.Empty(t, config.SigningKey)
}

func TestCreateJWTConfig(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		skipPaths []string
		enabled   bool
	}{
		{
			name:      "enabled JWT",
			secret:    "test-secret",
			skipPaths: []string{"/api/health"},
			enabled:   true,
		},
		{
			name:      "disabled JWT",
			secret:    "test-secret",
			skipPaths: []string{},
			enabled:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CreateJWTConfig(tt.secret, tt.skipPaths, tt.enabled)
			assert.NotNil(t, config)
			assert.Equal(t, []byte(tt.secret), config.SigningKey)
			assert.Equal(t, tt.skipPaths, config.SkipPaths)
			assert.Equal(t, tt.enabled, config.Enabled)
		})
	}
}

func TestApplyMiddlewares(t *testing.T) {
	e := echo.New()
	config := DefaultMiddlewareConfig()

	assert.NoError(t, ApplyMiddlewares(e, config))
	assert.NotNil(t, e)
}

func TestRequestLoggerPreservesRenderedApplicationErrorStatus(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = ErrorHandler
	config := DefaultMiddlewareConfig()

	assert.NoError(t, ApplyMiddlewares(e, config))
	e.GET("/bad", func(_ *echo.Context) error {
		return apperr.NewBadRequest("invalid payload")
	})

	req := httptest.NewRequest(http.MethodGet, "/bad", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assertErrorPayload(t, rec, apperr.ErrBadRequest, "bad request: invalid payload")
}

func TestRequestContextMiddlewareStoresMetadata(t *testing.T) {
	e := echo.New()
	e.Use(RequestContext())
	e.GET("/ping", func(c *echo.Context) error {
		info, ok := requestctx.FromContext(c.Request().Context())
		if !ok {
			t.Fatal("request context metadata missing")
		}
		if info.RequestID != "req-123" {
			t.Fatalf("RequestID = %q, want req-123", info.RequestID)
		}
		if info.TraceID != "trace-123" {
			t.Fatalf("TraceID = %q, want trace-123", info.TraceID)
		}
		if info.UserID != "" {
			t.Fatalf("UserID = %q, want empty because request metadata does not authenticate users", info.UserID)
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(headerRequestID, "req-123")
	req.Header.Set(headerTraceID, "trace-123")
	req.Header.Set("X-User-ID", "user-123")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "req-123", rec.Header().Get(headerRequestID))
}

func TestRequestContextMiddlewareCleansUntrustedMetadata(t *testing.T) {
	overlongRequestID := strings.Repeat("a", 129)
	overlongSpanID := strings.Repeat("b", 129)

	e := echo.New()
	e.Use(RequestContext())
	e.GET("/ping", func(c *echo.Context) error {
		info, ok := requestctx.FromContext(c.Request().Context())
		if !ok {
			t.Fatal("request context metadata missing")
		}
		if info.RequestID == "" {
			t.Fatal("RequestID is empty, want generated request id")
		}
		if info.RequestID == overlongRequestID {
			t.Fatal("RequestID used untrusted overlong header")
		}
		if info.TraceID != "" {
			t.Fatalf("TraceID = %q, want invalid trace id to be dropped", info.TraceID)
		}
		if info.SpanID != "" {
			t.Fatalf("SpanID = %q, want invalid span id to be dropped", info.SpanID)
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(headerRequestID, overlongRequestID)
	req.Header.Set(headerTraceID, "trace with space")
	req.Header.Set(headerSpanID, overlongSpanID)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.NotEmpty(t, rec.Header().Get(headerRequestID))
	assert.NotEqual(t, overlongRequestID, rec.Header().Get(headerRequestID))
}

func TestRequestLoggerIncludesActiveTraceMetadata(t *testing.T) {
	tracerProvider := trace.NewTracerProvider()
	defer func() {
		assert.NoError(t, tracerProvider.Shutdown(context.Background()))
	}()

	ctx := requestctx.WithInfo(context.Background(), requestctx.Info{RequestID: "req-123"})
	base := zerolog.New(&bytes.Buffer{})
	ctx = base.WithContext(ctx)
	ctx, span := tracerProvider.Tracer("test").Start(ctx, "request")
	defer span.End()

	logger := requestLoggerFromContext(ctx)
	var buf bytes.Buffer
	logger = logger.Output(&buf)
	logger.Info().Msg("request")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("decode log event: %v", err)
	}
	if got["request_id"] != "req-123" {
		t.Fatalf("request_id = %v, want req-123", got["request_id"])
	}
	if got["trace_id"] == "" {
		t.Fatal("trace_id is empty, want active span trace id")
	}
	if got["span_id"] == "" {
		t.Fatal("span_id is empty, want active span id")
	}
}

func TestRequestLoggerOmitsTraceMetadataWithoutActiveSpan(t *testing.T) {
	ctx := requestctx.WithInfo(context.Background(), requestctx.Info{RequestID: "req-123"})
	base := zerolog.New(&bytes.Buffer{})
	ctx = base.WithContext(ctx)

	logger := requestLoggerFromContext(ctx)
	var buf bytes.Buffer
	logger = logger.Output(&buf)
	logger.Info().Msg("request")

	var got map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("decode log event: %v", err)
	}
	if _, ok := got["trace_id"]; ok {
		t.Fatal("trace_id present without active span")
	}
	if _, ok := got["span_id"]; ok {
		t.Fatal("span_id present without active span")
	}
}

func TestJWTStoresSubjectAsAuthenticatedUserID(t *testing.T) {
	const secret = "test-secret"

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: "user-123",
	}).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	jwtMiddleware, err := JWT(CreateJWTConfig(secret, nil, true))
	if err != nil {
		t.Fatalf("JWT() error = %v", err)
	}

	e := echo.New()
	e.Use(RequestContext())
	e.Use(jwtMiddleware)
	e.GET("/me", func(c *echo.Context) error {
		if got := requestctx.GetUserID(c.Request().Context()); got != "user-123" {
			t.Fatalf("GetUserID() = %q, want user-123", got)
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+signedToken)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestJWTMissingTokenReturnsGenericUnauthorized(t *testing.T) {
	jwtMiddleware, err := JWT(CreateJWTConfig("test-secret", nil, true))
	if err != nil {
		t.Fatalf("JWT() error = %v", err)
	}

	e := echo.New()
	e.HTTPErrorHandler = ErrorHandler
	e.Use(RequestContext())
	e.Use(jwtMiddleware)
	e.GET("/me", func(c *echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assertErrorPayload(t, rec, apperr.ErrUnauthorized, "Unauthorized")
}

func TestErrorRecovery(t *testing.T) {
	e := echo.New()
	e.Use(ErrorRecovery())
	e.HTTPErrorHandler = ErrorHandler

	// 创建一个会panic的路由
	e.GET("/panic", func(_ *echo.Context) error {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	// 测试panic恢复
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assertErrorPayload(t, rec, apperr.ErrInternalServer, "Internal server error")
}

func TestErrorHandlerNormalizesApplicationAndBadRequestErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   int
		wantMsg    string
	}{
		{
			name:       "plain error becomes internal server error",
			err:        errors.New("database password leaked in raw error"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   apperr.ErrInternalServer,
			wantMsg:    "Internal server error",
		},
		{
			name:       "echo bad request uses generic public message",
			err:        echo.NewHTTPError(http.StatusBadRequest, "invalid query"),
			wantStatus: http.StatusBadRequest,
			wantCode:   apperr.ErrBadRequest,
			wantMsg:    "Bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertErrorHandlerNormalizes(t, tt.err, tt.wantStatus, tt.wantCode, tt.wantMsg)
		})
	}
}

func TestErrorHandlerNormalizesEchoHTTPClientErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   int
		wantMsg    string
	}{
		{
			name:       "echo method not allowed keeps client status",
			err:        echo.NewHTTPError(http.StatusMethodNotAllowed, "method not allowed"),
			wantStatus: http.StatusMethodNotAllowed,
			wantCode:   apperr.ErrMethodNotAllowed,
			wantMsg:    "Method not allowed",
		},
		{
			name:       "echo status coder method not allowed keeps client status",
			err:        echo.ErrMethodNotAllowed,
			wantStatus: http.StatusMethodNotAllowed,
			wantCode:   apperr.ErrMethodNotAllowed,
			wantMsg:    "Method not allowed",
		},
		{
			name:       "echo conflict keeps conflict status",
			err:        echo.NewHTTPError(http.StatusConflict, "already exists"),
			wantStatus: http.StatusConflict,
			wantCode:   apperr.ErrConflict,
			wantMsg:    "Conflict",
		},
		{
			name:       "unknown echo client error remains client error",
			err:        echo.NewHTTPError(http.StatusUnsupportedMediaType, "unsupported media type"),
			wantStatus: http.StatusBadRequest,
			wantCode:   apperr.ErrBadRequest,
			wantMsg:    "Bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertErrorHandlerNormalizes(t, tt.err, tt.wantStatus, tt.wantCode, tt.wantMsg)
		})
	}
}

func TestErrorHandlerSkipsCommittedResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Response().WriteHeader(http.StatusAccepted)

	ErrorHandler(c, errors.New("late error after response"))

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Empty(t, rec.Body.String())
}

func TestValidatorReturnsFieldValidationError(t *testing.T) {
	type request struct {
		Email string `validate:"required,email"`
	}

	err := (&Validator{Validator: validator.New()}).Validate(request{})
	if err == nil {
		t.Fatal("Validate() error = nil, want validation error")
	}

	appErr, ok := apperr.Parse(err)
	if !ok {
		t.Fatal("Validate() did not return application error")
	}
	if appErr.Code() != apperr.ErrValidation {
		t.Fatalf("Code = %d, want %d", appErr.Code(), apperr.ErrValidation)
	}
	if appErr.Message() != "Email is required" {
		t.Fatalf("Message = %q, want Email is required", appErr.Message())
	}
	if appErr.Detail() == "" {
		t.Fatal("Detail is empty, want validation detail")
	}
}

func TestValidatorReturnsInternalErrorWhenMisconfigured(t *testing.T) {
	err := (&Validator{}).Validate(struct{}{})
	if err == nil {
		t.Fatal("Validate() error = nil, want misconfiguration error")
	}

	appErr, ok := apperr.Parse(err)
	if !ok {
		t.Fatal("Validate() did not return application error")
	}
	if appErr.Code() != apperr.ErrInternalServer {
		t.Fatalf("Code = %d, want %d", appErr.Code(), apperr.ErrInternalServer)
	}
	if appErr.Message() != "Internal server error" {
		t.Fatalf("Message = %q, want Internal server error", appErr.Message())
	}
}

func TestNilValidatorReturnsInternalErrorWhenMisconfigured(t *testing.T) {
	var cv *Validator

	err := cv.Validate(struct{}{})
	if err == nil {
		t.Fatal("Validate() error = nil, want misconfiguration error")
	}
}

func assertErrorPayload(t *testing.T, rec *httptest.ResponseRecorder, wantCode int, wantMessage string) {
	t.Helper()

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if got["code"] != float64(wantCode) {
		t.Fatalf("code = %v, want %d", got["code"], wantCode)
	}
	if got["message"] != wantMessage {
		t.Fatalf("message = %v, want %q", got["message"], wantMessage)
	}
}

func assertErrorHandlerNormalizes(
	t *testing.T,
	err error,
	wantStatus int,
	wantCode int,
	wantMsg string,
) {
	t.Helper()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ErrorHandler(c, err)

	assert.Equal(t, wantStatus, rec.Code)
	assertErrorPayload(t, rec, wantCode, wantMsg)
}

func TestJWTConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *JWTConfig
	}{
		{
			name: "enabled JWT",
			config: &JWTConfig{
				SigningKey: []byte("test-secret"),
				SkipPaths:  []string{"/api/health"},
				Enabled:    true,
			},
		},
		{
			name: "disabled JWT",
			config: &JWTConfig{
				SigningKey: []byte("test-secret"),
				SkipPaths:  []string{},
				Enabled:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtMiddleware, err := JWT(tt.config)
			assert.NoError(t, err)
			assert.NotNil(t, jwtMiddleware)
		})
	}
}

func TestJWTRequiresSigningKeyWhenEnabled(t *testing.T) {
	jwtMiddleware, err := JWT(&JWTConfig{Enabled: true})
	assert.Error(t, err)
	assert.Nil(t, jwtMiddleware)
}

func TestApplyMiddlewaresRequiresCORSOriginsWhenEnabled(t *testing.T) {
	e := echo.New()
	config := DefaultMiddlewareConfig()
	config.EnableCORS = true

	err := ApplyMiddlewares(e, config)

	assert.Error(t, err)
}

func TestApplyMiddlewaresRejectsWildcardCORSOriginWithCredentials(t *testing.T) {
	e := echo.New()
	config := DefaultMiddlewareConfig()
	config.EnableCORS = true
	config.CORS.AllowOrigins = []string{"*"}
	config.CORS.AllowCredentials = true

	err := ApplyMiddlewares(e, config)

	assert.Error(t, err)
}

func TestRequestPathOmitsRawQuery(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users?token=secret", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users")

	if got := requestPath(c); got != "/users" {
		t.Fatalf("requestPath() = %q, want /users", got)
	}
}

func TestMiddlewareConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *MiddlewareConfig
	}{
		{
			name:   "default config",
			config: DefaultMiddlewareConfig(),
		},
		{
			name: "custom config",
			config: &MiddlewareConfig{
				EnableRecovery:       true,
				EnableRequestContext: true,
				EnableLogger:         true,
				EnableGzip:           true,
				EnableCORS:           true,
				CORS: middleware.CORSConfig{
					AllowOrigins: []string{"https://example.com"},
				},
				EnableJWT: false,
				JWT:       DefaultJWTConfig(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			assert.NoError(t, ApplyMiddlewares(e, tt.config))
			assert.NotNil(t, e)
		})
	}
}

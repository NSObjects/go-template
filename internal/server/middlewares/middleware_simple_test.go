package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/apperr"
	"github.com/NSObjects/go-template/internal/requestctx"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
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

func TestRequestContextMiddlewareStoresMetadata(t *testing.T) {
	e := echo.New()
	e.Use(RequestContext())
	e.GET("/ping", func(c echo.Context) error {
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

func TestErrorRecovery(t *testing.T) {
	e := echo.New()
	e.Use(ErrorRecovery())
	e.HTTPErrorHandler = ErrorHandler

	// 创建一个会panic的路由
	e.GET("/panic", func(_ echo.Context) error {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	// 测试panic恢复
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assertErrorPayload(t, rec, apperr.ErrInternalServer, "Internal server error")
}

func TestErrorHandlerNormalizesErrors(t *testing.T) {
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
		{
			name:       "validation error becomes validation code",
			err:        &ValidationError{Field: "email", Message: "invalid format"},
			wantStatus: http.StatusBadRequest,
			wantCode:   apperr.ErrValidation,
			wantMsg:    "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			ErrorHandler(tt.err, c)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assertErrorPayload(t, rec, tt.wantCode, tt.wantMsg)
		})
	}
}

func TestValidatorReturnsSafeBadRequest(t *testing.T) {
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
	if appErr.Code() != apperr.ErrBadRequest {
		t.Fatalf("Code = %d, want %d", appErr.Code(), apperr.ErrBadRequest)
	}
	if appErr.Message() != "invalid request" {
		t.Fatalf("Message = %q, want invalid request", appErr.Message())
	}
	if appErr.Detail() == "invalid request" {
		t.Fatal("Detail lost validator cause")
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

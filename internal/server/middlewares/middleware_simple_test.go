package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDefaultMiddlewareConfig(t *testing.T) {
	config := DefaultMiddlewareConfig()

	assert.NotNil(t, config)
	assert.True(t, config.EnableRecovery)
	assert.True(t, config.EnableLogger)
	assert.True(t, config.EnableGzip)
	assert.True(t, config.EnableCORS)
	assert.False(t, config.EnableJWT)
	assert.NotEmpty(t, config.LoggerFormat)
	assert.NotNil(t, config.JWT)
}

func TestDefaultJWTConfig(t *testing.T) {
	config := DefaultJWTConfig()

	assert.NotNil(t, config)
	assert.False(t, config.Enabled)
	assert.NotNil(t, config.SkipPaths)
	assert.NotEmpty(t, config.SigningKey) // 默认配置有默认密钥
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

	// 测试中间件应用不会panic
	assert.NotPanics(t, func() {
		ApplyMiddlewares(e, config)
	})

	assert.NotNil(t, e)
}

func TestErrorRecovery(t *testing.T) {
	e := echo.New()
	e.Use(ErrorRecovery())
	e.HTTPErrorHandler = ErrorHandler

	// 创建一个会panic的路由
	e.GET("/panic", func(c echo.Context) error {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	// 测试panic恢复
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assertErrorPayload(t, rec, code.ErrInternalServer, "Internal server error")
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
			wantCode:   code.ErrInternalServer,
			wantMsg:    "Internal server error",
		},
		{
			name:       "echo bad request becomes bad request code",
			err:        echo.NewHTTPError(http.StatusBadRequest, "invalid query"),
			wantStatus: http.StatusBadRequest,
			wantCode:   code.ErrBadRequest,
			wantMsg:    "invalid query",
		},
		{
			name:       "validation error becomes validation code",
			err:        &ValidationError{Field: "email", Message: "invalid format"},
			wantStatus: http.StatusBadRequest,
			wantCode:   code.ErrValidation,
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
		name     string
		config   *JWTConfig
		expected bool
	}{
		{
			name: "enabled JWT",
			config: &JWTConfig{
				SigningKey: []byte("test-secret"),
				SkipPaths:  []string{"/api/health"},
				Enabled:    true,
			},
			expected: true,
		},
		{
			name: "disabled JWT",
			config: &JWTConfig{
				SigningKey: []byte("test-secret"),
				SkipPaths:  []string{},
				Enabled:    false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := JWT(tt.config)
			assert.NotNil(t, middleware)
		})
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
				EnableRecovery: true,
				EnableLogger:   true,
				EnableGzip:     true,
				EnableCORS:     true,
				EnableJWT:      false,
				LoggerFormat:   "custom format",
				JWT:            DefaultJWTConfig(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ApplyMiddlewares(e, tt.config)
			assert.NotNil(t, e)
		})
	}
}

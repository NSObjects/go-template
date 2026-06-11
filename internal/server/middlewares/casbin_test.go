package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/utils"
	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const testCasbinModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`

func TestCasbinMiddlewareSkipsConfiguredPaths(t *testing.T) {
	enforcer := newTestCasbinEnforcer(t)
	e := echo.New()
	e.Use(Casbin(enforcer, CreateCasbinConfig(true, []string{"/health"}, nil)))
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestCasbinMiddlewareAllowsConfiguredAdminUser(t *testing.T) {
	enforcer := newTestCasbinEnforcer(t)
	e := echo.New()
	e.Use(setJWTUser(&utils.JwtCustomClaims{ID: 42, Admin: true}))
	e.Use(Casbin(enforcer, CreateCasbinConfig(true, nil, []string{"root"})))
	e.GET("/admin", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestCasbinMiddlewareRejectsMissingPolicy(t *testing.T) {
	enforcer := newTestCasbinEnforcer(t)
	e := echo.New()
	e.HTTPErrorHandler = ErrorHandler
	e.Use(setJWTUser(&utils.JwtCustomClaims{ID: 42}))
	e.Use(Casbin(enforcer, CreateCasbinConfig(true, nil, []string{"root"})))
	e.GET("/admin", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func newTestCasbinEnforcer(t *testing.T) *casbin.Enforcer {
	t.Helper()

	m, err := model.NewModelFromString(testCasbinModel)
	if err != nil {
		t.Fatalf("new casbin model: %v", err)
	}

	enforcer, err := casbin.NewEnforcer(m)
	if err != nil {
		t.Fatalf("new casbin enforcer: %v", err)
	}

	return enforcer
}

func setJWTUser(claims *utils.JwtCustomClaims) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", jwt.NewWithClaims(jwt.SigningMethodHS256, claims))
			return next(c)
		}
	}
}

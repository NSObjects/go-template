package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func TestControllerCreateBindsJSONAndReturnsOperateSuccess(t *testing.T) {
	controller := newController(controllerUseCaseStub{
		t: t,
		create: func(ctx context.Context, req CreateRequest) error {
			if ctx == nil {
				t.Fatal("context is nil")
			}
			if req.Username != "lintao" || req.Email != "lintao@example.com" || req.Age != 18 {
				t.Fatalf("request = %+v, want create payload", req)
			}
			return nil
		},
	})
	ctx, rec := newControllerTestContext(
		http.MethodPost,
		"/users",
		`{"username":"lintao","email":"lintao@example.com","age":18}`,
	)

	if err := controller.Create(ctx); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	body := decodeControllerJSONBody(t, rec)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body["code"] != float64(0) {
		t.Fatalf("code = %v, want 0", body["code"])
	}
	if body["msg"] != "success" {
		t.Fatalf("msg = %v, want success", body["msg"])
	}
}

type controllerUseCaseStub struct {
	t      *testing.T
	create func(context.Context, CreateRequest) error
	UseCase
}

func (s controllerUseCaseStub) Create(ctx context.Context, req CreateRequest) error {
	s.t.Helper()
	if s.create == nil {
		s.t.Fatal("unexpected Create call")
	}
	return s.create(ctx, req)
}

func newControllerTestContext(method, target, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = &middlewares.Validator{Validator: validator.New()}
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func decodeControllerJSONBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	return body
}

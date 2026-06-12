package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/configs"
	platformapp "github.com/NSObjects/go-template/internal/platform/app"
)

func TestExistingUserRouteReachableThroughModuleFirstPath(t *testing.T) {
	cfg := configs.Config{}
	modules, err := Modules(cfg)
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}
	app, err := platformapp.Assemble(platformapp.Options{
		Config:               cfg,
		Modules:              modules.Modules,
		CapabilitySelections: modules.CapabilitySelections,
	})
	if err != nil {
		t.Fatalf("platform app Assemble() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	app.Server().Server().ServeHTTP(rec, req)

	body := decodeResponseBody(t, rec)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", rec.Code, rec.Body.String())
	}
	if body["code"] != float64(0) {
		t.Fatalf("code = %v, want 0", body["code"])
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("data = %T, want object", body["data"])
	}
	if _, ok := data["list"].([]any); !ok {
		t.Fatalf("data.list = %T, want list", data["list"])
	}
	if _, ok := data["total"].(float64); !ok {
		t.Fatalf("data.total = %T, want number", data["total"])
	}
}

func TestInvalidUserInputUsesPlatformErrorEnvelope(t *testing.T) {
	cfg := configs.Config{}
	modules, err := Modules(cfg)
	if err != nil {
		t.Fatalf("Modules() error = %v", err)
	}
	app, err := platformapp.Assemble(platformapp.Options{
		Config:               cfg,
		Modules:              modules.Modules,
		CapabilitySelections: modules.CapabilitySelections,
	})
	if err != nil {
		t.Fatalf("platform app Assemble() error = %v", err)
	}

	body := strings.NewReader(`{"username":"li","email":"bad"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	app.Server().Server().ServeHTTP(rec, req)

	response := decodeResponseBody(t, rec)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body.String())
	}
	if response["code"] != float64(code.ErrValidation) {
		t.Fatalf("code = %v, want %d", response["code"], code.ErrValidation)
	}
	if _, ok := response["message"].(string); !ok {
		t.Fatalf("message = %T, want string", response["message"])
	}
	if _, ok := response["request_id"].(string); !ok {
		t.Fatalf("request_id = %T, want string", response["request_id"])
	}
	if _, exists := response["data"]; exists {
		t.Fatalf("data exists in error envelope: %+v", response["data"])
	}
}

func decodeResponseBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v, body = %s", err, rec.Body.String())
	}
	return body
}

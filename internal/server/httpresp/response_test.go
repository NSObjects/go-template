package httpresp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/apperr"
	"github.com/NSObjects/go-template/internal/requestctx"
	"github.com/labstack/echo/v4"
)

func GetContext() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestApiError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantErr    bool
		wantStatus int
		wantCode   float64
		wantMsg    string
	}{
		{
			name:       "plain error is rendered as unknown internal error",
			err:        errors.New("api error"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   float64(apperr.ErrUnknown),
			wantMsg:    "Internal server error",
		},
		{
			name:    "nil error returns error",
			err:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := GetContext()
			if err := APIError(c, tt.err); (err != nil) != tt.wantErr {
				t.Errorf("APIError() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			assertErrorResponse(t, rec, tt.wantStatus, tt.wantCode, tt.wantMsg)
		})
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name string
		kind apperr.Kind
		want int
	}{
		{name: "bad request", kind: apperr.KindBadRequest, want: http.StatusBadRequest},
		{name: "validation", kind: apperr.KindValidation, want: http.StatusBadRequest},
		{name: "unauthorized", kind: apperr.KindUnauthorized, want: http.StatusUnauthorized},
		{name: "forbidden", kind: apperr.KindForbidden, want: http.StatusForbidden},
		{name: "not found", kind: apperr.KindNotFound, want: http.StatusNotFound},
		{name: "conflict", kind: apperr.KindConflict, want: http.StatusConflict},
		{name: "internal", kind: apperr.KindInternal, want: http.StatusInternalServerError},
		{name: "unknown kind", kind: apperr.Kind("unknown"), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Status(tt.kind); got != tt.want {
				t.Fatalf("Status(%q) = %d, want %d", tt.kind, got, tt.want)
			}
		})
	}
}

func TestRequestIDUsesRequestContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(requestctx.WithInfo(context.Background(), requestctx.Info{
		RequestID: "req-from-context",
	}))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if got := RequestID(c); got != "req-from-context" {
		t.Fatalf("RequestID() = %q, want req-from-context", got)
	}
}

func TestRequestIDPrefersRequestContextOverHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "req-from-header")
	req = req.WithContext(requestctx.WithInfo(context.Background(), requestctx.Info{
		RequestID: "req-from-context",
	}))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if got := RequestID(c); got != "req-from-context" {
		t.Fatalf("RequestID() = %q, want req-from-context", got)
	}
}

func assertErrorResponse(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantCode float64, wantMessage string) {
	t.Helper()

	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d", rec.Code, wantStatus)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if got["code"] != wantCode {
		t.Fatalf("code = %v, want %v", got["code"], wantCode)
	}
	if got["message"] != wantMessage {
		t.Fatalf("message = %v, want %v", got["message"], wantMessage)
	}
	if got["request_id"] == "" {
		t.Fatal("request_id is empty")
	}
	if _, ok := got["timestamp"].(float64); !ok {
		t.Fatalf("timestamp = %T, want number", got["timestamp"])
	}
}

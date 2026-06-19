package httpresp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/requestctx"
)

const requestIDFromContext = "req-from-context"

func GetContext() (*echo.Context, *httptest.ResponseRecorder) {
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

func TestAPIErrorSkipsCommittedResponse(t *testing.T) {
	c, rec := GetContext()
	c.Response().WriteHeader(http.StatusAccepted)

	if err := APIError(c, errors.New("late error after response")); err != nil {
		t.Fatalf("APIError() error = %v, want nil for committed response", err)
	}

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if rec.Body.String() != "" {
		t.Fatalf("body = %q, want empty", rec.Body.String())
	}
}

func TestOKResponseIncludesStandardEnvelope(t *testing.T) {
	c, rec := GetContext()
	c.SetRequest(c.Request().WithContext(requestctx.WithInfo(context.Background(), requestctx.Info{
		RequestID: requestIDFromContext,
	})))

	if err := OK(c, map[string]string{"id": "order-1"}); err != nil {
		t.Fatalf("OK() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if got["code"] != float64(apperr.ErrSuccess) {
		t.Fatalf("code = %v, want success", got["code"])
	}
	if got["message"] != "OK" {
		t.Fatalf("message = %v, want OK", got["message"])
	}
	if got["request_id"] != requestIDFromContext {
		t.Fatalf("request_id = %v, want %s", got["request_id"], requestIDFromContext)
	}
	if got["data"] == nil {
		t.Fatal("data is nil")
	}
}

func TestListResponseIncludesPaginationMetadata(t *testing.T) {
	c, rec := GetContext()
	meta, err := NewPageMeta(2, 10, 25)
	if err != nil {
		t.Fatalf("NewPageMeta() error = %v", err)
	}

	if err := List(c, []string{"a", "b"}, meta); err != nil {
		t.Fatalf("List() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	page, ok := got["page"].(map[string]any)
	if !ok {
		t.Fatalf("page = %T, want object", got["page"])
	}
	if page["page"] != float64(2) {
		t.Fatalf("page.page = %v, want 2", page["page"])
	}
	if page["page_size"] != float64(10) {
		t.Fatalf("page.page_size = %v, want 10", page["page_size"])
	}
	if page["total"] != float64(25) {
		t.Fatalf("page.total = %v, want 25", page["total"])
	}
	if page["has_next"] != true {
		t.Fatalf("page.has_next = %v, want true", page["has_next"])
	}
}

func TestNewPageMetaRejectsInvalidPage(t *testing.T) {
	_, err := NewPageMeta(0, 10, 25)
	if err == nil {
		t.Fatal("NewPageMeta() error = nil, want invalid pagination error")
	}
	appErr, ok := apperr.Parse(err)
	if !ok {
		t.Fatal("NewPageMeta() did not return application error")
	}
	if appErr.Code() != apperr.ErrBadRequest {
		t.Fatalf("Code = %d, want %d", appErr.Code(), apperr.ErrBadRequest)
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
		{name: "method not allowed", kind: apperr.KindMethodNotAllowed, want: http.StatusMethodNotAllowed},
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
		RequestID: requestIDFromContext,
	}))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if got := RequestID(c); got != requestIDFromContext {
		t.Fatalf("RequestID() = %q, want %s", got, requestIDFromContext)
	}
	if got := rec.Header().Get("X-Request-ID"); got != requestIDFromContext {
		t.Fatalf("response X-Request-ID = %q, want %s", got, requestIDFromContext)
	}
}

func TestRequestIDPrefersRequestContextOverHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "req-from-header")
	req = req.WithContext(requestctx.WithInfo(context.Background(), requestctx.Info{
		RequestID: requestIDFromContext,
	}))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if got := RequestID(c); got != requestIDFromContext {
		t.Fatalf("RequestID() = %q, want %s", got, requestIDFromContext)
	}
}

func TestRequestIDRejectsInvalidHeaderFallback(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "req with spaces")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	got := RequestID(c)
	if got == "" {
		t.Fatal("RequestID() = empty, want generated request id")
	}
	if got == "req with spaces" {
		t.Fatal("RequestID() used invalid request header")
	}
	if rec.Header().Get("X-Request-ID") != got {
		t.Fatalf("response request id = %q, want %q", rec.Header().Get("X-Request-ID"), got)
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

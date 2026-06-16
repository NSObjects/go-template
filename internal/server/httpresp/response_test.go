/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package httpresp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NSObjects/go-template/internal/apperr"
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

func TestOperateSuccess(t *testing.T) {
	c, rec := GetContext()

	if err := OperateSuccess(c); err != nil {
		t.Fatalf("OperateSuccess() error = %v", err)
	}

	assertJSONResponse(t, rec, http.StatusOK, map[string]any{
		"code": float64(0),
		"msg":  "success",
		"data": map[string]any{},
	})
}

func TestListDataResponse(t *testing.T) {
	tests := []struct {
		name  string
		list  interface{}
		total int64
		want  map[string]any
	}{
		{
			name:  "nil list returns empty list",
			list:  nil,
			total: 0,
			want: map[string]any{
				"code": float64(0),
				"msg":  "success",
				"data": map[string]any{
					"list":  []any{},
					"total": float64(0),
				},
			},
		},
		{
			name:  "typed nil slice returns empty list",
			list:  []string(nil),
			total: 0,
			want: map[string]any{
				"code": float64(0),
				"msg":  "success",
				"data": map[string]any{
					"list":  []any{},
					"total": float64(0),
				},
			},
		},
		{
			name:  "list data preserves list and total",
			list:  []string{"alpha"},
			total: 1,
			want: map[string]any{
				"code": float64(0),
				"msg":  "success",
				"data": map[string]any{
					"list":  []any{"alpha"},
					"total": float64(1),
				},
			},
		},
		{
			name:  "non nilable list value is preserved",
			list:  struct{ Name string }{Name: "alpha"},
			total: 1,
			want: map[string]any{
				"code": float64(0),
				"msg":  "success",
				"data": map[string]any{
					"list": map[string]any{
						"Name": "alpha",
					},
					"total": float64(1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := GetContext()

			if err := ListDataResponse(c, tt.list, tt.total); err != nil {
				t.Fatalf("ListDataResponse() error = %v", err)
			}

			assertJSONResponse(t, rec, http.StatusOK, tt.want)
		})
	}
}

func TestOneDataResponse(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want map[string]any
	}{
		{
			name: "nil data is preserved",
			data: nil,
			want: map[string]any{
				"code": float64(0),
				"msg":  "success",
				"data": nil,
			},
		},
		{
			name: "object data is wrapped",
			data: map[string]string{"name": "lintao"},
			want: map[string]any{
				"code": float64(0),
				"msg":  "success",
				"data": map[string]any{"name": "lintao"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := GetContext()

			if err := OneDataResponse(c, tt.data); err != nil {
				t.Fatalf("OneDataResponse() error = %v", err)
			}

			assertJSONResponse(t, rec, http.StatusOK, tt.want)
		})
	}
}

func assertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, want map[string]any) {
	t.Helper()

	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d", rec.Code, wantStatus)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	assertMapEqual(t, got, want)
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

func assertMapEqual(t *testing.T, got, want map[string]any) {
	t.Helper()

	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("marshal got response: %v", err)
	}

	wantJSON, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal want response: %v", err)
	}

	if string(gotJSON) != string(wantJSON) {
		t.Fatalf("response = %s, want %s", gotJSON, wantJSON)
	}
}

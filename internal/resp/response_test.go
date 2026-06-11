/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{err: errors.New("api error")},
			wantErr: false,
		},
		{
			args:    args{err: nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := GetContext()
			if err := APIError(c, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("APIError() error = %v, wantErr %v", err, tt.wantErr)
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

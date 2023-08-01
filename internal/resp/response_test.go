/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package resp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func GetContext() (c echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestApiError(t *testing.T) {
	type args struct {
		err error
		c   echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{err: errors.New("api error"), c: GetContext()},
			wantErr: false,
		},
		{
			args:    args{err: nil, c: GetContext()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := APIError(tt.args.err, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("APIError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperateSuccess(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    struct{ c echo.Context }{c: GetContext()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OperateSuccess(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("OperateSuccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListDataResponse(t *testing.T) {
	type args struct {
		arr   interface{}
		total int64
		c     echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: struct {
				arr   interface{}
				total int64
				c     echo.Context
			}{arr: nil, total: 0, c: GetContext()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ListDataResponse(tt.args.arr, tt.args.total, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ListDataResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOneDataResponse(t *testing.T) {
	type args struct {
		data interface{}
		c    echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: struct {
				data interface{}
				c    echo.Context
			}{data: nil, c: GetContext()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := OneDataResponse(tt.args.data, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("OneDataResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

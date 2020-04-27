/*
 *
 * api_error.go
 * api_helper
 *
 * Created by lintao on 2019-01-26 16:58
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package tools

import (
	"errors"
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		this *Error
		want string
	}{
		{
			this: NewError(errors.New("test error"), StatusOK),
			want: "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.this.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	type args struct {
		err  error
		code ResponetsCode
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			args: struct {
				err  error
				code ResponetsCode
			}{err: errors.New("some error"), code: StatusDBErr},
			want: NewError(errors.New("some error"), StatusDBErr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.err, tt.args.code); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewParamError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			args: struct{ err error }{err: errors.New("param error")},
			want: NewParamError(errors.New("param error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParamError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParamError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDBError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			args: args{err: errors.New("db error")},
			want: NewDBError(errors.New("db error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDBError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMsgError(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			args: args{str: "str error"},
			want: NewMsgError("str error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMsgError(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMsgError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReloginError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			args: args{err: errors.New("relogin error")},
			want: NewReloginError(errors.New("relogin error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReloginError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReloginError() = %v, want %v", got, tt.want)
			}
		})
	}
}

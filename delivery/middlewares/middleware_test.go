/*
 *
 * middleware.go
 * api_helper
 *
 * Created by lintao on 2019-01-29 09:18
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package middlewares

import (
	"testing"

	"gopkg.in/go-playground/validator.v9"
)

type vData struct {
	Id    int    `validate:"gte=0,lte=130"`
	Name  string `validate:"required"`
	Phone string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestValidator_Validate(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		cv      *Validator
		args    args
		wantErr bool
	}{
		{
			cv:      &Validator{Validator: validator.New()},
			args:    struct{ i interface{} }{i: vData{Id: 1, Name: "string", Phone: "string", Email: "8888@qq.com"}},
			wantErr: false,
		},
		{
			cv:      &Validator{Validator: validator.New()},
			args:    struct{ i interface{} }{i: vData{Id: 131, Name: "string", Phone: "string", Email: "email"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cv.Validate(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("Validator.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

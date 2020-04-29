/*
 *
 * mongodb_test.go
 * db
 *
 * Created by lintao on 2020/4/27 11:33 上午
 * Copyright © 2020-2020 LINTAO. All rights reserved.
 *
 */

package db

import (
	"go-template/tools/configs"
	"testing"
)

func Test_mongoUrl(t *testing.T) {

	if err := configs.InitConfig("", "toml"); err != nil {
		panic(err)
	}
	tests := []struct {
		name string
		want string
	}{
		{
			name: "测试url",
			want: "mongodb://localhost:27017",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mongoUrl(); got != tt.want {
				t.Errorf("mongoUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

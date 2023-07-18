/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"github.com/NSObjects/go-template/tools/configs"
	"testing"
)

func Test_mongoUrl(t *testing.T) {

	if err := configs.InitConfig("toml"); err != nil {
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

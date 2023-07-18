/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package configs

import "testing"

func TestInitConfig(t *testing.T) {
	type args struct {
		configPath string
		configType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitConfig(tt.args.configPath); (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_viperInit(t *testing.T) {
	type args struct {
		configPath string
		configType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := viperInit(tt.args.configPath); (err != nil) != tt.wantErr {
				t.Errorf("viperInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunEnvironment(t *testing.T) {
	tests := []struct {
		name string
		want EnvironmentType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RunEnvironment(); got != tt.want {
				t.Errorf("RunEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}

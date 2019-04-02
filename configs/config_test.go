/*
 *
 * config.go
 * configs
 *
 * Created by lin on 2018/12/10 4:31 PM
 * Copyright © 2017-2018 PYL. All rights reserved.
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
			if err := InitConfig(tt.args.configPath, tt.args.configType); (err != nil) != tt.wantErr {
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
			if err := viperInit(tt.args.configPath, tt.args.configType); (err != nil) != tt.wantErr {
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

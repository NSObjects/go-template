/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package tools

import "testing"

func TestMd5Encode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantMd  string
		wantErr bool
	}{
		{
			args:    struct{ str string }{str: "str"},
			wantMd:  "341be97d9aff90c9978347f66f945b77",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMd, err := Md5Encode(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("Md5Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMd != tt.wantMd {
				t.Errorf("Md5Encode() = %v, want %v", gotMd, tt.wantMd)
			}
		})
	}
}

func TestCoin_Yuan(t *testing.T) {
	tests := []struct {
		name string
		c    Coin
		want float64
	}{
		{
			c:    ToCoin(38.99),
			want: 38.99,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Yuan(); got != tt.want {
				t.Errorf("Coin.Yuan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCoin(t *testing.T) {
	type args struct {
		yuan float64
	}
	tests := []struct {
		name string
		args args
		want Coin
	}{
		{
			args: struct{ yuan float64 }{yuan: 39.99},
			want: Coin(39.99 * Ratio),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCoin(tt.args.yuan); got != tt.want {
				t.Errorf("ToCoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRound(t *testing.T) {
	type args struct {
		x   float64
		pre int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: struct {
				x   float64
				pre int
			}{x: 39.99, pre: 0},
			want: 40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Round(tt.args.x, tt.args.pre); got != tt.want {
				t.Errorf("Round() = %v, want %v", got, tt.want)
			}
		})
	}
}

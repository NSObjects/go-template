/*
 *
 * db.go
 * db
 *
 * Created by lin on 2018-12-26 17:21
 * Copyright © 2017-2018 PYL. All rights reserved.
 *
 */

package db

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func TestNewDataSource(t *testing.T) {
	tests := []struct {
		name          string
		wantDatasouce *DataSource
		wantErr       bool
	}{
		{
			wantDatasouce: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDatasouce, err := NewDataSource()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDataSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDatasouce, tt.wantDatasouce) {
				t.Errorf("NewDataSource() = %v, want %v", gotDatasouce, tt.wantDatasouce)
			}
		})
	}
}

func TestNewTestDataSource(t *testing.T) {
	type args struct {
		db2 *sql.DB
	}
	tests := []struct {
		name          string
		args          args
		wantDatasouce *DataSource
		wantErr       bool
	}{
		{
			wantDatasouce: nil,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDatasouce, err := NewTestDataSource(tt.args.db2)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestDataSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDatasouce, tt.wantDatasouce) {
				t.Errorf("NewTestDataSource() = %v, want %v", gotDatasouce, tt.wantDatasouce)
			}
		})
	}
}

func Test_createEngin(t *testing.T) {
	type args struct {
		db  DataBase
		db2 []*sql.DB
	}
	tests := []struct {
		name       string
		args       args
		wantEngine *xorm.EngineGroup
		wantErr    bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEngine, err := createEngin(tt.args.db, tt.args.db2...)
			if (err != nil) != tt.wantErr {
				t.Errorf("createEngin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEngine, tt.wantEngine) {
				t.Errorf("createEngin() = %v, want %v", gotEngine, tt.wantEngine)
			}
		})
	}
}

func Test_dataSource(t *testing.T) {
	type args struct {
		db DataBase
	}
	tests := []struct {
		name     string
		args     args
		wantConn []string
		wantErr  bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConn, err := dataSource(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("dataSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConn, tt.wantConn) {
				t.Errorf("dataSource() = %v, want %v", gotConn, tt.wantConn)
			}
		})
	}
}

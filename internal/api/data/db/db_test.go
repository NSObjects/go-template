/*
 * Created by lintao on 2023/7/18 下午3:56
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package db

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func TestNewDataSource(t *testing.T) {
	tests := []struct {
		name           string
		wantDatasource *DataSource
		wantErr        bool
	}{
		{
			wantDatasource: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDatasource, err := NewDataSource()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDataSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDatasource, tt.wantDatasource) {
				t.Errorf("NewDataSource() = %v, want %v", gotDatasource, tt.wantDatasource)
			}
		})
	}
}

func TestNewTestDataSource(t *testing.T) {
	type args struct {
		db2 *sql.DB
	}
	tests := []struct {
		name           string
		args           args
		wantDatasource *DataSource
		wantErr        bool
	}{
		{
			wantDatasource: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDatasource, err := NewTestDataSource(tt.args.db2)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTestDataSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDatasource, tt.wantDatasource) {
				t.Errorf("NewTestDataSource() = %v, want %v", gotDatasource, tt.wantDatasource)
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

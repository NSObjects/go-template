/*
 *
 * user.go
 * repository
 *
 * Created by lintao on 2019-01-29 16:29
 * Copyright © 2017-2019 PYL. All rights reserved.
 *
 */

package repository

import (
	"fmt"
	"go-template/domain"
	"go-template/tools/db"
	"reflect"
	"regexp"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewUserDataSource(t *testing.T) {

	dbMock, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = dbMock.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	dataSource, err := db.NewTestDataSource(dbMock)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		dataSource *db.DataSource
	}

	tests := []struct {
		name string
		args args
		want *UserDataSource
	}{
		{
			args: args{dataSource: dataSource},
			want: NewUserDataSource(dataSource),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserDataSource(tt.args.dataSource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserDataSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDataSource_CreateUser(t *testing.T) {

	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = dbMock.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	mock.ExpectExec(
		"INSERT INTO `user`").
		WithArgs("string", "string", 1, "string", "string", nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	database, err := db.NewTestDataSource(dbMock)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		dataSource *db.DataSource
	}
	type args struct {
		param domain.UserParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  int64
		wantErr bool
	}{
		{
			fields: fields{dataSource: database},
			args: args{param: struct {
				domain.ApiQuery
				domain.User
			}{User: domain.User{
				Name:     "string",
				Account:  "string",
				Password: "string",
				Phone:    "string",
				Status:   1,
			}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := NewUserDataSource(tt.fields.dataSource)
			gotId, err := this.CreateUser(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDataSource.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("UserDataSource.CreateUser() = %v, want %v", gotId, tt.wantId)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserDataSource_GetUserById(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = dbMock.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows := sqlmock.NewRows([]string{"id", "name", "account", "password", "phone", "status", "created"}).
		AddRow(1, "string", "string", "string", "string", 0, nil)
	sql := regexp.QuoteMeta("SELECT `id`, `name`, `phone`, `status`, `account`, `password`, `created` FROM `user` WHERE `id`=? LIMIT 1")
	mock.ExpectQuery(sql).WithArgs(1).WillReturnRows(rows)

	database, err := db.NewTestDataSource(dbMock)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		dataSource *db.DataSource
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser domain.User
		wantErr  bool
	}{
		{
			fields:   fields{dataSource: database},
			args:     args{id: 1},
			wantUser: domain.User{Id: 1, Name: "string", Account: "string", Password: "string", Phone: "string", Status: 0},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := NewUserDataSource(tt.fields.dataSource)
			gotUser, err := this.GetUserById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDataSource.GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("UserDataSource.GetUserById() = %v, want %v", gotUser, tt.wantUser)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserDataSource_FindUser(t *testing.T) {

	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = dbMock.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows := sqlmock.NewRows([]string{"count(*)"}).AddRow(20)
	mock.ExpectQuery("^SELECT count(.+)$").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "name", "account", "password", "phone", "status", "created"}).
		AddRow(1, "string", "string", "string", "string", 0, nil)
	sql := regexp.QuoteMeta("SELECT `id`, `name`, `phone`, `status`, `account`, `password`, `created` FROM `user` LIMIT 20")

	mock.ExpectQuery(sql).WillReturnRows(rows)
	database, err := db.NewTestDataSource(dbMock)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		dataSource *db.DataSource
	}
	type args struct {
		param domain.UserParam
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantUser  []domain.User
		wantTotal int64
		wantErr   bool
	}{
		{
			fields:    fields{dataSource: database},
			args:      args{param: domain.UserParam{}},
			wantErr:   false,
			wantUser:  []domain.User{{Id: 1, Name: "string", Account: "string", Password: "string", Phone: "string", Status: 0}},
			wantTotal: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &UserDataSource{
				dataSource: tt.fields.dataSource,
			}
			gotUser, gotTotal, err := this.FindUser(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDataSource.FindUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("UserDataSource.FindUser() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("UserDataSource.FindUser() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserDataSource_DeleteUserById(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = dbMock.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	mock.ExpectExec("DELETE FROM `user` WHERE `id`=?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	database, err := db.NewTestDataSource(dbMock)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		dataSource *db.DataSource
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields:  fields{dataSource: database},
			args:    args{id: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &UserDataSource{
				dataSource: tt.fields.dataSource,
			}
			if err := this.DeleteUserById(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserDataSource.DeleteUserById() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserDataSource_UpdateUser(t *testing.T) {

	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	if err = dbMock.Close(); err != nil {
	//		t.Fatal(err)
	//	}
	//}()

	mock.ExpectExec("UPDATE `user`").
		WithArgs("lin", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dataSource, err := db.NewTestDataSource(dbMock)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		dataSource *db.DataSource
	}
	type args struct {
		param domain.UserParam
		id    int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields:  fields{dataSource: dataSource},
			args:    args{param: domain.UserParam{User: domain.User{Name: "lin"}}, id: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			this := NewUserDataSource(dataSource)
			if err := this.UpdateUser(tt.args.param, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserDataSource.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

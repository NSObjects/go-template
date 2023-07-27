/*
 * Created by lintao on 2023/7/18 下午3:59
 * Copyright © 2020-2023 LINTAO. All rights reserved.
 *
 */

package biz

//
//func TestUserHandler_GetUser(t *testing.T) {
//
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//	pmock := mock_repository.NewMockUserRepository(ctl)
//	pmock.EXPECT().FindUser(domain.UserParam{APIQuery: domain.APIQuery{Page: 1, Count: 10}})
//
//	type fields struct {
//		repository repository.UserRepository
//	}
//	type args struct {
//		c echo.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			args:    args{c: GetUsersContext()},
//			fields:  fields{repository: pmock},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := NewUserHandler(tt.fields.repository)
//			if err := u.ListUser(tt.args.c); (err != nil) != tt.wantErr {
//				t.Errorf("UserHandler.ListUser() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func GetUsersContext() (c echo.Context) {
//	e := echo.New()
//	q := make(url.Values)
//	q.Set("page", "1")
//	q.Set("count", "10")
//	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	return e.NewContext(req, rec)
//}
//
//func CreatedUserContext() echo.Context {
//	var user = `{
//    "id": 1,
//    "name": "lin",
//    "phone": "string",
//    "account": "string",
//    "password": "string",
//    "status": 1
//	}`
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(user))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	e.Validator = &middlewares.Validator{Validator: validator.New()}
//	return e.NewContext(req, rec)
//}
//
//func TestUserHandler_CreateUser(t *testing.T) {
//
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//	pmock := mock_repository.NewMockUserRepository(ctl)
//	pmock.EXPECT().CreateUser(domain.UserParam{
//		User: domain.User{Id: 1, Name: "lin", Phone: "string", Account: "string", Password: "string", Status: 1}})
//	type fields struct {
//		repository repository.UserRepository
//	}
//	type args struct {
//		c echo.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			args:    args{c: CreatedUserContext()},
//			fields:  fields{repository: pmock},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := NewUserHandler(tt.fields.repository)
//			if err := u.CreateUser(tt.args.c); (err != nil) != tt.wantErr {
//				t.Errorf("UserHandler.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func DeleteUserContext() echo.Context {
//
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodDelete, "/", nil)
//	rec := httptest.NewRecorder()
//
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	e.Validator = &middlewares.Validator{Validator: validator.New()}
//	c := e.NewContext(req, rec)
//
//	c.SetPath("/users/:id")
//	c.SetParamNames("id")
//	c.SetParamValues("1")
//	return c
//}
//
//func TestUserHandler_DeleteUser(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//	pmock := mock_repository.NewMockUserRepository(ctl)
//	pmock.EXPECT().DeleteUserByID(int64(1))
//
//	type fields struct {
//		repository repository.UserRepository
//	}
//	type args struct {
//		c echo.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			args:    args{c: DeleteUserContext()},
//			fields:  fields{repository: pmock},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := NewUserHandler(tt.fields.repository)
//			if err := u.DeleteUser(tt.args.c); (err != nil) != tt.wantErr {
//				t.Errorf("UserHandler.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func UpdateUserContext() echo.Context {
//	var user = `{
//    "id": 1,
//    "name": "lin",
//    "phone": "string",
//    "account": "string",
//    "password": "string",
//    "status": 1
//	}`
//	e := echo.New()
//	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(user))
//	rec := httptest.NewRecorder()
//
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	e.Validator = &middlewares.Validator{Validator: validator.New()}
//	c := e.NewContext(req, rec)
//
//	c.SetPath("/users/:id")
//	c.SetParamNames("id")
//	c.SetParamValues("1")
//	return c
//}
//
//func TestUserHandler_UpdateUser(t *testing.T) {
//	ctl := gomock.NewController(t)
//	defer ctl.Finish()
//	pmock := mock_repository.NewMockUserRepository(ctl)
//	pmock.EXPECT().UpdateUser(domain.UserParam{
//		User: domain.User{Id: 1, Name: "lin", Phone: "string", Account: "string", Password: "string", Status: 1}}, int64(1))
//	type fields struct {
//		repository repository.UserRepository
//	}
//	type args struct {
//		c echo.Context
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			args:    args{c: UpdateUserContext()},
//			fields:  fields{repository: pmock},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := NewUserHandler(tt.fields.repository)
//			if err := u.UpdateUser(tt.args.c); (err != nil) != tt.wantErr {
//				t.Errorf("UserHandler.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

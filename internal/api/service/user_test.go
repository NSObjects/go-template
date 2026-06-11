package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
	"github.com/NSObjects/go-template/internal/server/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	merrors "github.com/marmotedu/errors"
)

type userUseCaseStub struct {
	t *testing.T

	listUsers func(context.Context, param.UserListUsersRequest) ([]param.UserListItem, int64, error)
	create    func(context.Context, param.UserCreateRequest) error
	getByID   func(context.Context, int64) (param.UserData, error)
	update    func(context.Context, int64, param.UserUpdateRequest) error
	delete    func(context.Context, int64) error
}

func (s userUseCaseStub) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	s.t.Helper()
	if s.listUsers == nil {
		s.t.Fatal("unexpected ListUsers call")
	}
	return s.listUsers(ctx, req)
}

func (s userUseCaseStub) Create(ctx context.Context, req param.UserCreateRequest) error {
	s.t.Helper()
	if s.create == nil {
		s.t.Fatal("unexpected Create call")
	}
	return s.create(ctx, req)
}

func (s userUseCaseStub) GetByID(ctx context.Context, id int64) (param.UserData, error) {
	s.t.Helper()
	if s.getByID == nil {
		s.t.Fatal("unexpected GetByID call")
	}
	return s.getByID(ctx, id)
}

func (s userUseCaseStub) Update(ctx context.Context, id int64, req param.UserUpdateRequest) error {
	s.t.Helper()
	if s.update == nil {
		s.t.Fatal("unexpected Update call")
	}
	return s.update(ctx, id, req)
}

func (s userUseCaseStub) Delete(ctx context.Context, id int64) error {
	s.t.Helper()
	if s.delete == nil {
		s.t.Fatal("unexpected Delete call")
	}
	return s.delete(ctx, id)
}

func TestUserControllerListUsersBindsQueryAndReturnsListEnvelope(t *testing.T) {
	controller := &UserController{
		user: userUseCaseStub{
			t: t,
			listUsers: func(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
				if ctx == nil {
					t.Fatal("context is nil")
				}
				if req.Page != 2 || req.Count != 5 {
					t.Fatalf("request = %+v, want page=2 count=5", req)
				}
				return []param.UserListItem{{Id: 7, Username: "lintao", Email: "lintao@example.com"}}, 1, nil
			},
		},
	}
	ctx, rec := newServiceTestContext(http.MethodGet, "/users?page=2&count=5", "")

	if err := controller.ListUsers(ctx); err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}

	body := decodeJSONBody(t, rec)
	assertSuccessEnvelope(t, rec, body)
	data := body["data"].(map[string]interface{})
	if data["total"] != float64(1) {
		t.Fatalf("total = %v, want 1", data["total"])
	}
	list := data["list"].([]interface{})
	if len(list) != 1 {
		t.Fatalf("list length = %d, want 1", len(list))
	}
	item := list[0].(map[string]interface{})
	if item["id"] != float64(7) || item["username"] != "lintao" {
		t.Fatalf("list item = %+v, want user id=7 username=lintao", item)
	}
}

func TestUserControllerCreateBindsJSONAndReturnsOperateSuccess(t *testing.T) {
	controller := &UserController{
		user: userUseCaseStub{
			t: t,
			create: func(ctx context.Context, req param.UserCreateRequest) error {
				if req.Username != "lintao" || req.Email != "lintao@example.com" || req.Age != 18 {
					t.Fatalf("request = %+v, want create payload", req)
				}
				return nil
			},
		},
	}
	ctx, rec := newServiceTestContext(
		http.MethodPost,
		"/users",
		`{"username":"lintao","email":"lintao@example.com","age":18}`,
	)

	if err := controller.Create(ctx); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	body := decodeJSONBody(t, rec)
	assertSuccessEnvelope(t, rec, body)
	data := body["data"].(map[string]interface{})
	if len(data) != 0 {
		t.Fatalf("data = %+v, want empty object", data)
	}
}

func TestUserControllerCreateRejectsInvalidJSONBeforeUseCase(t *testing.T) {
	controller := &UserController{user: userUseCaseStub{t: t}}
	ctx, _ := newServiceTestContext(http.MethodPost, "/users", `{"username":`)

	err := controller.Create(ctx)
	assertCoder(t, err, code.ErrBind)
}

func TestUserControllerCreateRejectsInvalidPayloadBeforeUseCase(t *testing.T) {
	controller := &UserController{user: userUseCaseStub{t: t}}
	ctx, _ := newServiceTestContext(http.MethodPost, "/users", `{"username":"li","email":"bad"}`)

	err := controller.Create(ctx)
	assertCoder(t, err, code.ErrValidation)
}

func TestUserControllerGetByIDParsesPathParamAndReturnsData(t *testing.T) {
	controller := &UserController{
		user: userUseCaseStub{
			t: t,
			getByID: func(ctx context.Context, id int64) (param.UserData, error) {
				if id != 42 {
					t.Fatalf("id = %d, want 42", id)
				}
				return param.UserData{Id: 42, Username: "lintao", Email: "lintao@example.com"}, nil
			},
		},
	}
	ctx, rec := newServiceTestContext(http.MethodGet, "/users/42", "")
	setIDParam(ctx, "42")

	if err := controller.GetByID(ctx); err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	body := decodeJSONBody(t, rec)
	assertSuccessEnvelope(t, rec, body)
	data := body["data"].(map[string]interface{})
	if data["id"] != float64(42) || data["username"] != "lintao" {
		t.Fatalf("data = %+v, want user id=42 username=lintao", data)
	}
}

func TestUserControllerGetByIDRejectsInvalidPathParamBeforeUseCase(t *testing.T) {
	controller := &UserController{user: userUseCaseStub{t: t}}
	ctx, _ := newServiceTestContext(http.MethodGet, "/users/not-a-number", "")
	setIDParam(ctx, "not-a-number")

	err := controller.GetByID(ctx)
	assertCoder(t, err, code.ErrBadRequest)
}

func TestUserControllerUpdateParsesPathParamAndReturnsOperateSuccess(t *testing.T) {
	controller := &UserController{
		user: userUseCaseStub{
			t: t,
			update: func(ctx context.Context, id int64, req param.UserUpdateRequest) error {
				if id != 42 {
					t.Fatalf("id = %d, want 42", id)
				}
				if req.Email != "new@example.com" {
					t.Fatalf("request = %+v, want updated email", req)
				}
				return nil
			},
		},
	}
	ctx, rec := newServiceTestContext(http.MethodPut, "/users/42", `{"email":"new@example.com"}`)
	setIDParam(ctx, "42")

	if err := controller.Update(ctx); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	body := decodeJSONBody(t, rec)
	assertSuccessEnvelope(t, rec, body)
}

func TestUserControllerDeleteParsesPathParamAndReturnsOperateSuccess(t *testing.T) {
	controller := &UserController{
		user: userUseCaseStub{
			t: t,
			delete: func(ctx context.Context, id int64) error {
				if id != 42 {
					t.Fatalf("id = %d, want 42", id)
				}
				return nil
			},
		},
	}
	ctx, rec := newServiceTestContext(http.MethodDelete, "/users/42", "")
	setIDParam(ctx, "42")

	if err := controller.Delete(ctx); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	body := decodeJSONBody(t, rec)
	assertSuccessEnvelope(t, rec, body)
}

func TestUserControllerRegisterRouterUsesEchoPathParams(t *testing.T) {
	controller := &UserController{
		user: userUseCaseStub{
			t: t,
			delete: func(ctx context.Context, id int64) error {
				if id != 42 {
					t.Fatalf("id = %d, want 42", id)
				}
				return nil
			},
		},
	}
	e := newServiceTestEcho()
	controller.RegisterRouter(e.Group("/api"))

	req := httptest.NewRequest(http.MethodDelete, "/api/users/42", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	body := decodeJSONBody(t, rec)
	assertSuccessEnvelope(t, rec, body)
}

func TestUserControllerReturnsUseCaseError(t *testing.T) {
	wantErr := errors.New("usecase failed")
	controller := &UserController{
		user: userUseCaseStub{
			t: t,
			create: func(ctx context.Context, req param.UserCreateRequest) error {
				return wantErr
			},
		},
	}
	ctx, _ := newServiceTestContext(
		http.MethodPost,
		"/users",
		`{"username":"lintao","email":"lintao@example.com","age":18}`,
	)

	err := controller.Create(ctx)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Create() error = %v, want %v", err, wantErr)
	}
}

func newServiceTestContext(method, target, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := newServiceTestEcho()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func newServiceTestEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &middlewares.Validator{Validator: validator.New()}
	return e
}

func setIDParam(ctx echo.Context, id string) {
	ctx.SetParamNames("id")
	ctx.SetParamValues(id)
}

func decodeJSONBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	return body
}

func assertSuccessEnvelope(t *testing.T, rec *httptest.ResponseRecorder, body map[string]interface{}) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body["code"] != float64(0) {
		t.Fatalf("code = %v, want 0", body["code"])
	}
	if body["msg"] != "success" {
		t.Fatalf("msg = %v, want success", body["msg"])
	}
	if _, ok := body["data"]; !ok {
		t.Fatal("response missing data field")
	}
}

func assertCoder(t *testing.T, err error, wantCode int) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want code %d", wantCode)
	}
	coder := merrors.ParseCoder(err)
	if coder == nil {
		t.Fatalf("error = %v has no coder, want code %d", err, wantCode)
	}
	if coder.Code() != wantCode {
		t.Fatalf("code = %d, want %d", coder.Code(), wantCode)
	}
}

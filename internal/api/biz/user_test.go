package biz

import (
	"context"
	"errors"
	"testing"

	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
)

func TestUserHandlerDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := &fakeUserRepository{
		list: []param.UserListItem{{Id: 1, Username: "alice", Email: "alice@example.com", Age: 20}},
		data: param.UserData{Id: 1, Username: "alice", Email: "alice@example.com", Age: 20},
	}
	handler := NewUserHandler(repo)

	list, total, err := handler.ListUsers(ctx, param.UserListUsersRequest{APIQuery: param.APIQuery{Page: 1, Count: 10}})
	if err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Id != 1 {
		t.Fatalf("ListUsers() = (%#v, %d), want one user with total 1", list, total)
	}

	if err := handler.Create(ctx, param.UserCreateRequest{Username: "alice", Email: "alice@example.com", Age: 20}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	data, err := handler.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if data.Id != 1 {
		t.Fatalf("GetByID() id = %d, want 1", data.Id)
	}

	if err := handler.Update(ctx, 1, param.UserUpdateRequest{Username: "alice"}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := handler.Delete(ctx, 1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	wantCalls := []string{"ListUsers", "Create", "GetByID", "Update", "Delete"}
	if len(repo.calls) != len(wantCalls) {
		t.Fatalf("repo calls = %v, want %v", repo.calls, wantCalls)
	}
	for i, want := range wantCalls {
		if repo.calls[i] != want {
			t.Fatalf("repo calls[%d] = %s, want %s", i, repo.calls[i], want)
		}
	}
}

func TestUserHandlerPassesThroughRepositoryErrors(t *testing.T) {
	repoErr := code.WrapNotFoundError(errors.New("repository failed"), "user not found")
	handler := NewUserHandler(&fakeUserRepository{err: repoErr})

	_, _, err := handler.ListUsers(context.Background(), param.UserListUsersRequest{})
	if err == nil {
		t.Fatal("ListUsers() error = nil, want repository error")
	}
	if !errors.Is(err, repoErr) {
		t.Fatalf("ListUsers() error does not wrap repository error: %v", err)
	}
	coder, ok := code.ParseRegisteredCoder(err)
	if !ok {
		t.Fatalf("ListUsers() error has no code: %v", err)
	}
	if coder.Code() != code.ErrNotFound {
		t.Fatalf("ListUsers() error code = %d, want %d", coder.Code(), code.ErrNotFound)
	}
}

type fakeUserRepository struct {
	list  []param.UserListItem
	data  param.UserData
	err   error
	calls []string
}

func (f *fakeUserRepository) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	f.calls = append(f.calls, "ListUsers")
	if f.err != nil {
		return nil, 0, f.err
	}
	return f.list, int64(len(f.list)), nil
}

func (f *fakeUserRepository) Create(ctx context.Context, req param.UserCreateRequest) error {
	f.calls = append(f.calls, "Create")
	return f.err
}

func (f *fakeUserRepository) GetByID(ctx context.Context, id int64) (param.UserData, error) {
	f.calls = append(f.calls, "GetByID")
	if f.err != nil {
		return param.UserData{}, f.err
	}
	return f.data, nil
}

func (f *fakeUserRepository) Update(ctx context.Context, id int64, req param.UserUpdateRequest) error {
	f.calls = append(f.calls, "Update")
	return f.err
}

func (f *fakeUserRepository) Delete(ctx context.Context, id int64) error {
	f.calls = append(f.calls, "Delete")
	return f.err
}

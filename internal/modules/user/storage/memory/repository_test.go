package memory

import (
	"context"
	"testing"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
)

func TestRepositorySupportsUserLifecycle(t *testing.T) {
	var _ biz.UserRepository = NewRepository()

	repo := NewRepository()
	ctx := context.Background()

	if err := repo.Create(ctx, param.UserCreateRequest{
		Username: "lintao",
		Email:    "lintao@example.com",
		Age:      18,
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	list, total, err := repo.ListUsers(ctx, param.UserListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("ListUsers() total/list = %d/%d, want 1/1", total, len(list))
	}
	created := list[0]
	if created.Id == 0 || created.Username != "lintao" || created.Email != "lintao@example.com" || created.Age != 18 {
		t.Fatalf("created user = %+v, want lintao user data", created)
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatalf("created timestamps = %v/%v, want non-zero", created.CreatedAt, created.UpdatedAt)
	}

	got, err := repo.GetByID(ctx, created.Id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Id != created.Id || got.Username != created.Username || got.Email != created.Email {
		t.Fatalf("GetByID() = %+v, want created user", got)
	}

	if err := repo.Update(ctx, created.Id, param.UserUpdateRequest{
		Email: "new@example.com",
		Age:   19,
	}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	updated, err := repo.GetByID(ctx, created.Id)
	if err != nil {
		t.Fatalf("GetByID() after update error = %v", err)
	}
	if updated.Username != "lintao" || updated.Email != "new@example.com" || updated.Age != 19 {
		t.Fatalf("updated user = %+v, want patched data preserving username", updated)
	}
	if updated.UpdatedAt.Before(updated.CreatedAt) {
		t.Fatalf("updated timestamps = %v/%v, want updated_at >= created_at", updated.UpdatedAt, updated.CreatedAt)
	}

	if err := repo.Delete(ctx, created.Id); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	list, total, err = repo.ListUsers(ctx, param.UserListUsersRequest{})
	if err != nil {
		t.Fatalf("ListUsers() after delete error = %v", err)
	}
	if total != 0 || len(list) != 0 {
		t.Fatalf("ListUsers() after delete total/list = %d/%d, want 0/0", total, len(list))
	}
}

func TestRepositoryReturnsNotFoundForMissingUser(t *testing.T) {
	repo := NewRepository()

	_, err := repo.GetByID(context.Background(), 404)
	assertNotFound(t, err)

	err = repo.Update(context.Background(), 404, param.UserUpdateRequest{Email: "missing@example.com"})
	assertNotFound(t, err)

	err = repo.Delete(context.Background(), 404)
	assertNotFound(t, err)
}

func assertNotFound(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("error = nil, want not found")
	}
	coder, ok := code.ParseRegisteredCoder(err)
	if !ok {
		t.Fatalf("error has no registered code: %v", err)
	}
	if coder.Code() != code.ErrNotFound {
		t.Fatalf("error code = %d, want %d", coder.Code(), code.ErrNotFound)
	}
}

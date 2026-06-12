package user

import (
	"context"
	"errors"
	"testing"

	"github.com/NSObjects/go-template/internal/code"
)

func TestUseCaseDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRepository{
		list: []ListItem{{Id: 1, Username: "alice", Email: "alice@example.com", Age: 20}},
		data: Data{Id: 1, Username: "alice", Email: "alice@example.com", Age: 20},
	}
	useCase := NewUseCase(repo)

	list, total, err := useCase.ListUsers(ctx, ListUsersRequest{Query: Query{Page: 1, Count: 10}})
	if err != nil {
		t.Fatalf("ListUsers() error = %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Id != 1 {
		t.Fatalf("ListUsers() = (%#v, %d), want one user with total 1", list, total)
	}

	if err := useCase.Create(ctx, CreateRequest{Username: "alice", Email: "alice@example.com", Age: 20}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	data, err := useCase.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if data.Id != 1 {
		t.Fatalf("GetByID() id = %d, want 1", data.Id)
	}

	if err := useCase.Update(ctx, 1, UpdateRequest{Username: "alice"}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := useCase.Delete(ctx, 1); err != nil {
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

func TestUseCasePassesThroughRepositoryErrors(t *testing.T) {
	repoErr := code.WrapNotFoundError(errors.New("repository failed"), "user not found")
	useCase := NewUseCase(&fakeRepository{err: repoErr})

	_, _, err := useCase.ListUsers(context.Background(), ListUsersRequest{})
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

type fakeRepository struct {
	list  []ListItem
	data  Data
	err   error
	calls []string
}

func (f *fakeRepository) ListUsers(context.Context, ListUsersRequest) ([]ListItem, int64, error) {
	f.calls = append(f.calls, "ListUsers")
	if f.err != nil {
		return nil, 0, f.err
	}
	return f.list, int64(len(f.list)), nil
}

func (f *fakeRepository) Create(context.Context, CreateRequest) error {
	f.calls = append(f.calls, "Create")
	return f.err
}

func (f *fakeRepository) GetByID(context.Context, int64) (Data, error) {
	f.calls = append(f.calls, "GetByID")
	if f.err != nil {
		return Data{}, f.err
	}
	return f.data, nil
}

func (f *fakeRepository) Update(context.Context, int64, UpdateRequest) error {
	f.calls = append(f.calls, "Update")
	return f.err
}

func (f *fakeRepository) Delete(context.Context, int64) error {
	f.calls = append(f.calls, "Delete")
	return f.err
}

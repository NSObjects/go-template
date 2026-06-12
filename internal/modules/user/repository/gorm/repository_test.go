package gorm

import (
	"context"
	"testing"

	"github.com/NSObjects/go-template/internal/code"
	user "github.com/NSObjects/go-template/internal/modules/user"
)

func TestRepositoryReturnsDatabaseErrorWhenQueryIsMissing(t *testing.T) {
	repo := New(nil)
	ctx := context.Background()

	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "list users",
			run: func() error {
				_, _, err := repo.ListUsers(ctx, user.ListUsersRequest{})
				return err
			},
		},
		{
			name: "create user",
			run: func() error {
				return repo.Create(ctx, user.CreateRequest{})
			},
		},
		{
			name: "get user by id",
			run: func() error {
				_, err := repo.GetByID(ctx, 1)
				return err
			},
		},
		{
			name: "update user",
			run: func() error {
				return repo.Update(ctx, 1, user.UpdateRequest{})
			},
		},
		{
			name: "delete user",
			run: func() error {
				return repo.Delete(ctx, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if err == nil {
				t.Fatal("repository method error = nil, want database error")
			}
			coder, ok := code.ParseRegisteredCoder(err)
			if !ok {
				t.Fatalf("repository method error has no code: %v", err)
			}
			if coder.Code() != code.ErrDatabase {
				t.Fatalf("repository method error code = %d, want %d", coder.Code(), code.ErrDatabase)
			}
		})
	}
}

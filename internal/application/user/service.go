package user

import (
	"context"

	"github.com/NSObjects/go-template/internal/code"
	domain "github.com/NSObjects/go-template/internal/domain/user"
)

// Service defines the application use cases for the user aggregate.
type Service interface {
	ListUsers(ctx context.Context, query domain.ListUsersQuery) ([]domain.User, int64, error)
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id domain.ID) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id domain.ID) error
}

type service struct {
	repo domain.Repository
}

// NewService constructs a user Service implementation.
func NewService(repo domain.Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListUsers(ctx context.Context, query domain.ListUsersQuery) ([]domain.User, int64, error) {
	list, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, 0, code.WrapDatabaseError(err, "查询用户列表失败")
	}
	return list, total, nil
}

func (s *service) Create(ctx context.Context, user domain.User) error {
	if err := s.repo.Create(ctx, user); err != nil {
		return code.WrapDatabaseError(err, "创建用户失败")
	}
	return nil
}

func (s *service) GetByID(ctx context.Context, id domain.ID) (domain.User, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return result, code.WrapDatabaseError(err, "查询用户详情失败")
	}
	return result, nil
}

func (s *service) Update(ctx context.Context, user domain.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return code.WrapDatabaseError(err, "更新用户失败")
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id domain.ID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return code.WrapDatabaseError(err, "删除用户失败")
	}
	return nil
}

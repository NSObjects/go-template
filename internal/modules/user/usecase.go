package user

import "context"

// Repository stores and retrieves users for the user module.
type Repository interface {
	ListUsers(ctx context.Context, req ListUsersRequest) ([]ListItem, int64, error)
	Create(ctx context.Context, req CreateRequest) error
	GetByID(ctx context.Context, id int64) (Data, error)
	Update(ctx context.Context, id int64, req UpdateRequest) error
	Delete(ctx context.Context, id int64) error
}

// UseCase contains the user module's business operations.
type UseCase interface {
	ListUsers(ctx context.Context, req ListUsersRequest) ([]ListItem, int64, error)
	Create(ctx context.Context, req CreateRequest) error
	GetByID(ctx context.Context, id int64) (Data, error)
	Update(ctx context.Context, id int64, req UpdateRequest) error
	Delete(ctx context.Context, id int64) error
}

type useCase struct {
	repository Repository
}

// NewUseCase creates the user module use case backed by repository.
func NewUseCase(repository Repository) UseCase {
	return &useCase{repository: repository}
}

func (u *useCase) ListUsers(ctx context.Context, req ListUsersRequest) ([]ListItem, int64, error) {
	return u.repository.ListUsers(ctx, req)
}

func (u *useCase) Create(ctx context.Context, req CreateRequest) error {
	return u.repository.Create(ctx, req)
}

func (u *useCase) GetByID(ctx context.Context, id int64) (Data, error) {
	return u.repository.GetByID(ctx, id)
}

func (u *useCase) Update(ctx context.Context, id int64, req UpdateRequest) error {
	return u.repository.Update(ctx, id, req)
}

func (u *useCase) Delete(ctx context.Context, id int64) error {
	return u.repository.Delete(ctx, id)
}

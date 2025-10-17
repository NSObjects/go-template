package user

import "context"

// Repository defines the storage port for the user aggregate.
type Repository interface {
	List(ctx context.Context, query ListUsersQuery) ([]User, int64, error)
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id ID) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id ID) error
}

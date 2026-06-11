package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/code"
)

type Repository struct {
	mu     sync.RWMutex
	nextID int64
	users  map[int64]param.UserData
}

// NewRepository creates an empty in-memory user repository.
func NewRepository() *Repository {
	return &Repository{
		nextID: 1,
		users:  make(map[int64]param.UserData),
	}
}

func (r *Repository) ListUsers(ctx context.Context, req param.UserListUsersRequest) ([]param.UserListItem, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]param.UserData, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	total := int64(len(users))
	limit := req.Limit()
	offset := req.Offset()
	if offset > len(users) {
		offset = len(users)
	}
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	list := make([]param.UserListItem, 0, end-offset)
	for _, user := range users[offset:end] {
		list = append(list, toListItem(user))
	}
	return list, total, nil
}

func (r *Repository) Create(ctx context.Context, req param.UserCreateRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	id := r.nextID
	r.nextID++
	r.users[id] = param.UserData{
		Id:        id,
		Username:  req.Username,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (param.UserData, error) {
	if err := ctx.Err(); err != nil {
		return param.UserData{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return param.UserData{}, userNotFound()
	}
	return user, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req param.UserUpdateRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return userNotFound()
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Age != 0 {
		user.Age = req.Age
	}
	user.UpdatedAt = time.Now()
	r.users[id] = user
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return userNotFound()
	}
	delete(r.users, id)
	return nil
}

func toListItem(user param.UserData) param.UserListItem {
	return param.UserListItem{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func userNotFound() error {
	return code.WrapNotFoundError(errors.New("user not found"), "user not found")
}

var _ biz.UserRepository = (*Repository)(nil)

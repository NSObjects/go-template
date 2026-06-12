package memory

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/NSObjects/go-template/internal/code"
	user "github.com/NSObjects/go-template/internal/modules/user"
)

type Repository struct {
	mu     sync.RWMutex
	nextID int64
	users  map[int64]user.Data
}

// NewRepository creates an empty in-memory user repository.
func NewRepository() *Repository {
	return &Repository{
		nextID: 1,
		users:  make(map[int64]user.Data),
	}
}

func (r *Repository) ListUsers(ctx context.Context, req user.ListUsersRequest) ([]user.ListItem, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]user.Data, 0, len(r.users))
	for _, item := range r.users {
		users = append(users, item)
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

	list := make([]user.ListItem, 0, end-offset)
	for _, item := range users[offset:end] {
		list = append(list, toListItem(item))
	}
	return list, total, nil
}

func (r *Repository) Create(ctx context.Context, req user.CreateRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	id := r.nextID
	r.nextID++
	r.users[id] = user.Data{
		Id:        id,
		Username:  req.Username,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (user.Data, error) {
	if err := ctx.Err(); err != nil {
		return user.Data{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.users[id]
	if !ok {
		return user.Data{}, userNotFound()
	}
	return item, nil
}

func (r *Repository) Update(ctx context.Context, id int64, req user.UpdateRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.users[id]
	if !ok {
		return userNotFound()
	}
	if req.Username != "" {
		item.Username = req.Username
	}
	if req.Email != "" {
		item.Email = req.Email
	}
	if req.Age != 0 {
		item.Age = req.Age
	}
	item.UpdatedAt = time.Now()
	r.users[id] = item
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

func toListItem(data user.Data) user.ListItem {
	return user.ListItem{
		Id:        data.Id,
		Username:  data.Username,
		Email:     data.Email,
		Age:       data.Age,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func userNotFound() error {
	return code.WrapNotFoundError(errors.New("user not found"), "user not found")
}

var _ user.Repository = (*Repository)(nil)

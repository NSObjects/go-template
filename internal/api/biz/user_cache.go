package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/NSObjects/go-template/internal/api/service/param"
	"github.com/NSObjects/go-template/internal/cache"
)

// UserCacheService 用户缓存服务
type UserCacheService struct {
	cache *cache.RedisCache
}

// NewUserCacheService 创建用户缓存服务
func NewUserCacheService(cache *cache.RedisCache) *UserCacheService {
	return &UserCacheService{cache: cache}
}

// GetUserByID 从缓存获取用户
func (s *UserCacheService) GetUserByID(ctx context.Context, id int64) (*param.UserData, error) {
	key := fmt.Sprintf("user:%d", id)
	var user param.UserData
	err := s.cache.Get(ctx, key, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SetUser 设置用户缓存
func (s *UserCacheService) SetUser(ctx context.Context, user *param.UserData) error {
	key := fmt.Sprintf("user:%d", user.Id)
	return s.cache.Set(ctx, key, user, 30*time.Minute)
}

// DeleteUser 删除用户缓存
func (s *UserCacheService) DeleteUser(ctx context.Context, id int64) error {
	key := fmt.Sprintf("user:%d", id)
	return s.cache.Delete(ctx, key)
}

// InvalidateUserList 使用户列表缓存失效
func (s *UserCacheService) InvalidateUserList(ctx context.Context) error {
	pattern := "user_list:*"
	return s.cache.DeletePattern(ctx, pattern)
}

// GetUserList 获取用户列表缓存
func (s *UserCacheService) GetUserList(ctx context.Context, page, size int) ([]param.UserListItem, int64, error) {
	key := fmt.Sprintf("user_list:page_%d_size_%d", page, size)
	var result struct {
		Users []param.UserListItem `json:"users"`
		Total int64                `json:"total"`
	}
	err := s.cache.Get(ctx, key, &result)
	if err != nil {
		return nil, 0, err
	}
	return result.Users, result.Total, nil
}

// SetUserList 设置用户列表缓存
func (s *UserCacheService) SetUserList(ctx context.Context, page, size int, users []param.UserListItem, total int64) error {
	key := fmt.Sprintf("user_list:page_%d_size_%d", page, size)
	result := struct {
		Users []param.UserListItem `json:"users"`
		Total int64                `json:"total"`
	}{
		Users: users,
		Total: total,
	}
	return s.cache.Set(ctx, key, result, 10*time.Minute)
}

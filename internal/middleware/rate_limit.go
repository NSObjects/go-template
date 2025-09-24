package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// RateLimiter 限流器
type RateLimiter struct {
	redis *redis.Client
}

// NewRateLimiter 创建限流器
func NewRateLimiter(redis *redis.Client) *RateLimiter {
	return &RateLimiter{redis: redis}
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Requests int                         // 请求次数
	Window   time.Duration               // 时间窗口
	KeyFunc  func(c echo.Context) string // 限流键生成函数
}

// RateLimit 限流中间件
func (rl *RateLimiter) RateLimit(config RateLimitConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := config.KeyFunc(c)
			if key == "" {
				return next(c)
			}

			// 使用滑动窗口算法
			allowed, err := rl.isAllowed(c.Request().Context(), key, config)
			if err != nil {
				// Redis错误时允许通过，避免影响业务
				return next(c)
			}

			if !allowed {
				return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
					"code":    429,
					"message": "请求过于频繁，请稍后再试",
				})
			}

			return next(c)
		}
	}
}

// isAllowed 检查是否允许请求
func (rl *RateLimiter) isAllowed(ctx context.Context, key string, config RateLimitConfig) (bool, error) {
	now := time.Now()
	windowStart := now.Add(-config.Window)

	// 使用Redis的ZREMRANGEBYSCORE和ZADD实现滑动窗口
	pipe := rl.redis.Pipeline()

	// 删除过期的记录
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.Unix()))

	// 添加当前请求
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.Unix()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})

	// 获取当前窗口内的请求数
	pipe.ZCard(ctx, key)

	// 设置过期时间
	pipe.Expire(ctx, key, config.Window)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// 获取请求数
	count := results[2].(*redis.IntCmd).Val()

	return int(count) <= config.Requests, nil
}

// DefaultKeyFunc 默认的限流键生成函数
func DefaultKeyFunc(c echo.Context) string {
	// 使用IP地址作为限流键
	return fmt.Sprintf("rate_limit:%s", c.RealIP())
}

// UserKeyFunc 基于用户ID的限流键生成函数
func UserKeyFunc(c echo.Context) string {
	userID := c.Get("user_id")
	if userID == nil {
		return ""
	}
	return fmt.Sprintf("rate_limit:user:%v", userID)
}

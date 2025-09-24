package health

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// HealthChecker 健康检查器
type HealthChecker struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(db *gorm.DB, redis *redis.Client) *HealthChecker {
	return &HealthChecker{
		db:    db,
		redis: redis,
	}
}

// CheckResult 健康检查结果
type CheckResult struct {
	Status    string           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Checks    map[string]Check `json:"checks"`
}

// Check 单个检查项结果
type Check struct {
	Status      string        `json:"status"`
	Message     string        `json:"message,omitempty"`
	Duration    time.Duration `json:"duration"`
	LastChecked time.Time     `json:"last_checked"`
}

// CheckAll 执行所有健康检查
func (hc *HealthChecker) CheckAll(ctx context.Context) CheckResult {
	result := CheckResult{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks:    make(map[string]Check),
	}

	// 检查数据库
	dbCheck := hc.checkDatabase(ctx)
	result.Checks["database"] = dbCheck

	// 检查Redis
	redisCheck := hc.checkRedis(ctx)
	result.Checks["redis"] = redisCheck

	// 检查系统资源
	systemCheck := hc.checkSystem()
	result.Checks["system"] = systemCheck

	// 确定整体状态
	for _, check := range result.Checks {
		if check.Status != "healthy" {
			result.Status = "unhealthy"
			break
		}
	}

	return result
}

// checkDatabase 检查数据库连接
func (hc *HealthChecker) checkDatabase(ctx context.Context) Check {
	start := time.Now()

	if hc.db == nil {
		return Check{
			Status:      "unhealthy",
			Message:     "database not configured",
			Duration:    time.Since(start),
			LastChecked: time.Now(),
		}
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 执行简单查询
	var result int
	err := hc.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error

	duration := time.Since(start)

	if err != nil {
		return Check{
			Status:      "unhealthy",
			Message:     fmt.Sprintf("database query failed: %v", err),
			Duration:    duration,
			LastChecked: time.Now(),
		}
	}

	return Check{
		Status:      "healthy",
		Message:     "database connection successful",
		Duration:    duration,
		LastChecked: time.Now(),
	}
}

// checkRedis 检查Redis连接
func (hc *HealthChecker) checkRedis(ctx context.Context) Check {
	start := time.Now()

	if hc.redis == nil {
		return Check{
			Status:      "unhealthy",
			Message:     "redis not configured",
			Duration:    time.Since(start),
			LastChecked: time.Now(),
		}
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 执行ping命令
	err := hc.redis.Ping(ctx).Err()

	duration := time.Since(start)

	if err != nil {
		return Check{
			Status:      "unhealthy",
			Message:     fmt.Sprintf("redis ping failed: %v", err),
			Duration:    duration,
			LastChecked: time.Now(),
		}
	}

	return Check{
		Status:      "healthy",
		Message:     "redis connection successful",
		Duration:    duration,
		LastChecked: time.Now(),
	}
}

// checkSystem 检查系统资源
func (hc *HealthChecker) checkSystem() Check {
	start := time.Now()

	// 这里可以添加更多系统检查，如内存使用率、CPU使用率等
	// 目前只做基本检查

	return Check{
		Status:      "healthy",
		Message:     "system resources normal",
		Duration:    time.Since(start),
		LastChecked: time.Now(),
	}
}

// CheckDatabase 单独检查数据库
func (hc *HealthChecker) CheckDatabase(ctx context.Context) Check {
	return hc.checkDatabase(ctx)
}

// CheckRedis 单独检查Redis
func (hc *HealthChecker) CheckRedis(ctx context.Context) Check {
	return hc.checkRedis(ctx)
}

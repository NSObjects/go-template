package adapters

import (
	"context"

	"github.com/NSObjects/go-template/internal/infra/persistence"
	"github.com/NSObjects/go-template/internal/user/app"
	"gorm.io/gorm"
)

// gormTxManager GORM事务管理器实现
type gormTxManager struct {
	db *persistence.DataManager
}

// NewGormTxManager 创建GORM事务管理器
func NewGormTxManager(db *persistence.DataManager) app.TransactionManager {
	return &gormTxManager{db: db}
}

// ExecuteInTransaction 在事务中执行操作
func (tm *gormTxManager) ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.MySQLWithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建新的上下文，包含事务
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

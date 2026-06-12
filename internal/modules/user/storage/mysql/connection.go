package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NSObjects/go-template/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newMysqlConnection(cfg configs.MysqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	database, err := openMysqlConnection(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	sqlDB, err := database.DB()
	if err == nil {
		configureSQLPool(sqlDB, cfg)
	}
	return database, nil
}

func openMysqlConnection(dialector gorm.Dialector) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	return gorm.Open(dialector, &gorm.Config{Logger: newLogger})
}

func configureSQLPool(sqlDB *sql.DB, cfg configs.MysqlConfig) {
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}

package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(dsn string, driver string, maxOpenConns, maxIdleConns int, connMaxLifetime int64) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch driver {
	case "sqlite":
		// 确保数据库目录存在
		dbDir := filepath.Dir(dsn)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %v", err)
		}

		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
		}

		// 设置连接池参数
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("获取数据库连接失败: %v", err)
		}

		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetMaxIdleConns(maxIdleConns)
		if connMaxLifetime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
		}

	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", driver)
	}

	return db, nil
}

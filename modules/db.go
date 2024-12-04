package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"mind.com/log/models"
	"mind.com/log/utils"
)

func InitDatabase() (*gorm.DB, error) {
	dbConfig := utils.GetDatabaseConfig()

	var db *gorm.DB
	var err error

	switch dbConfig.Driver {
	case "sqlite":
		// 确保数据库目录存在
		dbDir := filepath.Dir(dbConfig.DSN)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %v", err)
		}

		db, err = gorm.Open(sqlite.Open(dbConfig.DSN), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
		}

	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", dbConfig.Driver)
	}

	// 初始化数据库表和基础数据
	if err := utils.InitDatabase(db); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&models.Log{}, &models.StatsRecord{}, &models.StatsInfo{}); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	return db, nil
}

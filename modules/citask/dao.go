package citask

import (
	"github.com/andycai/unitool/models"
)

// 数据迁移
func autoMigrate() error {
	return app.DB.AutoMigrate(&models.Task{}, &models.TaskLog{})
}

package note

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
	"gorm.io/gorm"
)

// 数据访问层，目前暂时没有特殊的数据访问逻辑
// 所有数据库操作都在 service 层完成

func autoMigrate() error {
	return app.DB.AutoMigrate(
		&models.NoteCategory{},
		&models.Note{},
		&models.NotePermission{},
		&models.CategoryPermission{},
	)
}

func initData() error {
	return nil

	// 检查是否已初始化
	var count int64
	app.DB.Model(&models.NoteCategory{}).Count(&count)
	if count > 0 {
		log.Println("笔记模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建笔记相关权限
		permissions := []models.Permission{
			{
				Name:        "笔记列表",
				Code:        "note:list",
				Description: "查看笔记列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建笔记",
				Code:        "note:create",
				Description: "创建新笔记",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新笔记",
				Code:        "note:update",
				Description: "更新笔记信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除笔记",
				Code:        "note:delete",
				Description: "删除笔记",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "分类列表",
				Code:        "note:category:list",
				Description: "查看笔记分类列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建分类",
				Code:        "note:category:create",
				Description: "创建笔记分类",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新分类",
				Code:        "note:category:update",
				Description: "更新笔记分类",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除分类",
				Code:        "note:category:delete",
				Description: "删除笔记分类",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		return nil
	})
}

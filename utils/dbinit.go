package utils

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mind.com/log/models"
)

// InitDatabase 初始化数据库
func InitDatabase(db *gorm.DB) error {
	// 创建表
	if err := createTables(db); err != nil {
		return err
	}

	// 初始化基础数据
	if err := initBaseData(db); err != nil {
		return err
	}

	return nil
}

// 创建表
func createTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
	)
}

// 初始化基础数据
func initBaseData(db *gorm.DB) error {
	// 检查是否已初始化
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建基础权限
		permissions := []models.Permission{
			{
				Name:        "用户列表",
				Code:        "user:list",
				Description: "查看用户列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建用户",
				Code:        "user:create",
				Description: "创建新用户",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新用户",
				Code:        "user:update",
				Description: "更新用户信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除用户",
				Code:        "user:delete",
				Description: "删除用户",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "角色列表",
				Code:        "role:list",
				Description: "查看角色列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建角色",
				Code:        "role:create",
				Description: "创建新角色",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新角色",
				Code:        "role:update",
				Description: "更新角色信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除角色",
				Code:        "role:delete",
				Description: "删除角色",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "权限列表",
				Code:        "permission:list",
				Description: "查看权限列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建权限",
				Code:        "permission:create",
				Description: "创建新权限",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新权限",
				Code:        "permission:update",
				Description: "更新权限信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除权限",
				Code:        "permission:delete",
				Description: "删除权限",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 2. 创建管理员角色
		adminRole := models.Role{
			Name:        "超级管理员",
			Description: "系统超级管理员",
			Permissions: permissions, // 赋予所有权限
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := tx.Create(&adminRole).Error; err != nil {
			return err
		}

		// 3. 创建管理员用户
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		adminUser := models.User{
			Username:  "admin",
			Password:  string(hashedPassword),
			Nickname:  "系统管理员",
			RoleID:    adminRole.ID,
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(&adminUser).Error; err != nil {
			return err
		}

		return nil
	})
}

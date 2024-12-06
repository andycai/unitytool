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
		&models.AdminLog{},
		&models.GameLog{},
		&models.StatsRecord{},
		&models.StatsInfo{},
		&models.Menu{},
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
			{
				Name:        "查看日志",
				Code:        "admin_log:list",
				Description: "查看操作日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除日志",
				Code:        "admin_log:delete",
				Description: "删除操作日志",
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

		// 4. 初始化菜单数据
		if err := initMenus(tx); err != nil {
			return err
		}

		return nil
	})
}

// initMenus 初始化菜单数据
func initMenus(tx *gorm.DB) error {
	now := time.Now()

	// 系统管理菜单组
	systemManage := models.Menu{
		ParentID:   0,
		Name:       "系统管理",
		Path:       "/admin",
		Icon:       "settings",
		Sort:       1,
		Permission: "",
		IsShow:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := tx.Create(&systemManage).Error; err != nil {
		return err
	}

	// 系统管理子菜单
	systemMenus := []models.Menu{
		{
			ParentID:   systemManage.ID,
			Name:       "用户管理",
			Path:       "/admin/users",
			Icon:       "user",
			Sort:       1,
			Permission: "user:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   systemManage.ID,
			Name:       "角色管理",
			Path:       "/admin/roles",
			Icon:       "users",
			Sort:       2,
			Permission: "role:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   systemManage.ID,
			Name:       "权限管理",
			Path:       "/admin/permissions",
			Icon:       "key",
			Sort:       3,
			Permission: "permission:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   systemManage.ID,
			Name:       "操作日志",
			Path:       "/admin/logs",
			Icon:       "list",
			Sort:       4,
			Permission: "admin_log:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := tx.Create(&systemMenus).Error; err != nil {
		return err
	}

	// 游戏管理菜单组
	gameManage := models.Menu{
		ParentID:   0,
		Name:       "游戏管理",
		Path:       "/admin/game",
		Icon:       "game",
		Sort:       2,
		Permission: "",
		IsShow:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := tx.Create(&gameManage).Error; err != nil {
		return err
	}

	// 游戏管理子菜单
	gameMenus := []models.Menu{
		{
			ParentID:   gameManage.ID,
			Name:       "游戏日志",
			Path:       "/admin/game/logs",
			Icon:       "file-text",
			Sort:       1,
			Permission: "",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   gameManage.ID,
			Name:       "性能统计",
			Path:       "/admin/game/stats",
			Icon:       "bar-chart",
			Sort:       2,
			Permission: "",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := tx.Create(&gameMenus).Error; err != nil {
		return err
	}

	// 系统工具菜单组
	toolsManage := models.Menu{
		ParentID:   0,
		Name:       "系统工具",
		Path:       "/admin/tools",
		Icon:       "tool",
		Sort:       3,
		Permission: "",
		IsShow:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := tx.Create(&toolsManage).Error; err != nil {
		return err
	}

	// 系统工具子菜单
	toolsMenus := []models.Menu{
		{
			ParentID:   toolsManage.ID,
			Name:       "文件浏览",
			Path:       "/admin/browse",
			Icon:       "folder",
			Sort:       1,
			Permission: "",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   toolsManage.ID,
			Name:       "FTP上传",
			Path:       "/admin/ftp",
			Icon:       "upload",
			Sort:       2,
			Permission: "",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   toolsManage.ID,
			Name:       "服务器配置",
			Path:       "/admin/serverconf",
			Icon:       "server",
			Sort:       3,
			Permission: "",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   toolsManage.ID,
			Name:       "命令执行",
			Path:       "/admin/cmd",
			Icon:       "terminal",
			Sort:       4,
			Permission: "",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   toolsManage.ID,
			Name:       "打包工具",
			Path:       "/admin/pack",
			Icon:       "package",
			Sort:       5,
			Permission: "",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := tx.Create(&toolsMenus).Error; err != nil {
		return err
	}

	return nil
}

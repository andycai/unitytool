package user

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

// 数据迁移
func autoMigrate() error {
	return app.DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.RolePermission{}, &models.ModuleInit{})
}

// 初始化数据
func initData() error {
	// 检查是否已初始化
	if app.IsInitializedModule("user") {
		log.Println("用户模块数据库已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建基础权限
		permissions := []models.Permission{
			// 用户
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
			// 角色
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
			// 权限
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
			// 菜单
			{
				Name:        "菜单列表",
				Code:        "menu:list",
				Description: "查看菜单列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "创建菜单",
				Code:        "menu:create",
				Description: "创建新菜单",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "更新菜单",
				Code:        "menu:update",
				Description: "更新菜单信息",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除菜单",
				Code:        "menu:delete",
				Description: "删除菜单",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			// 系统日志
			{
				Name:        "日志列表",
				Code:        "adminlog:list",
				Description: "查看操作日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除日志",
				Code:        "adminlog:delete",
				Description: "删除操作日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			// 游戏日志
			{
				Name:        "游戏日志",
				Code:        "gamelog:list",
				Description: "查看游戏日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除游戏日志",
				Code:        "gamelog:delete",
				Description: "删除游戏日志",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			// 游戏统计
			{
				Name:        "游戏统计",
				Code:        "stats:list",
				Description: "查看游戏统计",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "删除游戏统计",
				Code:        "stats:delete",
				Description: "删除游戏统计",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			// 文件浏览
			{
				Name:        "文件浏览",
				Code:        "browse:list",
				Description: "查看文件浏览",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "文件FTP上传",
				Code:        "browse:ftp",
				Description: "FTP上传文件",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "文件删除",
				Code:        "browse:delete",
				Description: "删除文件",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			// 服务器配置
			{
				Name:        "服务器配置",
				Code:        "serverconf:list",
				Description: "查看服务器配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器配置更新",
				Code:        "serverconf:update",
				Description: "更新服务器配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器配置删除",
				Code:        "serverconf:delete",
				Description: "删除服务器配置",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			// 任务管理
			{
				Name:        "构建任务查看",
				Code:        "citask:list",
				Description: "查看构建任务",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "构建任务创建",
				Code:        "citask:create",
				Description: "创建构建任务",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "构建任务更新",
				Code:        "citask:update",
				Description: "更新构建任务",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "构建任务删除",
				Code:        "citask:delete",
				Description: "删除构建任务",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "构建任务执行",
				Code:        "citask:run",
				Description: "执行构建任务",
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

		// 4. 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "user",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

package modules

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/andycai/unitool/dao"
	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/utils"
)

// ModuleConfig 模块配置接口
type ModuleConfig interface {
	IsEnabled() bool
}

// Module 模块接口
type Module interface {
	Init() error
	RegisterRoutes(app *fiber.App)
}

var adminGroup fiber.Router
var apiGroup fiber.Router

func GetAdminGroup() fiber.Router {
	return adminGroup
}

func GetApiGroup() fiber.Router {
	return apiGroup
}

// 初始化全局路由
func InitGlobalRoutes(app *fiber.App, db *gorm.DB) {
	// 登录页面路由（不需要认证）
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{}, "login")
	})

	// 登录 API 路由（不需要认证）
	app.Post("/api/login", func(c *fiber.Ctx) error {
		return handlers.Login(c, db)
	})

	// 创建管理路由组
	adminGroup = app.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(db))

	// API 路由组
	apiGroup = app.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware(db))
}

// 初始化模块
func InitModules(app *fiber.App, db *gorm.DB) {
	// 初始化菜单模块
	menuDao := dao.NewMenuDao(db)
	InitMenuModule(app, menuDao)

	// 初始化并注册模块
	moduleList := []Module{
		&GameLogsModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("game_logs"),
			},
		},
		&StatsModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("stats"),
			},
		},
		&BrowseModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("browse"),
			},
			ServerConfig: utils.GetServerConfig(),
		},
		&FTPModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("ftp"),
			},
			ServerConfig: utils.GetServerConfig(),
			FTPConfig:    utils.GetFTPConfig(),
		},
		&ServerConfModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("serverconf"),
			},
			ServerConfig: utils.GetServerConfig(),
			JSONPaths:    utils.GetJSONPathConfig(),
		},
		&CmdModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("cmd"),
			},
			ServerConfig: utils.GetServerConfig(),
		},
		&PackModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("pack"),
			},
		},
		&AuthModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("auth"),
			},
		},
		&AdminLogsModule{
			BaseModule: BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("admin_logs"),
			},
		},
	}

	// 初始化和注册所有模块
	for _, module := range moduleList {
		if err := module.Init(); err != nil {
			log.Printf("模块初始化失败: %v", err)
			continue
		}
		module.RegisterRoutes(app)
	}
}

// BaseModule 基础模块结构
type BaseModule struct {
	DB     *gorm.DB
	Config ModuleConfig
}

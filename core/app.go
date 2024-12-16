package core

import (
	"github.com/andycai/unitool/lib/authentication"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	App             *fiber.App
	DB              *gorm.DB
	DBs             []*gorm.DB
	Config          *Config
	RouterPublic    fiber.Router
	RouterPublicApi fiber.Router
	RouterApi       fiber.Router
	RouterAdmin     fiber.Router
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start(dbs []*gorm.DB, fiberApp *fiber.App) {
	a.Config = &config
	a.DBs = dbs
	a.DB = dbs[0]
	a.App = fiberApp

	sqlDb, _ := a.DB.DB()
	authentication.SessionSetup(config.Database.Driver, sqlDb, config.Database.DSN, "sessions")

	// 注册静态路由
	serverConfig := a.Config.Server
	for _, staticPath := range serverConfig.StaticPaths {
		fiberApp.Static(staticPath.Route, staticPath.Path)
	}

	AwakeModules(a)

	// 初始化公共路由
	a.RouterPublic = fiberApp.Group("/")
	a.RouterPublicApi = fiberApp.Group("/api")
	InitPublicRouters()

	// 初始化API路由
	a.RouterApi = fiberApp.Group("/api")
	a.RouterApi.Use(middleware.AuthMiddleware)

	// 初始化管理员路由
	a.RouterAdmin = fiberApp.Group("/admin")
	a.RouterAdmin.Use(middleware.AuthMiddleware)
	InitAuthRouters()
}

func (a *App) HasPermission(permissionCode string) fiber.Handler {
	return middleware.HasPermission(permissionCode, a.CurrentUser)
}

// Current 获取当前用户
func (a *App) CurrentUser(c *fiber.Ctx) *models.User {
	isAuthenticated, userID := authentication.AuthGet(c)

	if !isAuthenticated {
		return nil
	}

	var vo models.User
	if err := a.DB.Preload("Role.Permissions").First(&vo, userID).Error; err != nil {
		return nil
	}

	return &vo
}

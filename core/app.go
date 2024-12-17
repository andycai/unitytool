package core

import (
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
	SessionSetup(config.Database.Driver, sqlDb, config.Database.DSN, "sessions")

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
	a.RouterApi.Use(AuthMiddleware)

	// 初始化管理员路由
	a.RouterAdmin = fiberApp.Group("/admin")
	a.RouterAdmin.Use(AuthMiddleware)
	InitAuthRouters()
}

func (a *App) HasPermission(permissionCode string) fiber.Handler {
	return HasPermission(permissionCode, a.CurrentUser)
}

// Current 获取当前用户
func (a *App) CurrentUser(c *fiber.Ctx) *models.User {
	isAuthenticated, userID := GetSession(c)

	if !isAuthenticated {
		return nil
	}

	var vo models.User
	if err := a.DB.Preload("Role.Permissions").First(&vo, userID).Error; err != nil {
		return nil
	}

	return &vo
}

// 如果模块已初始化，则跳过
func (a *App) IsInitializedModule(module string) bool {
	if err := a.DB.Model(&models.ModuleInit{}).Where("module = ?", module).First(&models.ModuleInit{}).Error; err != nil {
		return false
	}
	return true
}

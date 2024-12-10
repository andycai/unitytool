package login

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.login.gorm.db"
	KeyNoCheckRouter = "admin.login.router.nocheck"
	KeyCheckRouter   = "admin.login.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initPublicRouter(publicGroup fiber.Router) {
	// 登录页面路由（不需要认证）
	publicGroup.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{}, "login")
	})

	// 登录 API 路由（不需要认证）
	publicGroup.Post("/api/login", func(c *fiber.Ctx) error {
		return login(c)
	})
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	// 管理后台主页路由
	adminGroup.Get("/", func(c *fiber.Ctx) error {
		return c.Render("admin/index", fiber.Map{
			"Title": "管理后台",
		}, "admin/layout")
	})
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterPublicRouter(KeyNoCheckRouter, initPublicRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
}

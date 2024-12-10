package user

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.user.gorm.db"
	KeyNoCheckRouter = "admin.user.router.nocheck"
	KeyCheckRouter   = "admin.user.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	// 其他管理后台页面路由
	adminGroup.Get("/users", middleware.HasPermission("user:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/users", fiber.Map{
			"Title": "用户管理",
			"Scripts": []string{
				"/static/js/admin/users.js",
			},
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Get("/users", middleware.HasPermission("user:list"), GetUsers)
	apiGroup.Post("/users", middleware.HasPermission("user:create"), CreateUser)
	apiGroup.Put("/users/:id", middleware.HasPermission("user:update"), UpdateUser)
	apiGroup.Delete("/users/:id", middleware.HasPermission("user:delete"), DeleteUser)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

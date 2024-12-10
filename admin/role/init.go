package role

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.role.gorm.db"
	KeyNoCheckRouter = "admin.role.router.nocheck"
	KeyCheckRouter   = "admin.role.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	adminGroup.Get("/roles", middleware.HasPermission("role:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/roles", fiber.Map{
			"Title": "角色管理",
			"Scripts": []string{
				"/static/js/admin/roles.js",
			},
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Get("/roles", middleware.HasPermission("role:list"), getRoles)
	apiGroup.Post("/roles", middleware.HasPermission("role:create"), createRole)
	apiGroup.Put("/roles/:id", middleware.HasPermission("role:update"), updateRole)
	apiGroup.Delete("/roles/:id", middleware.HasPermission("role:delete"), deleteRole)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

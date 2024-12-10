package permission

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.permission.gorm.db"
	KeyNoCheckRouter = "admin.permission.router.nocheck"
	KeyCheckRouter   = "admin.permission.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	adminGroup.Get("/permissions", middleware.HasPermission("permission:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/permissions", fiber.Map{
			"Title": "权限管理",
			"Scripts": []string{
				"/static/js/admin/permissions.js",
			},
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Get("/permissions", middleware.HasPermission("permission:list"), getPermissions)
	apiGroup.Post("/permissions", middleware.HasPermission("permission:create"), createPermission)
	apiGroup.Put("/permissions/:id", middleware.HasPermission("permission:update"), updatePermission)
	apiGroup.Delete("/permissions/:id", middleware.HasPermission("permission:delete"), deletePermission)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

package permission

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type permissionModule struct {
}

func (m *permissionModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *permissionModule) InitDB() error {
	// 数据迁移
	return nil
}

func (m *permissionModule) InitModule() error {
	// public

	// admin
	app.RouterAdmin.Get("/permissions", middleware.HasPermission("permission:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/permissions", fiber.Map{
			"Title": "权限管理",
			"Scripts": []string{
				"/static/js/admin/permissions.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/permissions", middleware.HasPermission("permission:list"), getPermissions)
	app.RouterApi.Post("/permissions", middleware.HasPermission("permission:create"), createPermission)
	app.RouterApi.Put("/permissions/:id", middleware.HasPermission("permission:update"), updatePermission)
	app.RouterApi.Delete("/permissions/:id", middleware.HasPermission("permission:delete"), deletePermission)

	return nil
}

func init() {
	core.RegisterModule(&permissionModule{})
}

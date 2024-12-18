package permission

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type permissionModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&permissionModule{}, enum.ModulePriorityPermission)
}

func (m *permissionModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *permissionModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/permissions", app.HasPermission("permission:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/permissions", fiber.Map{
			"Title": "权限管理",
			"Scripts": []string{
				"/static/js/admin/permissions.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/permissions", app.HasPermission("permission:list"), getPermissions)
	app.RouterApi.Post("/permissions", app.HasPermission("permission:create"), createPermission)
	app.RouterApi.Put("/permissions/:id", app.HasPermission("permission:update"), updatePermission)
	app.RouterApi.Delete("/permissions/:id", app.HasPermission("permission:delete"), deletePermission)

	return nil
}

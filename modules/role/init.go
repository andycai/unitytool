package role

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type roleModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&roleModule{}, enum.ModulePriorityRole)
}

func (m *roleModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *roleModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/roles", app.HasPermission("role:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/roles", fiber.Map{
			"Title": "角色管理",
			"Scripts": []string{
				"/static/js/admin/roles.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/roles", app.HasPermission("role:list"), getRoles)
	app.RouterApi.Post("/roles", app.HasPermission("role:create"), createRole)
	app.RouterApi.Put("/roles/:id", app.HasPermission("role:update"), updateRole)
	app.RouterApi.Delete("/roles/:id", app.HasPermission("role:delete"), deleteRole)

	return nil
}

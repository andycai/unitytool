package role

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type roleModule struct {
}

func (m *roleModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *roleModule) InitDB() error {
	// 数据迁移
	return app.DB.AutoMigrate(&models.Role{}, &models.Permission{}, &models.RolePermission{})
}

func (m *roleModule) InitModule() error {
	// public

	// admin
	app.RouterAdmin.Get("/roles", middleware.HasPermission("role:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/roles", fiber.Map{
			"Title": "角色管理",
			"Scripts": []string{
				"/static/js/admin/roles.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/roles", middleware.HasPermission("role:list"), getRoles)
	app.RouterApi.Post("/roles", middleware.HasPermission("role:create"), createRole)
	app.RouterApi.Put("/roles/:id", middleware.HasPermission("role:update"), updateRole)
	app.RouterApi.Delete("/roles/:id", middleware.HasPermission("role:delete"), deleteRole)

	return nil
}

func init() {
	core.RegisterModule(&roleModule{})
}

package user

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type userModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&userModule{}, enum.ModulePriorityUser)
}

func (m *userModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	return initData()
}

func (m *userModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/users", app.HasPermission("user:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/users", fiber.Map{
			"Title": "用户管理",
			"Scripts": []string{
				"/static/js/admin/users.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/users", app.HasPermission("user:list"), getUsersAction)
	app.RouterApi.Post("/users", app.HasPermission("user:create"), createUserAction)
	app.RouterApi.Put("/users/:id", app.HasPermission("user:update"), updateUserAction)
	app.RouterApi.Delete("/users/:id", app.HasPermission("user:delete"), deleteUserAction)

	return nil
}

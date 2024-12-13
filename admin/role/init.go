package role

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type roleModule struct {
}

func (u *roleModule) Init(a *core.App) error {
	app = a
	return nil
}

func (u *roleModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *roleModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *roleModule) InitRouter() error {
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
	core.RegisterModules(&roleModule{})
}

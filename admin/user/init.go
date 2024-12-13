package user

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type userModule struct {
}

func (u *userModule) Init(a *core.App) error {
	app = a
	return nil
}

func (u *userModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *userModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *userModule) InitRouter() error {
	// public

	// admin
	app.RouterAdmin.Get("/users", middleware.HasPermission("user:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/users", fiber.Map{
			"Title": "用户管理",
			"Scripts": []string{
				"/static/js/admin/users.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/users", middleware.HasPermission("user:list"), getUsers)
	app.RouterApi.Post("/users", middleware.HasPermission("user:create"), createUser)
	app.RouterApi.Put("/users/:id", middleware.HasPermission("user:update"), updateUser)
	app.RouterApi.Delete("/users/:id", middleware.HasPermission("user:delete"), deleteUser)

	return nil
}

func init() {
	core.RegisterModules(&userModule{})
}

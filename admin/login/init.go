package login

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type loginModule struct {
}

func (m *loginModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *loginModule) InitDB() error {
	// 数据迁移
	return nil
}

func (m *loginModule) InitModule() error {
	// public
	// 登录页面路由（不需要认证）
	app.RouterPublic.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{}, "login")
	})

	// 登录 API 路由（不需要认证）
	app.RouterPublic.Post("/api/login", func(c *fiber.Ctx) error {
		return login(c)
	})

	// 修改密码路由（不需要认证）
	app.RouterPublic.Post("/api/change-password", func(c *fiber.Ctx) error {
		return changePassword(c)
	})

	// admin
	app.RouterAdmin.Get("/", func(c *fiber.Ctx) error {
		return c.Render("admin/index", fiber.Map{
			"Title": "管理后台",
		}, "admin/layout")
	})

	// api

	return nil
}

func init() {
	core.RegisterModule(&loginModule{})
}

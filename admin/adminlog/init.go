package adminlog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type adminlogModule struct {
}

func (u *adminlogModule) Init(a *core.App) error {
	app = a
	return nil
}

func (u *adminlogModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *adminlogModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *adminlogModule) InitRouter() error {
	// public

	// admin
	app.RouterAdmin.Get("/adminlog", middleware.HasPermission("adminlog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/adminlog", fiber.Map{
			"Title": "操作日志",
			"Scripts": []string{
				"/static/js/admin/adminlog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/adminlog", middleware.HasPermission("adminlog:list"), getAdminLogs)
	app.RouterApi.Delete("/adminlog", middleware.HasPermission("adminlog:delete"), deleteAdminLogs)

	return nil
}

func init() {
	core.RegisterModules(&adminlogModule{})
}

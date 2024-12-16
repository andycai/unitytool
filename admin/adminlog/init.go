package adminlog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type adminlogModule struct {
}

func (m *adminlogModule) Awake(a *core.App) error {
	app = a
	return app.DB.AutoMigrate(&models.AdminLog{})
}

func (m *adminlogModule) Start() error {
	return nil
}

func (m *adminlogModule) AddPublicRouters() error {
	return nil
}

func (m *adminlogModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/adminlog", app.HasPermission("adminlog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/adminlog", fiber.Map{
			"Title": "操作日志",
			"Scripts": []string{
				"/static/js/admin/adminlog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/adminlog", app.HasPermission("adminlog:list"), getAdminLogs)
	app.RouterApi.Delete("/adminlog", app.HasPermission("adminlog:delete"), deleteAdminLogs)

	return nil
}

func init() {
	core.RegisterModule(&adminlogModule{})
}

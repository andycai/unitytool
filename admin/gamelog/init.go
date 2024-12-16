package gamelog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type gamelogModule struct {
}

func (m *gamelogModule) Awake(a *core.App) error {
	app = a

	return app.DB.AutoMigrate(&models.GameLog{})
}

func (m *gamelogModule) Start() error {
	return nil
}

func (m *gamelogModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/gamelog", createLog)
	return nil
}

func (m *gamelogModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/gamelog", app.HasPermission("gamelog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/gamelog", fiber.Map{
			"Title": "游戏日志",
			"Scripts": []string{
				"/static/js/admin/gamelog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/gamelog", app.HasPermission("gamelog:list"), getLogs)
	app.RouterApi.Delete("/gamelog/before", app.HasPermission("gamelog:delete"), deleteLogsBefore)
	app.RouterApi.Delete("/gamelog/:id", app.HasPermission("gamelog:list"), deleteLog)

	return nil
}

func init() {
	core.RegisterModule(&gamelogModule{})
}

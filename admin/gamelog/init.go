package gamelog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type gamelogModule struct {
}

func (u *gamelogModule) Init(a *core.App) error {
	app = a
	return nil
}

func (u *gamelogModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *gamelogModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *gamelogModule) InitRouter() error {
	// public
	app.RouterPublic.Post("/api/gamelog", createLog)

	// admin
	app.RouterAdmin.Get("/gamelog", middleware.HasPermission("gamelog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/gamelog", fiber.Map{
			"Title": "游戏日志",
			"Scripts": []string{
				"/static/js/admin/gamelog.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/gamelog", middleware.HasPermission("gamelog:list"), getLogs)
	app.RouterApi.Delete("/gamelog/before", middleware.HasPermission("gamelog:delete"), deleteLogsBefore)
	app.RouterApi.Delete("/gamelog/:id", middleware.HasPermission("gamelog:list"), deleteLog)

	return nil
}

func init() {
	core.RegisterModules(&gamelogModule{})
}

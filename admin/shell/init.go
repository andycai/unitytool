package shell

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type shellModule struct {
}

func (u *shellModule) Init(a *core.App) error {
	app = a
	return nil
}

func (u *shellModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *shellModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *shellModule) InitRouter() error {
	// public
	app.RouterPublic.Post("/api/shell", func(c *fiber.Ctx) error {
		return execShell(c, app.Config.Server.ScriptPath)
	})

	// admin

	// api

	return nil
}

func init() {
	core.RegisterModules(&shellModule{})
}

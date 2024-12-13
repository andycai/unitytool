package shell

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type shellModule struct {
}

func (m *shellModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *shellModule) InitDB() error {
	// 数据迁移
	return nil
}

func (m *shellModule) InitModule() error {
	// public
	app.RouterPublic.Post("/api/shell", func(c *fiber.Ctx) error {
		return execShell(c, app.Config.Server.ScriptPath)
	})

	// admin

	// api

	return nil
}

func init() {
	core.RegisterModule(&shellModule{})
}

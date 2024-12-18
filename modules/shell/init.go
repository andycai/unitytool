package shell

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type shellModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&shellModule{}, enum.ModulePriorityShell)
}

func (m *shellModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *shellModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/shell", func(c *fiber.Ctx) error {
		return execShell(c, app.Config.Server.ScriptPath)
	})

	return nil
}

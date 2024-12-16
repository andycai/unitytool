package shell

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type shellModule struct {
}

func (m *shellModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *shellModule) Start() error {
	return nil
}

func (m *shellModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/shell", func(c *fiber.Ctx) error {
		return execShell(c, app.Config.Server.ScriptPath)
	})

	return nil
}

func (m *shellModule) AddAuthRouters() error {
	return nil
}

func init() {
	core.RegisterModule(&shellModule{})
}

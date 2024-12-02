package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
	"mind.com/log/utils"
)

type CmdModule struct {
	BaseModule
	ServerConfig utils.ServerConfig
}

func (m *CmdModule) Init() error {
	return nil
}

func (m *CmdModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 脚本命令执行路由
	app.Post("/cmd", func(c *fiber.Ctx) error {
		return handlers.ExecShell(c, m.ServerConfig.ScriptPath)
	})
}

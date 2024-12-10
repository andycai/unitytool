package shell

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	KeyModule        = "admin.shell"
	KeyNoCheckRouter = "admin.shell.router.nocheck"
)

func initModule() {
}

func initPublicRouter(publicGroup fiber.Router) {
	publicGroup.Post("/api/shell", func(c *fiber.Ctx) error {
		return ExecShell(c, utils.GetServerConfig().ScriptPath)
	})
}

func init() {
	core.RegisterModule(KeyModule, initModule)
	core.RegisterPublicRouter(KeyNoCheckRouter, initPublicRouter)
}

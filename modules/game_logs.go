package modules

import (
	"github.com/andycai/unitool/handlers"
	"github.com/gofiber/fiber/v2"
)

type GameLogsModule struct {
	BaseModule
}

func (m *GameLogsModule) Init() error {
	return nil
}

func (m *GameLogsModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	adminGroup.Get("/game/logs", func(c *fiber.Ctx) error {
		return c.Render("admin/game_logs", fiber.Map{
			"Title": "游戏日志",
			"Scripts": []string{
				"/static/js/admin/game_logs.js",
			},
		}, "admin/layout")
	})

	apiGroup.Post("/game/logs", func(c *fiber.Ctx) error {
		return handlers.CreateLog(c, m.DB)
	})

	apiGroup.Get("/game/logs", func(c *fiber.Ctx) error {
		return handlers.GetLogs(c, m.DB)
	})

	apiGroup.Delete("/game/logs/before", func(c *fiber.Ctx) error {
		return handlers.DeleteLogsBefore(c, m.DB)
	})

	apiGroup.Delete("/game/logs/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteLog(c, m.DB)
	})
}

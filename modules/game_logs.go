package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
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

	apiGroup.Post("/game_logs", func(c *fiber.Ctx) error {
		return handlers.CreateLog(c, m.DB)
	})

	apiGroup.Get("/game_logs", func(c *fiber.Ctx) error {
		return handlers.GetLogs(c, m.DB)
	})

	apiGroup.Delete("/game_logs/before", func(c *fiber.Ctx) error {
		return handlers.DeleteLogsBefore(c, m.DB)
	})

	apiGroup.Delete("/game_logs/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteLog(c, m.DB)
	})
}

package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
)

type LogsModule struct {
	BaseModule
}

func (m *LogsModule) Init() error {
	return nil
}

func (m *LogsModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	apiGroup.Post("/logs", func(c *fiber.Ctx) error {
		return handlers.CreateLog(c, m.DB)
	})

	apiGroup.Get("/logs", func(c *fiber.Ctx) error {
		return handlers.GetLogs(c, m.DB)
	})

	apiGroup.Delete("/logs/before", func(c *fiber.Ctx) error {
		return handlers.DeleteLogsBefore(c, m.DB)
	})

	apiGroup.Delete("/logs/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteLog(c, m.DB)
	})
}

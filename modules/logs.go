package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
)

type LogsModule struct {
	BaseModule
}

func (m *LogsModule) Init() error {
	// 初始化日志模块
	return nil
}

func (m *LogsModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	app.Post("/api/logs", func(c *fiber.Ctx) error {
		return handlers.CreateLog(c, m.DB)
	})

	app.Get("/api/logs", func(c *fiber.Ctx) error {
		return handlers.GetLogs(c, m.DB)
	})

	app.Delete("/api/logs/before", func(c *fiber.Ctx) error {
		return handlers.DeleteLogsBefore(c, m.DB)
	})

	app.Delete("/api/logs/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteLog(c, m.DB)
	})
}

package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
)

type StatsModule struct {
	BaseModule
}

func (m *StatsModule) Init() error {
	// 初始化统计模块
	return nil
}

func (m *StatsModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 创建统计记录
	app.Post("/api/stats", func(c *fiber.Ctx) error {
		return handlers.CreateStats(c, m.DB)
	})

	// 获取统计记录列表
	app.Get("/api/stats", func(c *fiber.Ctx) error {
		return handlers.GetStats(c, m.DB)
	})

	// 删除指定日期前的统计记录
	app.Delete("/api/stats/before", func(c *fiber.Ctx) error {
		return handlers.DeleteStatsBefore(c, m.DB)
	})

	// 获取统计详情
	app.Get("/api/stats/details", func(c *fiber.Ctx) error {
		return handlers.GetStatDetails(c, m.DB)
	})

	// 删除单条统计记录
	app.Delete("/api/stats/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteStat(c, m.DB)
	})
}

package modules

import (
	"github.com/andycai/unitool/handlers"
	"github.com/gofiber/fiber/v2"
)

type StatsModule struct {
	BaseModule
}

func (m *StatsModule) Init() error {
	return nil
}

func (m *StatsModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	adminGroup.Get("/stats", func(c *fiber.Ctx) error {
		return c.Render("admin/stats", fiber.Map{
			"Title": "游戏统计",
			"Scripts": []string{
				"/static/js/chart-4.4.4.js",
				"/static/js/hammer-2.0.8.js",
				"/static/js/chartjs-plugin-zoom.min.js",
				"/static/js/chartjs-adapter-date-fns.bundle.min.js",
				"/static/js/admin/stats.js",
			},
		}, "admin/layout")
	})

	// 创建统计记录
	apiGroup.Post("/stats", func(c *fiber.Ctx) error {
		return handlers.CreateStats(c, m.DB)
	})

	// 获取统计记录列表
	apiGroup.Get("/stats", func(c *fiber.Ctx) error {
		return handlers.GetStats(c, m.DB)
	})

	// 删除指定日期前的统计记录
	apiGroup.Delete("/stats/before", func(c *fiber.Ctx) error {
		return handlers.DeleteStatsBefore(c, m.DB)
	})

	// 获取统计详情
	apiGroup.Get("/stats/details", func(c *fiber.Ctx) error {
		return handlers.GetStatDetails(c, m.DB)
	})

	// 删除单条统计记录
	apiGroup.Delete("/stats/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteStat(c, m.DB)
	})
}

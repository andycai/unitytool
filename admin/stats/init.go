package stats

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type statsModule struct {
}

func (m *statsModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *statsModule) InitDB() error {
	// 数据迁移
	return app.DB.AutoMigrate(&models.StatsRecord{}, &models.StatsInfo{})
}

func (m *statsModule) InitModule() error {
	// public
	app.RouterPublic.Post("/api/stats", CreateStats)

	// admin
	app.RouterAdmin.Get("/stats", middleware.HasPermission("stats:list"), func(c *fiber.Ctx) error {
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

	// api
	app.RouterApi.Get("/stats", middleware.HasPermission("stats:list"), getStats)
	app.RouterApi.Delete("/stats/before", middleware.HasPermission("stats:delete"), deleteStatsBefore)
	app.RouterApi.Get("/stats/details", middleware.HasPermission("stats:list"), getStatDetails)
	app.RouterApi.Delete("/stats/:id", middleware.HasPermission("stats:delete"), deleteStat)

	return nil
}

func init() {
	core.RegisterModule(&statsModule{})
}

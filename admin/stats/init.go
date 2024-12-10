package stats

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.stats.gorm.db"
	KeyNoCheckRouter = "admin.stats.router.nocheck"
	KeyCheckRouter   = "admin.stats.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initPublicRouter(publicGroup fiber.Router) {
	publicGroup.Post("/api/stats", CreateStats)
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	adminGroup.Get("/stats", middleware.HasPermission("stats:list"), func(c *fiber.Ctx) error {
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
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Get("/stats", middleware.HasPermission("stats:list"), getStats)
	apiGroup.Delete("/stats/before", middleware.HasPermission("stats:delete"), deleteStatsBefore)
	apiGroup.Get("/stats/details", middleware.HasPermission("stats:list"), getStatDetails)
	apiGroup.Delete("/stats/:id", middleware.HasPermission("stats:delete"), deleteStat)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterPublicRouter(KeyNoCheckRouter, initPublicRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

package gamelog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.gamelog.gorm.db"
	KeyNoCheckRouter = "admin.gamelog.router.nocheck"
	KeyCheckRouter   = "admin.gamelog.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initPublicRouter(publicGroup fiber.Router) {
	publicGroup.Post("/api/gamelog", createLog)
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	adminGroup.Get("/gamelog", middleware.HasPermission("gamelog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/gamelog", fiber.Map{
			"Title": "游戏日志",
			"Scripts": []string{
				"/static/js/admin/gamelog.js",
			},
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Get("/gamelog", middleware.HasPermission("gamelog:list"), getLogs)
	apiGroup.Delete("/gamelog/before", middleware.HasPermission("gamelog:delete"), deleteLogsBefore)
	apiGroup.Delete("/gamelog/:id", middleware.HasPermission("gamelog:list"), deleteLog)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterPublicRouter(KeyNoCheckRouter, initPublicRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

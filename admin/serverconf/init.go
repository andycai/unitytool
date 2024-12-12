package serverconf

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	KeyModule        = "admin.serverconf"
	KeyDB            = "admin.serverconf.gorm.db"
	KeyNoCheckRouter = "admin.serverconf.router.nocheck"
	KeyCheckRouter   = "admin.serverconf.router.check"
)

// var db *gorm.DB

// func initDB(dbs []*gorm.DB) {
// 	db = dbs[0]
// }

func initModule() {
	InitJSONPaths(utils.GetJSONPathConfig())
}

func initPublicRouter(publicGroup fiber.Router) {
	publicGroup.Get("/api/game/serverlist", getServerList)
	publicGroup.Get("/api/game/lastserver", getLastServer)
	publicGroup.Get("/api/game/serverinfo", getServerInfo)
	publicGroup.Get("/api/game/noticelist", getNoticeList)
	publicGroup.Get("/api/game/noticenum", getNoticeNum)
	publicGroup.Get("/api/serverlist", getServerList)
	publicGroup.Get("/api/lastserver", getLastServer)
	publicGroup.Get("/api/serverinfo", getServerInfo)
	publicGroup.Get("/api/noticelist", getNoticeList)
	publicGroup.Get("/api/noticenum", getNoticeNum)
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	adminGroup.Get("/serverconf", middleware.HasPermission("serverconf:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/serverconf", fiber.Map{
			"Title": "服务器配置",
			"Scripts": []string{
				"/static/js/admin/serverconf.js",
			},
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Post("/game/serverlist", middleware.HasPermission("serverconf:update"), updateServerList)
	apiGroup.Post("/game/lastserver", middleware.HasPermission("serverconf:update"), updateLastServer)
	apiGroup.Post("/game/serverinfo", middleware.HasPermission("serverconf:update"), updateServerInfo)
	apiGroup.Post("/game/noticelist", middleware.HasPermission("serverconf:update"), updateNoticeList)
	apiGroup.Post("/game/noticenum", middleware.HasPermission("serverconf:update"), updateNoticeNum)
}

func init() {
	core.RegisterModule(KeyModule, initModule)
	// core.RegisterDatabase(KeyDB, initDB)
	core.RegisterPublicRouter(KeyNoCheckRouter, initPublicRouter)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

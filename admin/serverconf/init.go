package serverconf

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type serverconfModule struct {
}

func (m *serverconfModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *serverconfModule) InitDB() error {
	// 数据迁移
	return nil
}

func (m *serverconfModule) InitModule() error {
	// public
	app.RouterPublic.Get("/api/game/serverlist", getServerList)
	app.RouterPublic.Get("/api/game/lastserver", getLastServer)
	app.RouterPublic.Get("/api/game/serverinfo", getServerInfo)
	app.RouterPublic.Get("/api/game/noticelist", getNoticeList)
	app.RouterPublic.Get("/api/game/noticenum", getNoticeNum)
	app.RouterPublic.Get("/api/serverlist", getServerList)
	app.RouterPublic.Get("/api/lastserver", getLastServer)
	app.RouterPublic.Get("/api/serverinfo", getServerInfo)
	app.RouterPublic.Get("/api/noticelist", getNoticeList)
	app.RouterPublic.Get("/api/noticenum", getNoticeNum)

	// admin
	app.RouterAdmin.Get("/serverconf", app.HasPermission("serverconf:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/serverconf", fiber.Map{
			"Title": "服务器配置",
			"Scripts": []string{
				"/static/js/admin/serverconf.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Post("/game/serverlist", app.HasPermission("serverconf:update"), updateServerList)
	app.RouterApi.Post("/game/lastserver", app.HasPermission("serverconf:update"), updateLastServer)
	app.RouterApi.Post("/game/serverinfo", app.HasPermission("serverconf:update"), updateServerInfo)
	app.RouterApi.Post("/game/noticelist", app.HasPermission("serverconf:update"), updateNoticeList)
	app.RouterApi.Post("/game/noticenum", app.HasPermission("serverconf:update"), updateNoticeNum)

	return nil
}

func init() {
	core.RegisterModule(&serverconfModule{})
}

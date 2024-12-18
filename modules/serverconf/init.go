package serverconf

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type serverconfModule struct {
	core.BaseModule
}

func init() {
	core.RegisterModule(&serverconfModule{}, enum.ModulePriorityServerconf)
}

func (m *serverconfModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *serverconfModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Get("/game/serverlist", getServerList)
	app.RouterPublicApi.Get("/game/lastserver", getLastServer)
	app.RouterPublicApi.Get("/game/serverinfo", getServerInfo)
	app.RouterPublicApi.Get("/game/noticelist", getNoticeList)
	app.RouterPublicApi.Get("/game/noticenum", getNoticeNum)
	app.RouterPublicApi.Get("/serverlist", getServerList)
	app.RouterPublicApi.Get("/lastserver", getLastServer)
	app.RouterPublicApi.Get("/serverinfo", getServerInfo)
	app.RouterPublicApi.Get("/noticelist", getNoticeList)
	app.RouterPublicApi.Get("/noticenum", getNoticeNum)

	return nil
}

func (m *serverconfModule) AddAuthRouters() error {
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

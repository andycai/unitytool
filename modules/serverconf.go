package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
	"mind.com/log/utils"
)

type ServerConfModule struct {
	BaseModule
	ServerConfig utils.ServerConfig
	JSONPaths    utils.JSONPathConfig
}

func (m *ServerConfModule) Init() error {
	// 初始化 JSON 路径配置
	handlers.InitJSONPaths(m.JSONPaths)
	return nil
}

func (m *ServerConfModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 配置页面路由
	app.Get("/serverconf", func(c *fiber.Ctx) error {
		return c.SendFile("templates/serverconf.html")
	})

	// API 路由
	app.Get("/api/serverlist", handlers.GetServerList)
	app.Post("/api/serverlist", handlers.UpdateServerList)

	app.Get("/api/lastserver", handlers.GetLastServer)
	app.Post("/api/lastserver", handlers.UpdateLastServer)

	app.Get("/api/serverinfo", handlers.GetServerInfo)
	app.Post("/api/serverinfo", handlers.UpdateServerInfo)

	// 添加公告列表路由
	app.Get("/api/noticelist", handlers.GetNoticeList)
	app.Post("/api/noticelist", handlers.UpdateNoticeList)

	// 添加公告数量路由
	app.Get("/api/noticenum", handlers.GetNoticeNum)
	app.Post("/api/noticenum", handlers.UpdateNoticeNum)
}

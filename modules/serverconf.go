package modules

import (
	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/utils"
	"github.com/gofiber/fiber/v2"
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
	adminGroup.Get("/serverconf", func(c *fiber.Ctx) error {
		return c.Render("admin/serverconf", fiber.Map{
			"Title": "服务器配置",
			"Scripts": []string{
				"/static/js/admin/serverconf.js",
			},
		}, "admin/layout")
	})

	// 私有 API 路由
	apiGroup.Post("/game/serverlist", handlers.UpdateServerList)
	apiGroup.Post("/game/lastserver", handlers.UpdateLastServer)
	apiGroup.Post("/game/serverinfo", handlers.UpdateServerInfo)
	apiGroup.Post("/game/noticelist", handlers.UpdateNoticeList)
	apiGroup.Post("/game/noticenum", handlers.UpdateNoticeNum)

	//  公开 API 路由
	openGroup.Get("/game/serverlist", handlers.GetServerList)
	openGroup.Get("/game/lastserver", handlers.GetLastServer)
	openGroup.Get("/game/serverinfo", handlers.GetServerInfo)
	openGroup.Get("/game/noticelist", handlers.GetNoticeList)
	openGroup.Get("/game/noticenum", handlers.GetNoticeNum)
}

package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/utils"
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
		return c.SendFile("templates/serverconf.html")
	})

	// API 路由
	apiGroup.Get("/serverlist", handlers.GetServerList)
	apiGroup.Post("/serverlist", handlers.UpdateServerList)

	apiGroup.Get("/lastserver", handlers.GetLastServer)
	apiGroup.Post("/lastserver", handlers.UpdateLastServer)

	apiGroup.Get("/serverinfo", handlers.GetServerInfo)
	apiGroup.Post("/serverinfo", handlers.UpdateServerInfo)

	// 添加公告列表路由
	apiGroup.Get("/noticelist", handlers.GetNoticeList)
	apiGroup.Post("/noticelist", handlers.UpdateNoticeList)

	// 添加公告数量路由
	apiGroup.Get("/noticenum", handlers.GetNoticeNum)
	apiGroup.Post("/noticenum", handlers.UpdateNoticeNum)
}

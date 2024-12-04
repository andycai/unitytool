package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
	"mind.com/log/utils"
)

type BrowseModule struct {
	BaseModule
	ServerConfig utils.ServerConfig
}

func (m *BrowseModule) Init() error {
	// 初始化文件浏览模块
	return nil
}

func (m *BrowseModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 文件浏览路由
	app.Get("/browse/*", func(c *fiber.Ctx) error {
		return handlers.HandleFileServer(c, m.ServerConfig.Output)
	})

	// 文件删除路由
	app.Delete("/browse/*", func(c *fiber.Ctx) error {
		return handlers.HandleFileServer(c, m.ServerConfig.Output)
	})
}

package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/utils"
)

type BrowseModule struct {
	BaseModule
	ServerConfig utils.ServerConfig
}

func (m *BrowseModule) Init() error {
	return nil
}

func (m *BrowseModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 文件浏览路由
	adminGroup.Get("/browse/*", func(c *fiber.Ctx) error {
		return handlers.HandleFileServer(c, m.ServerConfig.Output)
	})

	// 文件删除路由
	apiGroup.Delete("/browse/*", func(c *fiber.Ctx) error {
		return handlers.HandleFileServer(c, m.ServerConfig.Output)
	})
}

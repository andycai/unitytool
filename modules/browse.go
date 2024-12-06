package modules

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/utils"
	"github.com/gofiber/fiber/v2"
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

	// 配置页面路由
	// adminGroup.Get("/browse", func(c *fiber.Ctx) error {
	// 	return c.Render("admin/directory", fiber.Map{
	// 		"Path": "/",
	// 	})
	// })

	adminGroup.Get("/browse/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		if path == "" {
			path = "."
		}

		// URL 解码路径
		decodedPath, err := url.QueryUnescape(path)
		if err != nil {
			return c.Status(400).SendString("Invalid path encoding")
		}

		fullPath := filepath.Join(m.ServerConfig.Output, decodedPath)

		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("File not found: %s", fullPath))
		}

		// 如果是目录，显示目录内容
		if fileInfo.IsDir() {
			return handlers.HandleBrowseDirectory(c, fullPath)
		}

		// 如果是文件，显示文件内容
		return handlers.HandleBrowseFile(c, fullPath)
	})

	// 文件删除路由
	apiGroup.Delete("/browse/*", func(c *fiber.Ctx) error {
		return handlers.HandleBrowseDelete(c, m.ServerConfig.Output)
	})
}

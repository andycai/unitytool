package modules

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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

	// 浏览目录和文件的路由
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

		// 获取配置的根目录的绝对路径
		rootDir := m.ServerConfig.Output
		absRootDir, err := filepath.Abs(rootDir)
		if err != nil {
			return c.Status(500).SendString("Invalid root directory configuration")
		}

		// 构建完整路径
		fullPath := filepath.Join(rootDir, decodedPath)

		// 获取完整路径的绝对路径
		absPath, err := filepath.Abs(fullPath)
		if err != nil {
			return c.Status(400).SendString("Invalid path")
		}

		// 确保访问路径在根目录内
		if !strings.HasPrefix(absPath, absRootDir) {
			return c.Status(403).SendString("Access denied: Path outside root directory")
		}

		fileInfo, err := os.Stat(absPath)
		if err != nil {
			return c.Status(404).SendString(fmt.Sprintf("File not found: %s", decodedPath))
		}

		// 如果是目录，显示目录内容
		if fileInfo.IsDir() {
			return handlers.HandleBrowseDirectory(c, absPath)
		}

		// 如果是文件，显示文件内容
		return handlers.HandleBrowseFile(c, absPath)
	})

	// 文件删除路由
	adminGroup.Delete("/browse/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		if path == "" {
			return c.Status(400).SendString("Path is required")
		}

		// URL 解码路径
		decodedPath, err := url.QueryUnescape(path)
		if err != nil {
			return c.Status(400).SendString("Invalid path encoding")
		}

		// 获取配置的根目录的绝对路径
		rootDir := m.ServerConfig.Output
		absRootDir, err := filepath.Abs(rootDir)
		if err != nil {
			return c.Status(500).SendString("Invalid root directory configuration")
		}

		// 构建完整路径
		fullPath := filepath.Join(rootDir, decodedPath)

		// 获取完整路径的绝对路径
		absPath, err := filepath.Abs(fullPath)
		if err != nil {
			return c.Status(400).SendString("Invalid path")
		}

		// 确保删除路径在根目录内
		if !strings.HasPrefix(absPath, absRootDir) {
			return c.Status(403).SendString("Access denied: Path outside root directory")
		}

		// 检查是否是目录
		fileInfo, err := os.Stat(absPath)
		if err != nil {
			return c.Status(404).SendString("File not found")
		}
		if fileInfo.IsDir() {
			return c.Status(400).SendString("Cannot delete directories")
		}

		return handlers.HandleBrowseDelete(c, absPath)
	})

	// FTP 上传路由
	adminGroup.Post("/ftp/upload", func(c *fiber.Ctx) error {
		return handlers.HandleFTPUpload(c, m.ServerConfig.Output)
	})
}

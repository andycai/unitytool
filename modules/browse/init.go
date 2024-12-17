package browse

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type browseModule struct {
}

func init() {
	core.RegisterModule(&browseModule{}, enum.ModulePriorityBrowse)
}

func (m *browseModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *browseModule) Start() error {
	return nil
}

func (m *browseModule) AddPublicRouters() error {
	return nil
}

func (m *browseModule) AddAuthRouters() error {
	// admin
	// 浏览目录和文件的路由
	app.RouterAdmin.Get("/browse/*", app.HasPermission("browse:list"), func(c *fiber.Ctx) error {
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
		rootDir := app.Config.Server.Output
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
			return fiber.NewError(fiber.StatusForbidden, "Access denied: Path outside root directory")
		}

		fileInfo, err := os.Stat(absPath)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("File not found: %s", decodedPath))
		}

		// 如果是目录，显示目录内容
		if fileInfo.IsDir() {
			return handleBrowseDirectory(c, absPath)
		}

		// 如果是文件，显示文件内容
		return handleBrowseFile(c, absPath)
	})

	// 文件删除路由
	app.RouterAdmin.Delete("/browse/*", app.HasPermission("browse:delete"), func(c *fiber.Ctx) error {
		path := c.Params("*")
		if path == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Path is required")
		}

		// URL 解码路径
		decodedPath, err := url.QueryUnescape(path)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid path encoding")
		}

		// 获取配置的根目录的绝对路径
		rootDir := app.Config.Server.Output
		absRootDir, err := filepath.Abs(rootDir)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Invalid root directory configuration")
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
			return fiber.NewError(fiber.StatusForbidden, "Access denied: Path outside root directory")
		}

		// 检查是否是目录
		fileInfo, err := os.Stat(absPath)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}
		if fileInfo.IsDir() {
			return c.Status(fiber.StatusBadRequest).SendString("Cannot delete directories")
		}

		return handleBrowseDelete(c, absPath)
	})

	// FTP 上传路由
	app.RouterAdmin.Post("/ftp/upload", app.HasPermission("browse:ftp"), func(c *fiber.Ctx) error {
		return uploadByFTP(c, app.Config.Server.Output)
	})

	// api

	return nil
}

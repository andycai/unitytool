package modules

import (
	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
)

type AdminLogsModule struct {
	BaseModule
}

func (m *AdminLogsModule) Init() error {
	return nil
}

func (m *AdminLogsModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 管理员日志页面
	adminGroup.Get("/logs", middleware.HasPermission("admin_log:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/logs", fiber.Map{
			"Title": "操作日志",
			"Scripts": []string{
				"/static/js/admin/logs.js",
			},
		}, "admin/layout")
	})

	// 获取日志列表
	apiGroup.Get("/logs", middleware.HasPermission("admin_log:list"), func(c *fiber.Ctx) error {
		return handlers.GetAdminLogs(c, m.DB)
	})

	// 删除日志
	apiGroup.Delete("/logs", middleware.HasPermission("admin_log:delete"), func(c *fiber.Ctx) error {
		return handlers.DeleteAdminLogs(c, m.DB)
	})
}

package modules

import (
	"github.com/andycai/unitool/handlers"
	"github.com/gofiber/fiber/v2"
)

type CITaskModule struct {
	BaseModule
}

func (m *CITaskModule) Init() error {
	return nil
}

func (m *CITaskModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 页面路由
	adminGroup.Get("/citask", func(c *fiber.Ctx) error {
		return c.Render("admin/citask", fiber.Map{
			"Title": "任务管理",
			"Scripts": []string{
				"/static/js/admin/citask.js",
			},
		}, "admin/layout")
	})

	// 任务管理 API
	apiGroup.Get("/citask", func(c *fiber.Ctx) error { return handlers.GetTasks(c, m.DB) })                 // 获取任务列表
	apiGroup.Post("/citask", func(c *fiber.Ctx) error { return handlers.CreateTask(c, m.DB) })              // 创建任务
	apiGroup.Get("/citask/:id", func(c *fiber.Ctx) error { return handlers.GetTask(c, m.DB) })              // 获取任务详情
	apiGroup.Put("/citask/:id", func(c *fiber.Ctx) error { return handlers.UpdateTask(c, m.DB) })           // 更新任务
	apiGroup.Delete("/citask/:id", func(c *fiber.Ctx) error { return handlers.DeleteTask(c, m.DB) })        // 删除任务
	apiGroup.Post("/citask/:id/run", func(c *fiber.Ctx) error { return handlers.RunTask(c, m.DB) })         // 执行任务
	apiGroup.Get("/citask/:id/logs", func(c *fiber.Ctx) error { return handlers.GetTaskLogs(c, m.DB) })     // 获取任务日志
	apiGroup.Get("/citask/:id/status", func(c *fiber.Ctx) error { return handlers.GetTaskStatus(c, m.DB) }) // 获取任务状态
}

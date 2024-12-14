package citask

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type taskModule struct {
}

func (m *taskModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *taskModule) InitDB() error {
	// 数据迁移
	return app.DB.AutoMigrate(&models.Task{}, &models.TaskLog{})
}

func (m *taskModule) InitModule() error {
	// logic
	initCron()

	// public

	// admin
	app.RouterAdmin.Get("/citask", app.HasPermission("citask:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/citask", fiber.Map{
			"Title": "任务管理",
			"Scripts": []string{
				"/static/js/admin/citask.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/citask", app.HasPermission("citask:list"), getTasks)                        // 获取任务列表
	app.RouterApi.Post("/citask", app.HasPermission("citask:create"), createTask)                   // 创建任务
	app.RouterApi.Get("/citask/running", app.HasPermission("citask:list"), GetRunningTasks)         // 获取正在执行的任务
	app.RouterApi.Get("/citask/next-run", app.HasPermission("citask:list"), getNextRunTime)         // 计算下次执行时间
	app.RouterApi.Get("/citask/search", app.HasPermission("citask:list"), searchTasks)              // 添加搜索接口
	app.RouterApi.Get("/citask/:id", app.HasPermission("citask:list"), getTask)                     // 获取任务详情
	app.RouterApi.Put("/citask/:id", app.HasPermission("citask:update"), updateTask)                // 更新任务
	app.RouterApi.Delete("/citask/:id", app.HasPermission("citask:delete"), deleteTask)             // 删除任务
	app.RouterApi.Post("/citask/run/:id", app.HasPermission("citask:run"), runTask)                 // 执行任务
	app.RouterApi.Get("/citask/logs/:id", app.HasPermission("citask:list"), getTaskLogs)            // 获取任务日志
	app.RouterApi.Get("/citask/progress/:logId", app.HasPermission("citask:list"), getTaskProgress) // 获取任务进度
	app.RouterApi.Post("/citask/stop/:logId", app.HasPermission("citask:run"), stopTask)            // 停止任务

	return nil
}

func init() {
	core.RegisterModule(&taskModule{})
}

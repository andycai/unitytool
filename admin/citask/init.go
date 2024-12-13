package citask

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
)

var app *core.App

type taskModule struct {
}

func (u *taskModule) Init(a *core.App) error {
	app = a

	// 初始化定时任务
	initCron()

	return nil
}

func (u *taskModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *taskModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *taskModule) InitRouter() error {
	// public

	// admin
	app.RouterAdmin.Get("/citask", middleware.HasPermission("citask:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/citask", fiber.Map{
			"Title": "任务管理",
			"Scripts": []string{
				"/static/js/admin/citask.js",
			},
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/citask", middleware.HasPermission("citask:list"), getTasks)                        // 获取任务列表
	app.RouterApi.Post("/citask", middleware.HasPermission("citask:create"), createTask)                   // 创建任务
	app.RouterApi.Get("/citask/running", middleware.HasPermission("citask:list"), GetRunningTasks)         // 获取正在执行的任务
	app.RouterApi.Get("/citask/next-run", middleware.HasPermission("citask:list"), getNextRunTime)         // 计算下次执行时间
	app.RouterApi.Get("/citask/search", middleware.HasPermission("citask:list"), searchTasks)              // 添加搜索接口
	app.RouterApi.Get("/citask/:id", middleware.HasPermission("citask:list"), getTask)                     // 获取任务详情
	app.RouterApi.Put("/citask/:id", middleware.HasPermission("citask:update"), updateTask)                // 更新任务
	app.RouterApi.Delete("/citask/:id", middleware.HasPermission("citask:delete"), deleteTask)             // 删除任务
	app.RouterApi.Post("/citask/run/:id", middleware.HasPermission("citask:run"), runTask)                 // 执行任务
	app.RouterApi.Get("/citask/logs/:id", middleware.HasPermission("citask:list"), getTaskLogs)            // 获取任务日志
	app.RouterApi.Get("/citask/progress/:logId", middleware.HasPermission("citask:list"), getTaskProgress) // 获取任务进度
	app.RouterApi.Post("/citask/stop/:logId", middleware.HasPermission("citask:run"), stopTask)            // 停止任务

	return nil
}

func init() {
	core.RegisterModules(&taskModule{})
}

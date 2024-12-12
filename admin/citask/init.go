package citask

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.citask.gorm.db"
	KeyNoCheckRouter = "admin.citask.router.nocheck"
	KeyCheckRouter   = "admin.citask.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]

	initCron()
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	adminGroup.Get("/citask", middleware.HasPermission("citask:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/citask", fiber.Map{
			"Title": "任务管理",
			"Scripts": []string{
				"/static/js/admin/citask.js",
			},
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Get("/citask", middleware.HasPermission("citask:list"), getTasks)                        // 获取任务列表
	apiGroup.Post("/citask", middleware.HasPermission("citask:create"), createTask)                   // 创建任务
	apiGroup.Get("/citask/running", middleware.HasPermission("citask:list"), GetRunningTasks)         // 获取正在执行的任务
	apiGroup.Get("/citask/next-run", middleware.HasPermission("citask:list"), getNextRunTime)         // 计算下次执行时间
	apiGroup.Get("/citask/search", middleware.HasPermission("citask:list"), searchTasks)              // 添加搜索接口
	apiGroup.Get("/citask/:id", middleware.HasPermission("citask:list"), getTask)                     // 获取任务详情
	apiGroup.Put("/citask/:id", middleware.HasPermission("citask:update"), updateTask)                // 更新任务
	apiGroup.Delete("/citask/:id", middleware.HasPermission("citask:delete"), deleteTask)             // 删除任务
	apiGroup.Post("/citask/run/:id", middleware.HasPermission("citask:run"), runTask)                 // 执行任务
	apiGroup.Get("/citask/logs/:id", middleware.HasPermission("citask:list"), getTaskLogs)            // 获取任务日志
	apiGroup.Get("/citask/progress/:logId", middleware.HasPermission("citask:list"), getTaskProgress) // 获取任务进度
	apiGroup.Post("/citask/stop/:logId", middleware.HasPermission("citask:run"), stopTask)            // 停止任务
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

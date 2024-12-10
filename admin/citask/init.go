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
	apiGroup.Get("/citask", middleware.HasPermission("citask:list"), getTasks)                 // 获取任务列表
	apiGroup.Post("/citask", middleware.HasPermission("citask:create"), createTask)            // 创建任务
	apiGroup.Get("/citask/:id", middleware.HasPermission("citask:list"), getTask)              // 获取任务详情
	apiGroup.Put("/citask/:id", middleware.HasPermission("citask:update"), updateTask)         // 更新任务
	apiGroup.Delete("/citask/:id", middleware.HasPermission("citask:delete"), deleteTask)      // 删除任务
	apiGroup.Post("/citask/:id/run", middleware.HasPermission("citask:run"), runTask)          // 执行任务
	apiGroup.Get("/citask/:id/logs", middleware.HasPermission("citask:list"), getTaskLogs)     // 获取任务日志
	apiGroup.Get("/citask/:id/status", middleware.HasPermission("citask:list"), getTaskStatus) // 获取任务状态
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

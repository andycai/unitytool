package adminlog

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.adminlog.gorm.db"
	KeyNoCheckRouter = "admin.adminlog.router.nocheck"
	KeyCheckRouter   = "admin.adminlog.router.check"
)

var db *gorm.DB

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	// 管理员日志页面
	adminGroup.Get("/adminlog", middleware.HasPermission("adminlog:list"), func(c *fiber.Ctx) error {
		return c.Render("admin/adminlog", fiber.Map{
			"Title": "操作日志",
			"Scripts": []string{
				"/static/js/admin/adminlog.js",
			},
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	// 获取日志列表
	apiGroup.Get("/adminlog", middleware.HasPermission("adminlog:list"), getAdminLogs)

	// 删除日志
	apiGroup.Delete("/adminlog", middleware.HasPermission("adminlog:delete"), deleteAdminLogs)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

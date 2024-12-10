package menu

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	KeyDB            = "admin.menu.gorm.db"
	KeyNoCheckRouter = "admin.menu.router.nocheck"
	KeyCheckRouter   = "admin.menu.router.check"
)

var db *gorm.DB
var menuDao *MenuDao

func initDB(dbs []*gorm.DB) {
	db = dbs[0]
	menuDao = NewMenuDao(db)
}

func initAdminCheckRouter(adminGroup fiber.Router) {
	adminGroup.Get("/menus", middleware.HasPermission("menu:list"), func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		return c.Render("admin/menus", fiber.Map{
			"Title": "菜单管理",
			"Scripts": []string{
				"/static/js/admin/menus.js",
			},
			"user": user,
		}, "admin/layout")
	})
}

func initAPICheckRouter(apiGroup fiber.Router) {
	apiGroup.Get("/menus", middleware.HasPermission("menu:list"), listMenus)
	apiGroup.Get("/menus/tree", middleware.HasPermission("menu:list"), getMenuTree)
	apiGroup.Post("/menus", middleware.HasPermission("menu:create"), createMenu)
	apiGroup.Put("/menus/:id", middleware.HasPermission("menu:update"), updateMenu)
	apiGroup.Delete("/menus/:id", middleware.HasPermission("menu:delete"), deleteMenu)
}

func init() {
	core.RegisterDatabase(KeyDB, initDB)
	core.RegisterAdminCheckRouter(KeyCheckRouter, initAdminCheckRouter)
	core.RegisterAPICheckRouter(KeyCheckRouter, initAPICheckRouter)
}

package menu

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var app *core.App
var menuDao *MenuDao

type menuModule struct {
}

func (u *menuModule) Init(a *core.App) error {
	app = a
	menuDao = NewMenuDao()
	return nil
}

func (u *menuModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *menuModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *menuModule) InitRouter() error {
	// public

	// admin
	app.RouterAdmin.Get("/menus", middleware.HasPermission("menu:list"), func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)
		return c.Render("admin/menus", fiber.Map{
			"Title": "菜单管理",
			"Scripts": []string{
				"/static/js/admin/menus.js",
			},
			"user": user,
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/menus", middleware.HasPermission("menu:list"), listMenus)
	app.RouterApi.Get("/menus/tree", middleware.HasPermission("menu:list"), getMenuTree)
	app.RouterApi.Post("/menus", middleware.HasPermission("menu:create"), createMenu)
	app.RouterApi.Put("/menus/:id", middleware.HasPermission("menu:update"), updateMenu)
	app.RouterApi.Delete("/menus/:id", middleware.HasPermission("menu:delete"), deleteMenu)

	return nil
}

func init() {
	core.RegisterModules(&menuModule{})
}

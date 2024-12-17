package menu

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
)

var app *core.App
var menuDao *MenuDao

type menuModule struct {
}

func (m *menuModule) Awake(a *core.App) error {
	app = a
	if err := autoMigrate(); err != nil {
		return err
	}
	menuDao = NewMenuDao()

	return initData()
}

func (m *menuModule) Start() error {
	return nil
}

func (m *menuModule) AddPublicRouters() error {
	app.RouterPublicApi.Get("/menus/public/tree", getMenuTree)
	return nil
}

func (m *menuModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/menus", app.HasPermission("menu:list"), func(c *fiber.Ctx) error {
		user := app.CurrentUser(c)

		return c.Render("admin/menus", fiber.Map{
			"Title": "菜单管理",
			"Scripts": []string{
				"/static/js/admin/menus.js",
			},
			"user": user,
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/menus", app.HasPermission("menu:list"), listMenus)
	app.RouterApi.Get("/menus/tree", app.HasPermission("menu:list"), getMenuTree)
	app.RouterApi.Post("/menus", app.HasPermission("menu:create"), createMenu)
	app.RouterApi.Put("/menus/:id", app.HasPermission("menu:update"), updateMenu)
	app.RouterApi.Delete("/menus/:id", app.HasPermission("menu:delete"), deleteMenu)

	return nil
}

func init() {
	core.RegisterModule(&menuModule{}, 997)
}

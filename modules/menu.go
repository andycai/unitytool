package modules

import (
	"fmt"

	"github.com/andycai/unitool/dao"
	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/middleware"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var menuDao *dao.MenuDao

// InitMenuModule 初始化菜单模块
func InitMenuModule(app *fiber.App, d *dao.MenuDao) {
	menuDao = d
	RegisterMenuRoutes(app)
}

// RegisterMenuRoutes 注册菜单路由
func RegisterMenuRoutes(app *fiber.App) {
	// 菜单管理页面路由
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

	// 菜单 API 路由
	apiGroup.Get("/menus", middleware.HasPermission("menu:list"), listMenus)
	apiGroup.Get("/menus/tree", middleware.HasPermission("menu:list"), getMenuTree)
	apiGroup.Post("/menus", middleware.HasPermission("menu:create"), createMenu)
	apiGroup.Put("/menus/:id", middleware.HasPermission("menu:update"), updateMenu)
	apiGroup.Delete("/menus/:id", middleware.HasPermission("menu:delete"), deleteMenu)
}

// listMenus 获取菜单列表
func listMenus(c *fiber.Ctx) error {
	menus, err := menuDao.GetMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取菜单列表失败",
		})
	}
	return c.JSON(menus)
}

// getMenuTree 获取菜单树
func getMenuTree(c *fiber.Ctx) error {
	menus, err := menuDao.GetMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取菜单列表失败",
		})
	}
	tree := menuDao.BuildMenuTree(menus, 0)
	return c.JSON(tree)
}

// createMenu 创建菜单
func createMenu(c *fiber.Ctx) error {
	menu := new(models.Menu)
	if err := c.BodyParser(menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := menuDao.CreateMenu(menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建菜单失败",
		})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	handlers.CreateAdminLog(c, menuDao.DB, currentUser, "create", "menu", menu.ID, fmt.Sprintf("创建菜单：%s", menu.Name))

	return c.JSON(menu)
}

// updateMenu 更新菜单
func updateMenu(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的菜单ID",
		})
	}

	menu := new(models.Menu)
	if err := c.BodyParser(menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	menu.ID = uint(id)
	if err := menuDao.UpdateMenu(menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新菜单失败",
		})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	handlers.CreateAdminLog(c, menuDao.DB, currentUser, "update", "menu", menu.ID, fmt.Sprintf("更新菜单：%s", menu.Name))

	return c.JSON(menu)
}

// deleteMenu 删除菜单
func deleteMenu(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的菜单ID",
		})
	}

	// 获取菜单信息用于日志记录
	menu, err := menuDao.GetMenuByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "菜单不存在",
		})
	}

	if err := menuDao.DeleteMenu(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除菜单失败",
		})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	handlers.CreateAdminLog(c, menuDao.DB, currentUser, "delete", "menu", uint(id), fmt.Sprintf("删除菜单：%s", menu.Name))

	return c.SendStatus(fiber.StatusNoContent)
}

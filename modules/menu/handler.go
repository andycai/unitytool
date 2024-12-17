package menu

import (
	"fmt"

	"github.com/andycai/unitool/models"
	"github.com/andycai/unitool/modules/adminlog"
	"github.com/gofiber/fiber/v2"
)

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
	adminlog.CreateAdminLog(c, "create", "menu", menu.ID, fmt.Sprintf("创建菜单：%s", menu.Name))

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
	adminlog.CreateAdminLog(c, "update", "menu", menu.ID, fmt.Sprintf("更新菜单：%s", menu.Name))

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
	adminlog.CreateAdminLog(c, "delete", "menu", uint(id), fmt.Sprintf("删除菜单：%s", menu.Name))

	return c.SendStatus(fiber.StatusNoContent)
}

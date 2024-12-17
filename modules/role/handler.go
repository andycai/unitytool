package role

import (
	"fmt"
	"time"

	"github.com/andycai/unitool/models"
	"github.com/andycai/unitool/modules/adminlog"
	"github.com/gofiber/fiber/v2"
)

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions []uint `json:"permissions"` // 权限ID列表
}

type UpdateRoleRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Permissions []uint `json:"permissions,omitempty"`
}

// getRoles 获取角色列表
func getRoles(c *fiber.Ctx) error {
	var roles []models.Role
	if err := app.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取角色列表失败"})
	}
	return c.JSON(roles)
}

// createRole 创建角色
func createRole(c *fiber.Ctx) error {
	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 检查角色名是否已存在
	var count int64
	app.DB.Model(&models.Role{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "角色名已存在"})
	}

	// 开始事务
	tx := app.DB.Begin()

	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "创建角色失败"})
	}

	// 添加权限关联
	if len(req.Permissions) > 0 {
		var permissions []models.Permission
		if err := tx.Find(&permissions, req.Permissions).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "获取权限失败"})
		}

		if err := tx.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "设置权限失败"})
		}
	}

	tx.Commit()

	// 记录操作日志
	adminlog.CreateAdminLog(c, "create", "role", role.ID, fmt.Sprintf("创建角色：%s", role.Name))

	return c.JSON(role)
}

// updateRole 更新角色
func updateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	var role models.Role
	if err := app.DB.First(&role, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "角色不存在"})
	}

	// 开始事务
	tx := app.DB.Begin()

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := tx.Model(&role).Updates(updates).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "更新角色失败"})
	}

	// 更新权限关联
	if len(req.Permissions) > 0 {
		var permissions []models.Permission
		if err := tx.Find(&permissions, req.Permissions).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "获取权限失败"})
		}

		if err := tx.Model(&role).Association("Permissions").Replace(permissions); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "更新权限失败"})
		}
	}

	tx.Commit()

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "role", role.ID, fmt.Sprintf("更新角色：%s", role.Name))

	return c.JSON(role)
}

// deleteRole 删除角色
func deleteRole(c *fiber.Ctx) error {
	id := c.Params("id")

	// 检查是否有用户使用此角色
	var count int64
	if err := app.DB.Model(&models.User{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "检查角色使用状态失败"})
	}

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "该角色正在使用中，无法删除"})
	}

	var role models.Role
	if err := app.DB.First(&role, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "角色不存在"})
	}

	// 开始事务
	tx := app.DB.Begin()

	// 清除权限关联
	if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "清除权限关联失败"})
	}

	// 删除角色
	if err := tx.Delete(&role).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "删除角色失败"})
	}

	tx.Commit()

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "role", role.ID, fmt.Sprintf("删除角色：%s", role.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

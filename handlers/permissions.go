package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/andycai/unitool/models"
)

type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

// GetPermissions 获取权限列表
func GetPermissions(c *fiber.Ctx, db *gorm.DB) error {
	var permissions []models.Permission
	if err := db.Find(&permissions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取权限列表失败"})
	}
	return c.JSON(permissions)
}

// CreatePermission 创建权限
func CreatePermission(c *fiber.Ctx, db *gorm.DB) error {
	var req CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 检查权限编码是否已存在
	var count int64
	db.Model(&models.Permission{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "权限编码已存在"})
	}

	permission := models.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&permission).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建权限失败"})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	CreateAdminLog(c, db, currentUser, "create", "permission", permission.ID, fmt.Sprintf("创建权限：%s", permission.Name))

	return c.JSON(permission)
}

// UpdatePermission 更新权限
func UpdatePermission(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var req UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	var permission models.Permission
	if err := db.First(&permission, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "权限不存在"})
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		// 检查新的权限编码是否已存在
		var count int64
		db.Model(&models.Permission{}).Where("code = ? AND id != ?", req.Code, id).Count(&count)
		if count > 0 {
			return c.Status(400).JSON(fiber.Map{"error": "权限编码已存在"})
		}
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := db.Model(&permission).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新权限失败"})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	CreateAdminLog(c, db, currentUser, "update", "permission", permission.ID, fmt.Sprintf("更新权限：%s", permission.Name))

	return c.JSON(permission)
}

// DeletePermission 删除权限
func DeletePermission(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")

	// 检查权限是否被角色使用
	var count int64
	if err := db.Table("role_permissions").Where("permission_id = ?", id).Count(&count).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "检查权限使用状态失败"})
	}

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "该权限正在被角色使用，无法删除"})
	}

	var permission models.Permission
	if err := db.First(&permission, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "权限不存在"})
	}

	if err := db.Delete(&permission).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除权限失败"})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	CreateAdminLog(c, db, currentUser, "delete", "permission", permission.ID, fmt.Sprintf("删除权限：%s", permission.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

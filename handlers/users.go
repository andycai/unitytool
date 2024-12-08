package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/andycai/unitool/models"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	RoleID   uint   `json:"role_id"`
}

type UpdateUserRequest struct {
	Password string `json:"password,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	RoleID   uint   `json:"role_id,omitempty"`
	Status   *int   `json:"status,omitempty"`
}

// GetUsers 获取用户列表
func GetUsers(c *fiber.Ctx, db *gorm.DB) error {
	var users []models.User
	if err := db.Preload("Role").Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取用户列表失败"})
	}
	return c.JSON(users)
}

// CreateUser 创建用户
func CreateUser(c *fiber.Ctx, db *gorm.DB) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 检查用户名是否已存在
	var count int64
	db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "用户名已存在"})
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "密码加密失败"})
	}

	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Nickname:  req.Nickname,
		RoleID:    req.RoleID,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建用户失败"})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	CreateAdminLog(c, db, currentUser, "create", "user", user.ID, fmt.Sprintf("创建用户：%s", user.Username))

	return c.JSON(user)
}

// UpdateUser 更新用户
func UpdateUser(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.RoleID != 0 {
		updates["role_id"] = req.RoleID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "密码加密失败"})
		}
		updates["password"] = string(hashedPassword)
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "用户不存在"})
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新用户失败"})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	CreateAdminLog(c, db, currentUser, "update", "user", user.ID, fmt.Sprintf("更新用户：%s", user.Username))

	return c.JSON(user)
}

// DeleteUser 删除用户
func DeleteUser(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	var user models.User

	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "用户不存在",
		})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "删除用户失败",
		})
	}

	// 记录操作日志
	currentUser := c.Locals("user").(models.User)
	CreateAdminLog(c, db, currentUser, "delete", "user", user.ID, fmt.Sprintf("删除用户：%s", user.Username))

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}

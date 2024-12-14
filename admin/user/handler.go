package user

import (
	"fmt"
	"time"

	"github.com/andycai/unitool/admin/adminlog"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

// getUsersAction 获取用户列表
func getUsersAction(c *fiber.Ctx) error {
	var users []models.User
	if err := app.DB.Preload("Role").Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取用户列表失败"})
	}
	return c.JSON(users)
}

// createUserAction 创建用户
func createUserAction(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 检查用户名是否已存在
	var count int64
	app.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
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

	if err := app.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建用户失败"})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "create", "user", user.ID, fmt.Sprintf("创建用户：%s", user.Username))

	return c.JSON(user)
}

// updateUserAction 更新用户
func updateUserAction(c *fiber.Ctx) error {
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
	if err := app.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "用户不存在"})
	}

	if err := app.DB.Model(&user).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新用户失败"})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "user", user.ID, fmt.Sprintf("更新用户：%s", user.Username))

	return c.JSON(user)
}

// deleteUserAction 删除用户
func deleteUserAction(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := app.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "用户不存在",
		})
	}

	if err := app.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "删除用户失败",
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "user", user.ID, fmt.Sprintf("删除用户：%s", user.Username))

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}

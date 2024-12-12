package login

import (
	"time"

	"github.com/andycai/unitool/models"
	"github.com/andycai/unitool/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// login 处理登录请求
func login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	var user models.User
	if err := db.Preload("Role.Permissions").Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "用户名或密码错误"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "用户名或密码错误"})
	}

	// 打印用户信息用于调试
	// fmt.Printf("User info: %+v\n", user)

	// 生成 JWT token
	claims := jwt.MapClaims{
		"sub":      user.ID, // 使用标准声明
		"user_id":  user.ID,
		"username": user.Username,
		"role_id":  user.RoleID,
		"exp":      time.Now().Add(time.Duration(utils.GetConfig().Auth.TokenExpire) * time.Second).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(utils.GetConfig().Auth.JWTSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "生成token失败"})
	}

	// 更新最后登录时间
	db.Model(&user).Update("last_login", time.Now())

	// 清除密码字段
	user.Password = ""

	// 生成 token 后，设置 cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(time.Duration(utils.GetConfig().Auth.TokenExpire) * time.Second)
	cookie.HTTPOnly = true
	cookie.Secure = true // 如果是 HTTPS
	cookie.Path = "/"
	c.Cookie(cookie)

	// 构建响应数据
	responseData := fiber.Map{
		"code":    0,
		"message": "登录成功",
		"data": fiber.Map{
			"token": tokenString,
			"user": fiber.Map{
				"id":              user.ID,
				"username":        user.Username,
				"nickname":        user.Nickname,
				"role_id":         user.RoleID,
				"role":            user.Role,
				"status":          user.Status,
				"last_login":      user.LastLogin,
				"has_changed_pwd": user.HasChangedPwd,
				"created_at":      user.CreatedAt,
				"updated_at":      user.UpdatedAt,
			},
		},
	}

	// 打印响应数据用于调试
	// fmt.Printf("Response data: %+v\n", responseData)

	// 返回响应
	return c.JSON(responseData)
}

// changePassword 修改密码
func changePassword(c *fiber.Ctx) error {
	var req ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 获取用户信息
	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "用户不存在"})
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "当前密码错误"})
	}

	// 生成新密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "密码加密失败"})
	}

	// 更新密码和密码修改状态
	updates := map[string]interface{}{
		"password":        string(hashedPassword),
		"has_changed_pwd": true,
		"updated_at":      time.Now(),
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新密码失败"})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "密码修改成功",
	})
}

// 生成 JWT token
func generateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role_id":  user.RoleID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.GetConfig().Auth.JWTSecret))
}

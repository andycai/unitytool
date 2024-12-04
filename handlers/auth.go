package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mind.com/log/models"
	"mind.com/log/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Login 处理登录请求
func Login(c *fiber.Ctx, db *gorm.DB) error {
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

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "登录成功",
		"data": fiber.Map{
			"token": tokenString,
			"user":  user,
		},
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

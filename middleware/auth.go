package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"mind.com/log/models"
	"mind.com/log/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取 token
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "未授权访问"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.GetConfig().Auth.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "无效的token"})
		}

		// 获取用户信息
		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))

		var user models.User
		if err := db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "用户不存在"})
		}

		// 将用户信息存储到上下文
		c.Locals("user", user)

		return c.Next()
	}
}

// HasPermission 权限检查中间件
func HasPermission(permissionCode string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)

		// 检查用户权限
		hasPermission := false
		for _, perm := range user.Role.Permissions {
			if perm.Code == permissionCode {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(403).JSON(fiber.Map{"error": "没有权限"})
		}

		return c.Next()
	}
}

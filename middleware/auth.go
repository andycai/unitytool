package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mind.com/log/models"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取 token
		// authHeader := c.Get("Authorization")
		// log.Printf("Authorization header: %s", authHeader)

		// // 检查所有请求头
		// headers := c.GetReqHeaders()
		// log.Printf("All headers: %v", headers)

		// if authHeader == "" {
		// 	// 对于页面请求，重定向到登录页
		// 	if c.Method() == "GET" && !strings.HasPrefix(c.Path(), "/api/") {
		// 		return c.Redirect("/login")
		// 	}
		// 	return c.Status(401).JSON(fiber.Map{"error": "未授权访问"})
		// }

		// // 解析 token
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// claims := jwt.MapClaims{}
		// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(utils.GetConfig().Auth.JWTSecret), nil
		// })

		// // token 验证失败
		// if err != nil || !token.Valid {
		// 	// 对于页面请求，重定向到登录页
		// 	if c.Method() == "GET" && !strings.HasPrefix(c.Path(), "/api/") {
		// 		return c.Redirect("/login")
		// 	}
		// 	return c.Status(401).JSON(fiber.Map{"error": "无效的token"})
		// }

		// // 获取用户信息
		// userID := uint(claims["user_id"].(float64))
		// var user models.User
		// if err := db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
		// 	// 对于页面请求，重定向到登录页
		// 	if c.Method() == "GET" && !strings.HasPrefix(c.Path(), "/api/") {
		// 		return c.Redirect("/login")
		// 	}
		// 	return c.Status(401).JSON(fiber.Map{"error": "用户不存在"})
		// }

		// // 将用户信息存储到上下文
		// c.Locals("user", user)

		// 继续处理请求
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

package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// SkipAuth 用于标记不需要认证的路由
var skipAuthRoutes = make(map[string]bool)

// SkipAuthRoute 添加不需要认证的路由
func SkipAuthRoute(path string) {
	skipAuthRoutes[path] = true
}

func AuthMiddleware(c *fiber.Ctx) error {
	// 检查是否是不需要认证的路由
	if skipAuthRoutes[c.Path()] {
		return c.Next()
	}

	isAuthenticated, _ := GetSession(c)

	if isAuthenticated {
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "未授权访问",
	})
}

// AuthMiddleware 认证中间件
func AuthMiddleware2(db *gorm.DB, JWTSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 首先从请求头获取 token
		authHeader := c.Get("Authorization")

		// 如果请求头没有 token，尝试从 cookie 获取
		if authHeader == "" {
			if token := c.Cookies("token"); token != "" {
				authHeader = "Bearer " + token
				// 设置到请求头中，方便后续使用
				c.Request().Header.Set("Authorization", authHeader)
			}
		}

		// 如果还是没有 token
		if authHeader == "" {
			// 对于页面请求，重定向到登录页
			if c.Method() == "GET" && !strings.HasPrefix(c.Path(), "/api/") {
				return c.Redirect("/login")
			}
			return c.Status(401).JSON(fiber.Map{"error": "未授权访问"})
		}

		// 解析 token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		})

		// token 验证失败
		if err != nil || !token.Valid {
			// 对于页面请求，重定向到登录页
			if c.Method() == "GET" && !strings.HasPrefix(c.Path(), "/api/") {
				return c.Redirect("/login")
			}
			return c.Status(401).JSON(fiber.Map{"error": fmt.Sprintf("无效的token: %v", err)})
		}

		// 获取用户信息
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "无效的token格式"})
		}

		userID := uint(claims["user_id"].(float64))
		var user models.User
		if err := db.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
			// 对于页面请求，重定向到登录页
			if c.Method() == "GET" && !strings.HasPrefix(c.Path(), "/api/") {
				return c.Redirect("/login")
			}
			return c.Status(401).JSON(fiber.Map{"error": "用户不存在"})
		}

		// 将用户信息存储到上下文
		c.Locals("user", user)

		// 继续处理请求
		return c.Next()
	}
}

// HasPermission 权限检查中间件
func HasPermission(permissionCode string, userFunc func(c *fiber.Ctx) *models.User) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if userFunc == nil || userFunc(c) == nil {
			return errors.New("请先登录")
		}
		currentUser := userFunc(c)

		// 检查用户权限
		hasPermission := false
		for _, perm := range currentUser.Role.Permissions {
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

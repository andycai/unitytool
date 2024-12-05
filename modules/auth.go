package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
	"mind.com/log/middleware"
)

type AuthModule struct {
	BaseModule
}

func (m *AuthModule) Init() error {
	return nil
}

func (m *AuthModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// 登录页面路由
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	// 管理后台页面路由组
	admin := app.Group("/admin", middleware.AuthMiddleware(m.DB))

	// 管理后台主页路由
	admin.Get("/", func(c *fiber.Ctx) error {
		return c.Render("admin/index", fiber.Map{
			"Title": "管理后台",
		}, "admin/layout")
	})

	// 其他管理后台页面路由
	admin.Get("/users", func(c *fiber.Ctx) error {
		return c.Render("admin/users", fiber.Map{
			"Title": "用户管理",
			"Scripts": []string{
				"/static/js/admin/users.js",
			},
		}, "admin/layout")
	})

	admin.Get("/roles", func(c *fiber.Ctx) error {
		return c.Render("admin/roles", fiber.Map{
			"Title": "角色管理",
			"Scripts": []string{
				"/static/js/admin/roles.js",
			},
		}, "admin/layout")
	})

	admin.Get("/permissions", func(c *fiber.Ctx) error {
		return c.Render("admin/permissions", fiber.Map{
			"Title": "权限管理",
			"Scripts": []string{
				"/static/js/admin/permissions.js",
			},
		}, "admin/layout")
	})

	// API 路由
	api := app.Group("/api")

	// 公开路由
	api.Post("/login", func(c *fiber.Ctx) error {
		return handlers.Login(c, m.DB)
	})

	// 需要认证的路由
	auth := api.Use(middleware.AuthMiddleware(m.DB))

	// 用户管理
	auth.Get("/users", middleware.HasPermission("user:list"), func(c *fiber.Ctx) error {
		return handlers.GetUsers(c, m.DB)
	})
	auth.Post("/users", middleware.HasPermission("user:create"), func(c *fiber.Ctx) error {
		return handlers.CreateUser(c, m.DB)
	})
	auth.Put("/users/:id", middleware.HasPermission("user:update"), func(c *fiber.Ctx) error {
		return handlers.UpdateUser(c, m.DB)
	})
	auth.Delete("/users/:id", middleware.HasPermission("user:delete"), func(c *fiber.Ctx) error {
		return handlers.DeleteUser(c, m.DB)
	})

	// 角色管理
	auth.Get("/roles", middleware.HasPermission("role:list"), func(c *fiber.Ctx) error {
		return handlers.GetRoles(c, m.DB)
	})
	auth.Post("/roles", middleware.HasPermission("role:create"), func(c *fiber.Ctx) error {
		return handlers.CreateRole(c, m.DB)
	})
	auth.Put("/roles/:id", middleware.HasPermission("role:update"), func(c *fiber.Ctx) error {
		return handlers.UpdateRole(c, m.DB)
	})
	auth.Delete("/roles/:id", middleware.HasPermission("role:delete"), func(c *fiber.Ctx) error {
		return handlers.DeleteRole(c, m.DB)
	})

	// 权限管理
	auth.Get("/permissions", middleware.HasPermission("permission:list"), func(c *fiber.Ctx) error {
		return handlers.GetPermissions(c, m.DB)
	})
	auth.Post("/permissions", middleware.HasPermission("permission:create"), func(c *fiber.Ctx) error {
		return handlers.CreatePermission(c, m.DB)
	})
	auth.Put("/permissions/:id", middleware.HasPermission("permission:update"), func(c *fiber.Ctx) error {
		return handlers.UpdatePermission(c, m.DB)
	})
	auth.Delete("/permissions/:id", middleware.HasPermission("permission:delete"), func(c *fiber.Ctx) error {
		return handlers.DeletePermission(c, m.DB)
	})
}

package menu

import (
	"time"

	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var app *core.App
var menuDao *MenuDao

type menuModule struct {
}

func (m *menuModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := app.DB.AutoMigrate(&models.Menu{}); err != nil {
		return err
	}
	menuDao = NewMenuDao()

	return initData()
}

func initData() error {
	var count int64
	app.DB.Model(&models.Menu{}).Count(&count)
	if count > 0 {
		return nil
	}

	// 初始化数据
	now := time.Now()

	// 系统管理菜单组
	systemManage := models.Menu{
		ParentID:   0,
		Name:       "系统管理",
		Path:       "/admin",
		Icon:       "system",
		Sort:       1,
		Permission: "",
		IsShow:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := app.DB.Create(&systemManage).Error; err != nil {
		return err
	}

	// 系统管理子菜单
	systemMenus := []models.Menu{
		{
			ParentID:   systemManage.ID,
			Name:       "用户管理",
			Path:       "/admin/users",
			Icon:       "user",
			Sort:       1,
			Permission: "user:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   systemManage.ID,
			Name:       "角色管理",
			Path:       "/admin/roles",
			Icon:       "role",
			Sort:       2,
			Permission: "role:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   systemManage.ID,
			Name:       "权限管理",
			Path:       "/admin/permissions",
			Icon:       "permission",
			Sort:       3,
			Permission: "permission:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   systemManage.ID,
			Name:       "菜单管理",
			Path:       "/admin/menus",
			Icon:       "menu",
			Sort:       4,
			Permission: "menu:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   systemManage.ID,
			Name:       "操作日志",
			Path:       "/admin/adminlog",
			Icon:       "adminlog",
			Sort:       5,
			Permission: "adminlog:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := app.DB.Create(&systemMenus).Error; err != nil {
		return err
	}

	// 游戏管理菜单组
	gameManage := models.Menu{
		ParentID:   0,
		Name:       "游戏管理",
		Path:       "/admin/game",
		Icon:       "game",
		Sort:       2,
		Permission: "",
		IsShow:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := app.DB.Create(&gameManage).Error; err != nil {
		return err
	}

	// 游戏管理子菜单
	gameMenus := []models.Menu{
		{
			ParentID:   gameManage.ID,
			Name:       "游戏日志",
			Path:       "/admin/gamelog",
			Icon:       "gamelog",
			Sort:       1,
			Permission: "gamelog:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   gameManage.ID,
			Name:       "性能统计",
			Path:       "/admin/stats",
			Icon:       "stats",
			Sort:       2,
			Permission: "stats:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := app.DB.Create(&gameMenus).Error; err != nil {
		return err
	}

	// 系统工具菜单组
	toolsManage := models.Menu{
		ParentID:   0,
		Name:       "系统工具",
		Path:       "/admin/tools",
		Icon:       "tools",
		Sort:       3,
		Permission: "",
		IsShow:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := app.DB.Create(&toolsManage).Error; err != nil {
		return err
	}

	// 系统工具子菜单
	toolsMenus := []models.Menu{
		{
			ParentID:   toolsManage.ID,
			Name:       "构建任务",
			Path:       "/admin/citask",
			Icon:       "citask",
			Sort:       1,
			Permission: "citask:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   toolsManage.ID,
			Name:       "文件浏览",
			Path:       "/admin/browse",
			Icon:       "browse",
			Sort:       2,
			Permission: "browse:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ParentID:   toolsManage.ID,
			Name:       "服务器配置",
			Path:       "/admin/serverconf",
			Icon:       "serverconf",
			Sort:       3,
			Permission: "serverconf:list",
			IsShow:     true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := app.DB.Create(&toolsMenus).Error; err != nil {
		return err
	}

	return nil
}

func (m *menuModule) Start() error {
	return nil
}

func (m *menuModule) AddPublicRouters() error {
	return nil
}

func (m *menuModule) AddAuthRouters() error {
	// admin
	app.RouterAdmin.Get("/menus", app.HasPermission("menu:list"), func(c *fiber.Ctx) error {
		user := app.CurrentUser(c)

		return c.Render("admin/menus", fiber.Map{
			"Title": "菜单管理",
			"Scripts": []string{
				"/static/js/admin/menus.js",
			},
			"user": user,
		}, "admin/layout")
	})

	// api
	app.RouterApi.Get("/menus", app.HasPermission("menu:list"), listMenus)
	app.RouterApi.Get("/menus/tree", app.HasPermission("menu:list"), getMenuTree)
	app.RouterApi.Post("/menus", app.HasPermission("menu:create"), createMenu)
	app.RouterApi.Put("/menus/:id", app.HasPermission("menu:update"), updateMenu)
	app.RouterApi.Delete("/menus/:id", app.HasPermission("menu:delete"), deleteMenu)

	return nil
}

func init() {
	core.RegisterModule(&menuModule{})
}

package menu

import (
	"log"
	"time"

	"github.com/andycai/unitool/models"
)

type MenuDao struct {
}

func NewMenuDao() *MenuDao {
	return &MenuDao{}
}

// 数据迁移
func autoMigrate() error {
	return app.DB.AutoMigrate(&models.Menu{})
}

// 初始化数据
func initData() error {
	// 检查是否已初始化
	if app.IsInitializedModule("menu") {
		log.Println("菜单模块数据库已初始化，跳过")
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

	// 标记模块已初始化
	if err := app.DB.Create(&models.ModuleInit{
		Module:      "menu",
		Initialized: 1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}).Error; err != nil {
		return err
	}

	return nil
}

// GetMenus 获取所有菜单
func (d *MenuDao) GetMenus() ([]*models.Menu, error) {
	var menus []*models.Menu
	result := app.DB.Order("sort asc").Find(&menus)
	return menus, result.Error
}

// GetMenusByPermissions 根据权限获取菜单
func (d *MenuDao) GetMenusByPermissions(permissions []string) ([]*models.Menu, error) {
	var menus []*models.Menu
	result := app.DB.Where("permission IN ? OR permission = ''", permissions).
		Where("is_show = ?", true).
		Order("sort asc").
		Find(&menus)
	return menus, result.Error
}

// BuildMenuTree 构建菜单树
func (d *MenuDao) BuildMenuTree(menus []*models.Menu, parentID uint) []*models.MenuTree {
	var tree []*models.MenuTree
	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := &models.MenuTree{
				Menu:     menu,
				Children: d.BuildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// CreateMenu 创建菜单
func (d *MenuDao) CreateMenu(menu *models.Menu) error {
	return app.DB.Create(menu).Error
}

// UpdateMenu 更新菜单
func (d *MenuDao) UpdateMenu(menu *models.Menu) error {
	return app.DB.Save(menu).Error
}

// DeleteMenu 删除菜单
func (d *MenuDao) DeleteMenu(id uint) error {
	// 先删除子菜单
	if err := app.DB.Where("parent_id = ?", id).Delete(&models.Menu{}).Error; err != nil {
		return err
	}
	// 再删除当前菜单
	return app.DB.Delete(&models.Menu{}, id).Error
}

// GetMenuByID 根据ID获取菜单
func (d *MenuDao) GetMenuByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	result := app.DB.First(&menu, id)
	return &menu, result.Error
}

package menu

import (
	"github.com/andycai/unitool/models"
)

type MenuDao struct {
}

func NewMenuDao() *MenuDao {
	return &MenuDao{}
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

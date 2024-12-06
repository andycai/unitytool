package dao

import (
	"github.com/andycai/unitool/models"
	"gorm.io/gorm"
)

type MenuDao struct {
	DB *gorm.DB
}

func NewMenuDao(db *gorm.DB) *MenuDao {
	return &MenuDao{DB: db}
}

// GetMenus 获取所有菜单
func (d *MenuDao) GetMenus() ([]*models.Menu, error) {
	var menus []*models.Menu
	result := d.DB.Order("sort asc").Find(&menus)
	return menus, result.Error
}

// GetMenusByPermissions 根据权限获取菜单
func (d *MenuDao) GetMenusByPermissions(permissions []string) ([]*models.Menu, error) {
	var menus []*models.Menu
	result := d.DB.Where("permission IN ? OR permission = ''", permissions).
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
	return d.DB.Create(menu).Error
}

// UpdateMenu 更新菜单
func (d *MenuDao) UpdateMenu(menu *models.Menu) error {
	return d.DB.Save(menu).Error
}

// DeleteMenu 删除菜单
func (d *MenuDao) DeleteMenu(id uint) error {
	// 先删除子菜单
	if err := d.DB.Where("parent_id = ?", id).Delete(&models.Menu{}).Error; err != nil {
		return err
	}
	// 再删除当前菜单
	return d.DB.Delete(&models.Menu{}, id).Error
}

// GetMenuByID 根据ID获取菜单
func (d *MenuDao) GetMenuByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	result := d.DB.First(&menu, id)
	return &menu, result.Error
}

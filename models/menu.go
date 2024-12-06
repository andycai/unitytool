package models

import (
	"time"
)

// Menu 菜单模型
type Menu struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	ParentID   uint      `gorm:"index" json:"parent_id"`        // 父菜单ID，0表示根菜单
	Name       string    `gorm:"size:50;not null" json:"name"`  // 菜单名称
	Path       string    `gorm:"size:100;not null" json:"path"` // 菜单路径
	Icon       string    `gorm:"size:100" json:"icon"`          // 菜单图标
	Sort       int       `gorm:"default:0" json:"sort"`         // 排序
	Permission string    `gorm:"size:100" json:"permission"`    // 绑定的权限
	IsShow     bool      `gorm:"default:true" json:"is_show"`   // 是否显示
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MenuTree 菜单树结构
type MenuTree struct {
	Menu     *Menu       `json:"menu"`
	Children []*MenuTree `json:"children"`
}

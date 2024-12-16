package models

import (
	"time"
)

// NoteCategory 笔记分类
type NoteCategory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	ParentID    uint      `json:"parent_id" gorm:"default:0"`
	IsPublic    uint8     `json:"is_public" gorm:"type:tinyint(1);default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`

	// 关联
	Parent   *NoteCategory   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []*NoteCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Notes    []*Note         `json:"notes,omitempty" gorm:"foreignKey:CategoryID"`
	Roles    []*Role         `json:"roles,omitempty" gorm:"many2many:category_permissions"`
}

// Note 笔记
type Note struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"size:200;not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	CategoryID uint      `json:"category_id"`
	ParentID   uint      `json:"parent_id" gorm:"default:0"`
	IsPublic   uint8     `json:"is_public" gorm:"type:tinyint(1);default:0"`
	ViewCount  int       `json:"view_count" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  uint      `json:"created_by"`
	UpdatedBy  uint      `json:"updated_by"`

	// 关联
	Category *NoteCategory `json:"category,omitempty"`
	Parent   *Note         `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []*Note       `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Roles    []*Role       `json:"roles,omitempty" gorm:"many2many:note_permissions"`
}

// NotePermission 笔记权限
type NotePermission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	NoteID    uint      `json:"note_id"`
	RoleID    uint      `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint      `json:"created_by"`

	// 关联
	Note *Note `json:"note,omitempty"`
	Role *Role `json:"role,omitempty"`
}

// CategoryPermission 分类权限
type CategoryPermission struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CategoryID uint      `json:"category_id"`
	RoleID     uint      `json:"role_id"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  uint      `json:"created_by"`

	// 关联
	Category *NoteCategory `json:"category,omitempty"`
	Role     *Role         `json:"role,omitempty"`
}

// TableName 指定表名
func (NoteCategory) TableName() string {
	return "note_categories"
}

func (Note) TableName() string {
	return "notes"
}

func (NotePermission) TableName() string {
	return "note_permissions"
}

func (CategoryPermission) TableName() string {
	return "category_permissions"
}

// 检查用户是否有权限访问笔记
func (n *Note) HasPermission(userRoles []uint) bool {
	// 如果笔记是公开的，直接返回true
	if n.IsPublic == 1 {
		return true
	}

	// 检查用户角色是否有权限
	for _, role := range n.Roles {
		for _, userRole := range userRoles {
			if role.ID == userRole {
				return true
			}
		}
	}

	return false
}

// 检查用户是否有权限访问分类
func (c *NoteCategory) HasPermission(userRoles []uint) bool {
	// 如果分类是公开的，直接返回true
	if c.IsPublic == 1 {
		return true
	}

	// 检查用户角色是否有权限
	for _, role := range c.Roles {
		for _, userRole := range userRoles {
			if role.ID == userRole {
				return true
			}
		}
	}

	return false
}

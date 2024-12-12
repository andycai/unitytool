package models

import (
	"time"
)

// User 用户表
type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Username      string    `json:"username" gorm:"uniqueIndex;size:50"`
	Password      string    `json:"-" gorm:"size:100"` // 密码不返回给前端
	Nickname      string    `json:"nickname" gorm:"size:50"`
	RoleID        uint      `json:"role_id"`
	Role          Role      `json:"role" gorm:"foreignKey:RoleID"`
	Status        int       `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	LastLogin     time.Time `json:"last_login"`
	HasChangedPwd bool      `json:"has_changed_pwd" gorm:"default:false"` // 是否已修改初始密码
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Role 角色表
type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"uniqueIndex;size:50"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Permission 权限表
type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:50"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:50"` // 权限编码
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RolePermission 角色-权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`
}

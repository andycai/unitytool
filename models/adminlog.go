package models

import "time"

// AdminLog 管理员操作日志
type AdminLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`     // 操作用户ID
	Username   string    `json:"username"`    // 操作用户名
	Action     string    `json:"action"`      // 操作类型
	Resource   string    `json:"resource"`    // 资源类型
	ResourceID uint      `json:"resource_id"` // 资源ID
	Details    string    `json:"details"`     // 操作详情
	IP         string    `json:"ip"`          // 操作IP
	UserAgent  string    `json:"user_agent"`  // 用户代理
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
}

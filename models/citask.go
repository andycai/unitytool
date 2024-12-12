package models

import (
	"time"
)

// Task 任务表
type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`                 // 任务名称
	Description string    `json:"description" gorm:"type:text"`                  // 任务描述
	Type        string    `json:"type" gorm:"size:20;not null;default:'script'"` // 任务类型：script(脚本), http(远程调用)
	Script      string    `json:"script" gorm:"type:text"`                       // 脚本内容
	URL         string    `json:"url" gorm:"size:255"`                           // HTTP URL
	Method      string    `json:"method" gorm:"size:10;default:'GET'"`           // HTTP 方法
	Headers     string    `json:"headers" gorm:"type:text"`                      // HTTP 请求头
	Body        string    `json:"body" gorm:"type:text"`                         // HTTP 请求体
	Timeout     int       `json:"timeout" gorm:"default:300"`                    // 超时时间(秒)
	Status      string    `json:"status" gorm:"size:20;default:'active'"`        // 状态：active, inactive
	EnableCron  uint8     `json:"enable_cron" gorm:"type:tinyint;default:0"`     // 是否启用定时执行：0-否，1-是
	CronExpr    string    `json:"cron_expr"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskLog 任务执行日志表
type TaskLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"task_id" gorm:"index"`                    // 任务ID
	Task      Task      `json:"task" gorm:"foreignKey:TaskID"`           // 任务关联
	Status    string    `json:"status" gorm:"size:20;default:'pending'"` // 执行状态：success, failed, running
	Output    string    `json:"output" gorm:"type:text"`                 // 执行输出
	Error     string    `json:"error" gorm:"type:text"`                  // 错误信息
	StartTime time.Time `json:"start_time"`                              // 开始时间
	EndTime   time.Time `json:"end_time"`                                // 结束时间
	Duration  int       `json:"duration"`                                // 执行时长(秒)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

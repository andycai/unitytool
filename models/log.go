package models

import (
	"time"

	"gorm.io/gorm"
)

type Log struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	AppID      string    `json:"app_id"`
	Package    string    `json:"package"`
	RoleName   string    `json:"role_name"`
	Device     string    `json:"device"`
	LogType    string    `json:"log_type"`
	LogMessage string    `json:"log_message"`
	LogStack   string    `json:"log_stack"`
	LogTime    time.Time `json:"log_time" gorm:"type:datetime"`
}

func (l *Log) BeforeCreate(tx *gorm.DB) (err error) {
	if l.LogTime.IsZero() {
		l.LogTime = time.Now()
	}
	return
}

package models

type GameLog struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	AppID      string `json:"app_id"`
	Package    string `json:"package"`
	RoleName   string `json:"role_name" gorm:"index:idx_role_name_message, priority:1"`
	Device     string `json:"device"`
	LogMessage string `json:"log_message" gorm:"index:idx_role_name_message, priority:2"`
	LogTime    int64  `json:"log_time" gorm:"index"`
	LogType    string `json:"log_type"`
	LogStack   string `json:"log_stack"`
	CreateAt   int64  `json:"create_at" gorm:"index"`
}

package models

type Log struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	AppID      string `json:"app_id"`
	Package    string `json:"package"`
	RoleName   string `json:"role_name"`
	Device     string `json:"device"`
	LogMessage string `json:"log_message"`
	LogTime    int64  `json:"log_time"`
	LogType    string `json:"log_type"`
	LogStack   string `json:"log_stack"`
}

type StatsRecord struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	LoginID   int    `json:"login_id" gorm:"unique"`
	AppID     int    `json:"app_id"`
	Package   string `json:"package"`
	RoleName  string `json:"role_name"`
	Device    string `json:"device"`
	CPU       string `json:"cpu"`
	GPU       string `json:"gpu"`
	Memory    string `json:"memory"`
	CreatedAt int64  `json:"created_at"`
}

type StatsInfo struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	LoginID      int    `json:"login_id"`
	FPS          int    `json:"fps"`
	TotalMem     int    `json:"total_mem"`
	UsedMem      int    `json:"used_mem"`
	MonoUsedMem  int    `json:"mono_used_mem"`
	MonoStackMem int    `json:"mono_stack_mem"`
	Texture      int    `json:"texture"`
	Audio        int    `json:"audio"`
	TextAsset    int    `json:"text_asset"`
	Shader       int    `json:"shader"`
	Pic          string `json:"pic"`
	Process      string `json:"process"`
	CreatedAt    int64  `json:"created_at"`
}

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
	ID          uint   `json:"id" gorm:"primaryKey"`
	LoginID     int    `json:"login_id" gorm:"unique"`
	AppID       int    `json:"app_id"`
	Package     string `json:"package_name"`
	ProductName string `json:"product_name"`
	RoleName    string `json:"role_name"`
	Device      string `json:"device_name"`
	CPU         string `json:"system_cpu"`
	GPU         string `json:"graphics_divice"`
	Memory      string `json:"system_mem"`
	GPUMemory   string `json:"graphics_mem"`
	CreatedAt   int64  `json:"mtime"`
}

type StatsInfo struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	LoginID     int    `json:"login_id"`
	FPS         int    `json:"fps"`
	TotalMem    int    `json:"total_mem"`
	UsedMem     int    `json:"used_mem"`
	MonoUsedMem int    `json:"mono_used_mem"`
	MonoHeapMem int    `json:"mono_heap_mem"`
	Texture     int    `json:"texture"`
	Mesh        int    `json:"mesh"`
	Animation   int    `json:"animation"`
	Audio       int    `json:"audio"`
	Font        int    `json:"font"`
	TextAsset   int    `json:"text_asset"`
	Shader      int    `json:"shader"`
	Pic         string `json:"pic"`
	Process     string `json:"list"`
	CreatedAt   int64  `json:"mtime"`
}

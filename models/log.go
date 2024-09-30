package models

type Log struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	AppID      string `json:"app_id"`
	Package    string `json:"package"`
	RoleName   string `json:"role_name" gorm:"index"`
	Device     string `json:"device"`
	LogMessage string `json:"log_message" gorm:"index"`
	LogTime    int64  `json:"log_time"`
	LogType    string `json:"log_type"`
	LogStack   string `json:"log_stack"`
	CreateAt   int64  `json:"create_at"`
}

type StatsRecord struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	LoginID     int64  `json:"login_id" gorm:"uniqueIndex"`
	AppID       int    `json:"app_id"`
	Package     string `json:"package_name"`
	ProductName string `json:"product_name"`
	RoleName    string `json:"role_name"`
	Device      string `json:"device_name"`
	CPU         string `json:"system_cpu"`
	GPU         string `json:"graphics_divice"`
	Memory      int    `json:"system_mem"`
	GPUMemory   int    `json:"graphics_mem"`
	CreatedAt   int64  `json:"mtime" gorm:"index"`
}

type StatsInfo struct {
	ID          uint                     `json:"id" gorm:"primaryKey"`
	LoginID     int64                    `json:"login_id" gorm:"index"`
	FPS         int                      `json:"fps"`
	TotalMem    int                      `json:"total_mem"`
	UsedMem     int                      `json:"used_mem"`
	MonoUsedMem int                      `json:"mono_used_mem"`
	MonoHeapMem int                      `json:"mono_heap_mem"`
	Texture     int                      `json:"texture"`
	Mesh        int                      `json:"mesh"`
	Animation   int                      `json:"animation"`
	Audio       int                      `json:"audio"`
	Font        int                      `json:"font"`
	TextAsset   int                      `json:"text_asset"`
	Shader      int                      `json:"shader"`
	Pic         string                   `json:"pic"`
	Process2    []map[string]interface{} `json:"list" gorm:"-"`
	Process     string                   `json:"process"`
	CreatedAt   int64                    `json:"mtime" gorm:"index"`
}

type ProcessItem struct {
	Name string   `json:"name"`
	List []string `json:"list"`
}

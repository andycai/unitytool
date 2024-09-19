package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	AppID      string
	Package    string
	RoleName   string
	Device     string
	LogMessage string
	LogTime    int64
	LogType    string
	LogStack   string
}

type StatsRecord struct {
	ID        uint `gorm:"primaryKey"`
	LoginID   int  `gorm:"unique"`
	AppID     int
	Package   string
	RoleName  string
	Device    string
	CPU       string
	GPU       string
	Memory    string
	CreatedAt int64
}

type StatsInfo struct {
	ID           uint `gorm:"primaryKey"`
	LoginID      int
	FPS          int
	TotalMem     int
	UsedMem      int
	MonoUsedMem  int
	MonoStackMem int
	Texture      int
	Audio        int
	TextAsset    int
	Shader       int
	Pic          string
	Process      string
	CreatedAt    int64
}

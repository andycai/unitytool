package modules

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ModuleConfig 模块配置接口
type ModuleConfig interface {
	IsEnabled() bool
}

// Module 模块接口
type Module interface {
	Init() error
	RegisterRoutes(app *fiber.App)
}

// BaseModule 基础模块结构
type BaseModule struct {
	DB     *gorm.DB
	Config ModuleConfig
}

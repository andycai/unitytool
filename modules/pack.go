package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/andycai/unitool/handlers"
)

type PackModule struct {
	BaseModule
}

func (m *PackModule) Init() error {
	return nil
}

func (m *PackModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// Unity打包接口
	apiGroup.Post("/pack/ab", handlers.HandlePackAB)
	apiGroup.Post("/pack/apk", handlers.HandlePackAPK)
}

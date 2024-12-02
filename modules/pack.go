package modules

import (
	"github.com/gofiber/fiber/v2"
	"mind.com/log/handlers"
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
	app.Post("/pack/ab", handlers.HandlePackAB)
	app.Post("/pack/apk", handlers.HandlePackAPK)
}

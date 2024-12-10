package unibuild

import (
	"github.com/andycai/unitool/core"
	"github.com/gofiber/fiber/v2"
)

const (
	KeyModule        = "admin.unibuild"
	KeyNoCheckRouter = "admin.unibuild.router.nocheck"
)

func initModule() {
}

func initPublicRouter(publicGroup fiber.Router) {
	// Unity打包接口
	publicGroup.Post("/api/unibuild/ab", HandlePackAB)
	publicGroup.Post("/api/unibuild/apk", HandlePackAPK)
}

func init() {
	core.RegisterModule(KeyModule, initModule)
	core.RegisterPublicRouter(KeyNoCheckRouter, initPublicRouter)
}

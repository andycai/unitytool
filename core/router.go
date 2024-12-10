package core

import (
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRouter(app *fiber.App, db *gorm.DB) {
	public := app.Group("/")
	for _, f := range routerPublicNoCheckMap {
		f(public)
	}

	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware(db))
	for _, f := range routerAPICheckMap {
		f(api)
	}

	admin := app.Group("/admin")
	admin.Use(middleware.AuthMiddleware(db))
	for _, f := range routerAdminCheckMap {
		f(admin)
	}
}

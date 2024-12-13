package core

import (
	"github.com/andycai/unitool/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	App          *fiber.App
	DB           *gorm.DB
	DBs          []*gorm.DB
	Config       *Config
	RouterPublic fiber.Router
	RouterAdmin  fiber.Router
	RouterApi    fiber.Router
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(dbs []*gorm.DB, fiberApp *fiber.App) {
	a.Config = &config
	a.DBs = dbs
	a.DB = dbs[0]
	a.App = fiberApp

	a.RouterPublic = fiberApp.Group("/")
	a.RouterApi = fiberApp.Group("/api")
	a.RouterApi.Use(middleware.AuthMiddleware(dbs[0], config.Auth.JWTSecret))
	a.RouterAdmin = fiberApp.Group("/admin")
	a.RouterAdmin.Use(middleware.AuthMiddleware(dbs[0], config.Auth.JWTSecret))
}

package main

import (
	"mind.com/log/handlers"
	"mind.com/log/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("logs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Log{})

	app := fiber.New()

	// Serve static files
	app.Static("/", "./public")

	app.Post("/api/logs", func(c *fiber.Ctx) error {
		return handlers.CreateLog(c, db)
	})

	app.Get("/api/logs", func(c *fiber.Ctx) error {
		return handlers.GetLogs(c, db)
	})

	app.Delete("/api/logs", func(c *fiber.Ctx) error {
		return handlers.DeleteLogsBefore(c, db)
	})

	app.Listen(":3000")
}

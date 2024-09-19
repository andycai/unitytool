package main

import (
	"flag"
	"fmt"

	"mind.com/log/handlers"
	"mind.com/log/models"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	// 定义命令行参数
	port := flag.Int("port", 3000, "端口号")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open("db/logs.db"), &gorm.Config{})
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

	// 使用命令行参数设置端口
	app.Listen(fmt.Sprintf(":%d", *port))
}

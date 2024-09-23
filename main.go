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
	host := flag.String("host", "0.0.0.0", "主机地址")
	port := flag.Int("port", 3000, "端口号")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open("db/logs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Log{}, &models.StatsRecord{}, &models.StatsInfo{})

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

	app.Post("/api/stats", func(c *fiber.Ctx) error {
		return handlers.CreateStats(c, db)
	})

	app.Get("/api/stats", func(c *fiber.Ctx) error {
		return handlers.GetStats(c, db)
	})

	app.Delete("/api/stats", func(c *fiber.Ctx) error {
		return handlers.DeleteStatsBefore(c, db)
	})

	app.Get("/api/stats/details", func(c *fiber.Ctx) error {
		return handlers.GetStatDetails(c, db)
	})

	// 使用命令行参数设置端口
	app.Listen(fmt.Sprintf("%s:%d", *host, *port))
}

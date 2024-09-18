package handlers

import (
	"time"

	"mind.com/log/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateLog(c *fiber.Ctx, db *gorm.DB) error {
	log := new(models.Log)
	if err := c.BodyParser(log); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	result := db.Create(log)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create log"})
	}

	return c.Status(201).JSON(log)
}

func GetLogs(c *fiber.Ctx, db *gorm.DB) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")

	var logs []models.Log
	var total int64

	query := db.Model(&models.Log{})

	if search != "" {
		query = query.Where("log_message LIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit
	result := query.Offset(offset).Order("log_time DESC").Limit(limit).Find(&logs)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}

	return c.JSON(fiber.Map{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func DeleteLogsBefore(c *fiber.Ctx, db *gorm.DB) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Date parameter is required"})
	}

	// Parse the date string to time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format"})
	}

	// Set the time to the end of the day (23:59:59)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	// Delete logs before the end of the selected day
	result := db.Where("log_time < ?", endOfDay.Format("2006-01-02 15:04:05")).Delete(&models.Log{})
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete logs"})
	}

	return c.JSON(fiber.Map{"message": "Logs deleted successfully", "count": result.RowsAffected})
}

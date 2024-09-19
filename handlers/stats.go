package handlers

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mind.com/log/models"
)

func CreateStats(c *fiber.Ctx, db *gorm.DB) error {
	var stats models.StatsRecord
	if err := c.BodyParser(&stats); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	if err := db.Create(&stats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create stats record"})
	}
	return c.Status(fiber.StatusCreated).JSON(stats)
}

func GetStats(c *fiber.Ctx, db *gorm.DB) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	var stats []models.StatsRecord

	offset := (page - 1) * limit
	if err := db.Offset(offset).Limit(limit).Find(&stats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch stats records"})
	}
	return c.JSON(fiber.Map{
		"stats": stats,
		"total": len(stats),
		"page":  page,
		"limit": limit,
	})
}

func DeleteStatsBefore(c *fiber.Ctx, db *gorm.DB) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "date query parameter is required"})
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date format"})
	}

	// Convert date to milliseconds timestamp
	timestamp := date.UnixNano() / int64(time.Millisecond)

	// Delete from stats_record table
	if err := db.Where("created_at <= ?", timestamp).Delete(&models.StatsRecord{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete stats records"})
	}

	// Fetch stats_info records to delete
	var statsInfo []models.StatsInfo
	if err := db.Where("created_at <= ?", timestamp).Find(&statsInfo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch stats info records"})
	}

	// Delete images from file system
	for _, info := range statsInfo {
		if info.Pic != "" {
			if err := os.Remove(filepath.Join("public", info.Pic)); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete image file"})
			}
		}
	}

	// Delete from stats_info table
	if err := db.Where("created_at <= ?", timestamp).Delete(&models.StatsInfo{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete stats info records"})
	}

	return c.JSON(fiber.Map{"message": "records deleted successfully"})
}

func GetStatDetails(c *fiber.Ctx, db *gorm.DB) error {
	loginID := c.Query("login_id")
	if loginID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "login_id is required"})
	}

	var statsRecord models.StatsRecord
	if err := db.Where("login_id = ?", loginID).First(&statsRecord).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "stats record not found"})
	}

	var statsInfo []models.StatsInfo
	if err := db.Where("login_id = ?", loginID).Find(&statsInfo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch stats info"})
	}

	return c.JSON(fiber.Map{
		"statsRecord": statsRecord,
		"statsInfo":   statsInfo,
	})
}

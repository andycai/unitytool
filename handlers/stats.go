package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mind.com/log/models"
)

// 创建资源占用记录
func CreateStats(c *fiber.Ctx, db *gorm.DB) error {
	var record models.StatsRecord
	if err := c.BodyParser(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	var existingRecord models.StatsRecord
	if err := db.Where("login_id = ?", record.LoginID).First(&existingRecord).Error; err != nil {
		if err := db.Create(&record).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create stats record"})
		}
	}

	var info models.StatsInfo
	if err := c.BodyParser(&info); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	tmp, _ := json.Marshal(&info.Process2)
	info.Process = string(tmp)

	// Handle the pic field
	if info.Pic != "" {
		// Decode base64 to image data
		imgData, err := base64.StdEncoding.DecodeString(info.Pic)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid base64 image data"})
		}

		// Create directory if it doesn't exist
		uploadDir := "./public/uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create upload directory"})
		}

		// Generate a unique filename
		filename := time.Now().Format("20060102150405") + ".jpg"
		filePath := filepath.Join(uploadDir, filename)

		// Save the image file
		if err := os.WriteFile(filePath, imgData, 0644); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save image file"})
		}

		// Update the pic field with the file path
		info.Pic = filepath.Join("uploads", filename)
	}

	if err := db.Create(&info).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create stats record"})
	}

	return c.Status(fiber.StatusCreated).JSON(record)
}

// 获取资源占用记录
func GetStats(c *fiber.Ctx, db *gorm.DB) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	var stats []models.StatsRecord

	offset := (page - 1) * limit
	if err := db.Offset(offset).Order("created_at desc").Limit(limit).Find(&stats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch stats records"})
	}
	return c.JSON(fiber.Map{
		"stats": stats,
		"total": len(stats),
		"page":  page,
		"limit": limit,
	})
}

// 删除资源占用记录
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
		fmt.Println("err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete stats records"})
	}

	// Fetch stats_info records to delete
	var statsInfo []models.StatsInfo
	if err := db.Where("created_at <= ?", timestamp).Find(&statsInfo).Error; err != nil {
		fmt.Println("err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot fetch stats info records"})
	}

	// Delete images from file system
	for _, info := range statsInfo {
		picPath := filepath.Join("public", info.Pic)
		fmt.Println("picPath", picPath)
		if info.Pic != "" {
			if err := os.Remove(picPath); err != nil {
				fmt.Println("err", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete image file"})
			}
		}
	}

	// Delete from stats_info table
	result := db.Where("created_at <= ?", timestamp).Delete(&models.StatsInfo{})
	if result.Error != nil {
		fmt.Println("err", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot delete stats info records"})
	}

	return c.JSON(fiber.Map{"message": "records deleted successfully", "count": result.RowsAffected})
}

// 获取资源占用详情
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

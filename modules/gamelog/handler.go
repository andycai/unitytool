package gamelog

import (
	"fmt"
	"time"

	"github.com/andycai/unitool/models"
	"github.com/andycai/unitool/modules/adminlog"

	"github.com/gofiber/fiber/v2"
)

type LogReq struct {
	AppID    string           `json:"app_id"`
	Package  string           `json:"package"`
	RoleName string           `json:"role_name"`
	Device   string           `json:"device"`
	Logs     []models.GameLog `json:"list"`
}

// 创建日志记录
func createLog(c *fiber.Ctx) error {
	logReq := new(LogReq)
	if err := c.BodyParser(logReq); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Loop through LogReq.Logs and copy appid, package, role_name, device to each log
	for i := range logReq.Logs {
		logReq.Logs[i].AppID = logReq.AppID
		logReq.Logs[i].Package = logReq.Package
		logReq.Logs[i].RoleName = logReq.RoleName
		logReq.Logs[i].Device = logReq.Device
		logReq.Logs[i].CreateAt = time.Now().UnixMilli()
	}

	result := app.DB.CreateInBatches(logReq.Logs, len(logReq.Logs))
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create log"})
	}

	// 记录操作日志
	// adminlog.CreateAdminLog(c, "create", "gamelog", 0, fmt.Sprintf("批量创建游戏日志，角色：%s，数量：%d", logReq.RoleName, len(logReq.Logs)))

	return c.Status(201).JSON(logReq.Logs)
}

// 获取日志记录
func getLogs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)
	search := c.Query("search", "")

	var logs []models.GameLog
	var total int64

	query := app.DB.Model(&models.GameLog{})

	if search != "" {
		query = query.Where("role_name LIKE ? OR log_message LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	result := query.Offset(offset).Order("create_at DESC").Limit(pageSize).Find(&logs)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}

	return c.JSON(fiber.Map{
		"logs":     logs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// 删除日志记录
func deleteLogsBefore(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(400).JSON(fiber.Map{"code": 1, "error": "Date parameter is required"})
	}

	// Parse the date string to time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"code": 2, "error": "Invalid date format"})
	}

	// Set the time to the end of the day (23:59:59)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())
	endOfDayMilli := endOfDay.UnixMilli()

	// Delete logs before the end of the selected day
	result := app.DB.Where("create_at < ?", endOfDayMilli).Delete(&models.GameLog{})

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"code": 3, "error": "Failed to delete logs"})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "gamelog", 0, fmt.Sprintf("批量删除%s之前的游戏日志，共%d条", dateStr, result.RowsAffected))

	return c.JSON(fiber.Map{"code": 0, "message": "Logs deleted successfully", "count": result.RowsAffected})
}

// 删除单条日志记录
func deleteLog(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"code": 4, "error": "Log ID is required"})
	}

	// 获取日志信息用于记录
	var log models.GameLog
	if err := app.DB.First(&log, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"code": 6, "error": "Log not found"})
	}

	result := app.DB.Delete(&log)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"code": 5, "error": "Failed to delete log"})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"code": 6, "error": "Log not found"})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "gamelog", uint(log.ID), fmt.Sprintf("删除游戏日志，角色：%s", log.RoleName))

	return c.JSON(fiber.Map{"code": 0, "message": "Log deleted successfully"})
}

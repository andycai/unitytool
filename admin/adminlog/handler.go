package adminlog

import (
	"fmt"

	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

// CreateAdminLog 创建操作日志
func CreateAdminLog(c *fiber.Ctx, action string, resource string, resourceID uint, details string) error {
	user := c.Locals("user").(models.User)

	if user.ID == 0 {
		return fmt.Errorf("登录已过期，请重新登录")
	}

	log := models.AdminLog{
		UserID:     user.ID,
		Username:   user.Username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IP:         c.IP(),
		UserAgent:  c.Get("User-Agent"),
		CreatedAt:  db.NowFunc(),
	}

	if err := db.Create(&log).Error; err != nil {
		return err
	}

	return nil
}

// getAdminLogs 获取操作日志列表
func getAdminLogs(c *fiber.Ctx) error {
	var total int64

	// 获取查询参数
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	username := c.Query("username")
	action := c.Query("action")
	resource := c.Query("resource")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 构建查询
	query := db.Model(&models.AdminLog{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// 获取总数
	query.Count(&total)

	// 获取分页数据
	var logs []models.AdminLog
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "获取操作日志失败",
		})
	}

	return c.JSON(fiber.Map{
		"total": total,
		"data":  logs,
	})
}

// deleteAdminLogs 删除指定日期之前的操作日志
func deleteAdminLogs(c *fiber.Ctx) error {
	beforeDate := c.Query("beforeDate")
	if beforeDate == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "请指定日期",
		})
	}

	if err := db.Where("created_at < ?", beforeDate).Delete(&models.AdminLog{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "删除操作日志失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}

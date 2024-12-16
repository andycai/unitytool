package user

import (
	"github.com/andycai/unitool/lib/authentication"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

// GetByID 获取用户
func GetByID(id uint) *models.User {
	var vo models.User
	app.DB.Model(&vo).
		Where("id", id).
		First(&vo)

	return &vo
}

func Current(c *fiber.Ctx) *models.User {
	isAuthenticated, userID := authentication.AuthGet(c)

	if !isAuthenticated {
		return nil
	}

	return GetByID(userID)
}
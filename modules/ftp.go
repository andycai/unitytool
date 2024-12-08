package modules

import (
	"github.com/gofiber/fiber/v2"
	"github.com/andycai/unitool/handlers"
	"github.com/andycai/unitool/utils"
)

type FTPModule struct {
	BaseModule
	ServerConfig utils.ServerConfig
	FTPConfig    utils.FTPConfig
}

func (m *FTPModule) Init() error {
	// 初始化 FTP 配置
	handlers.InitFTP(handlers.FTPConfig{
		Host:       m.FTPConfig.Host,
		Port:       m.FTPConfig.Port,
		User:       m.FTPConfig.User,
		Password:   m.FTPConfig.Password,
		APKPath:    m.FTPConfig.APKPath,
		ZIPPath:    m.FTPConfig.ZIPPath,
		LogDir:     m.FTPConfig.LogDir,
		MaxLogSize: m.FTPConfig.MaxLogSize,
	})
	return nil
}

func (m *FTPModule) RegisterRoutes(app *fiber.App) {
	if !m.Config.IsEnabled() {
		return
	}

	// FTP 上传路由
	adminGroup.Post("/ftp/upload", func(c *fiber.Ctx) error {
		return handlers.HandleFTPUpload(c, m.ServerConfig.Output)
	})
}

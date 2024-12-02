package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"mind.com/log/modules"
	"mind.com/log/utils"
)

func main() {
	// 加载配置文件
	if err := utils.LoadConfig(); err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 初始化数据库
	db, err := modules.InitDatabase()
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建 Fiber 应用
	app := fiber.New()

	// 初始化并注册模块
	moduleList := []modules.Module{
		&modules.LogsModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("logs"),
			},
		},
		&modules.StatsModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("stats"),
			},
		},
		&modules.BrowseModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("browse"),
			},
			ServerConfig: utils.GetServerConfig(),
		},
		&modules.FTPModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("ftp"),
			},
			ServerConfig: utils.GetServerConfig(),
			FTPConfig:    utils.GetFTPConfig(),
		},
		&modules.ServerConfModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("serverconf"),
			},
			ServerConfig: utils.GetServerConfig(),
			JSONPaths:    utils.GetJSONPathConfig(),
		},
		&modules.CmdModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("cmd"),
			},
			ServerConfig: utils.GetServerConfig(),
		},
		&modules.PackModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("pack"),
			},
		},
	}

	// 初始化和注册所有模块
	for _, module := range moduleList {
		if err := module.Init(); err != nil {
			log.Printf("模块初始化失败: %v", err)
			continue
		}
		module.RegisterRoutes(app)
	}

	// 启动服务器
	serverConfig := utils.GetServerConfig()
	app.Listen(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}

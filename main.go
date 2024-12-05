package main

import (
	"fmt"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
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

	// 初始化模板引擎
	engine := html.New("./templates", ".html")
	engine.Reload(true) // 开发模式下启用模板重载
	engine.Debug(true)  // 开发模式下启用调试信息

	// 添加模板函数
	engine.AddFunc("yield", func() string { return "" })
	engine.AddFunc("partial", func(name string, data interface{}) template.HTML {
		return template.HTML("")
	})

	// 创建 Fiber 应用，并配置模板引擎
	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "admin/layout", // 设置默认布局
	})

	// 注册静态路由
	serverConfig := utils.GetServerConfig()
	for _, staticPath := range serverConfig.StaticPaths {
		app.Static(staticPath.Route, staticPath.Path)
	}

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
		&modules.AuthModule{
			BaseModule: modules.BaseModule{
				DB:     db,
				Config: utils.GetModuleConfig("auth"),
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
	app.Listen(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}

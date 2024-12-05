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

	// 初始化全局路由
	modules.InitGlobalRoutes(app, db)

	// 初始化模块
	modules.InitModules(app, db)

	// 启动服务器
	app.Listen(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}

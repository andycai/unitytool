package main

import (
	"fmt"
	"html/template"
	"log"
	"strings"

	_ "github.com/andycai/unitool/admin"
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/lib/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm"
)

func main() {
	// 加载配置文件
	if err := core.LoadConfig(); err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 初始化数据库
	db, err := database.InitDatabase(core.GetConfig().Database.DSN, core.GetConfig().Database.Driver)
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
	// 添加 hasSuffix 函数用于检查文件扩展名
	engine.AddFunc("hasSuffix", strings.HasSuffix)
	// 添加 splitPath 函数用于分割路径
	engine.AddFunc("splitPath", func(path string) []string {
		if path == "" {
			return []string{}
		}
		return strings.Split(path, "/")
	})
	// 添加 sub 函数用于数字减法
	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})

	// 创建 Fiber 应用，并配置模板引擎
	fiberApp := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// 根据错误码返回对应的错误页面
			switch code {
			case fiber.StatusForbidden:
				return c.Status(code).Render("errors/403", fiber.Map{}, "")
			case fiber.StatusNotFound:
				return c.Status(code).Render("errors/404", fiber.Map{}, "")
			default:
				return c.Status(code).Render("errors/500", fiber.Map{}, "")
			}
		},
	})

	app := core.NewApp()
	app.Start([]*gorm.DB{db}, fiberApp)

	// 启动服务器
	fiberApp.Listen(fmt.Sprintf("%s:%d", app.Config.Server.Host, app.Config.Server.Port))
}

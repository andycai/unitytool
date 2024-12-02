package main

import (
	"flag"
	"fmt"
	"log"

	"mind.com/log/handlers"
	"mind.com/log/models"
	"mind.com/log/utils"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func main() {
	// 定义命令行参数
	configPath := flag.String("config", "conf.toml", "配置文件路径")
	host := flag.String("host", "", "主机地址")
	port := flag.Int("port", 0, "端口号")
	output := flag.String("output", "", "输出目录")
	scriptPath := flag.String("script_path", "", "脚本路径")
	userDataPath := flag.String("user_data_path", "", "用户数据路径")
	flag.Parse()

	// 加载配置文件
	if err := utils.LoadConfig(*configPath); err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 获取配置
	serverConfig := utils.GetServerConfig()
	dbConfig := utils.GetDatabaseConfig()

	// 命令行参数覆盖配置文件
	if *host != "" {
		serverConfig.Host = *host
	}
	if *port != 0 {
		serverConfig.Port = *port
	}
	if *output != "" {
		serverConfig.Output = *output
	}
	if *scriptPath != "" {
		serverConfig.ScriptPath = *scriptPath
	}
	if *userDataPath != "" {
		serverConfig.UserDataPath = *userDataPath
	}

	// 更新配置
	utils.UpdateServerConfig(serverConfig)

	// 初始化数据库连接
	var db *gorm.DB
	var err error
	switch dbConfig.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbConfig.DSN), &gorm.Config{})
	default:
		log.Fatalf("不支持的数据库驱动: %s", dbConfig.Driver)
	}

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&models.Log{}, &models.StatsRecord{}, &models.StatsInfo{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	app := fiber.New()

	// Serve static files
	app.Static("/", serverConfig.StaticPath)

	// 处理目录浏览请求
	app.Get("/browse/*", func(c *fiber.Ctx) error {
		return handlers.HandleFileServer(c, serverConfig.Output)
	})

	// 添加文件删除路由
	app.Delete("/browse/*", func(c *fiber.Ctx) error {
		return handlers.HandleFileServer(c, serverConfig.Output)
	})

	// begin 脚本命令
	app.Post("/cmd", func(c *fiber.Ctx) error {
		return handlers.ExecShell(c, serverConfig.ScriptPath)
	})
	// end

	// begin Unity打包接口
	app.Post("/pack/ab", handlers.HandlePackAB)
	app.Post("/pack/apk", handlers.HandlePackAPK)
	// end

	// begin 日志接口
	app.Post("/api/logs", func(c *fiber.Ctx) error {
		return handlers.CreateLog(c, db)
	})

	app.Get("/api/logs", func(c *fiber.Ctx) error {
		return handlers.GetLogs(c, db)
	})

	app.Delete("/api/logs/before", func(c *fiber.Ctx) error {
		return handlers.DeleteLogsBefore(c, db)
	})

	app.Delete("/api/logs/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteLog(c, db)
	})
	// end 日志接口

	// begin 统计接口
	app.Post("/api/stats", func(c *fiber.Ctx) error {
		return handlers.CreateStats(c, db)
	})

	app.Get("/api/stats", func(c *fiber.Ctx) error {
		return handlers.GetStats(c, db)
	})

	app.Delete("/api/stats/before", func(c *fiber.Ctx) error {
		return handlers.DeleteStatsBefore(c, db)
	})

	app.Get("/api/stats/details", func(c *fiber.Ctx) error {
		return handlers.GetStatDetails(c, db)
	})

	app.Delete("/api/stats/:id", func(c *fiber.Ctx) error {
		return handlers.DeleteStat(c, db)
	})
	// end 统计接口

	// 添加 FTP 上传路由
	app.Post("/ftp/upload", func(c *fiber.Ctx) error {
		return handlers.HandleFTPUpload(c, serverConfig.Output)
	})

	username, password, err := utils.ReadFromBinaryFile(serverConfig.UserDataPath)
	if err != nil {
		log.Fatalf("读取用户数据文件失败: %v", err)
	}
	// 初始化 FTP 配置
	ftpConfig := utils.GetFTPConfig()
	handlers.InitFTP(handlers.FTPConfig{
		Host:       ftpConfig.Host,
		Port:       ftpConfig.Port,
		User:       username,
		Password:   password,
		APKPath:    ftpConfig.APKPath,
		ZIPPath:    ftpConfig.ZIPPath,
		LogDir:     ftpConfig.LogDir,
		MaxLogSize: ftpConfig.MaxLogSize,
	})

	// 初始化 JSON 路径配置
	jsonPaths := utils.GetJSONPathConfig()
	handlers.InitJSONPaths(jsonPaths)

	// begin 服务器配置接口
	app.Get("/server-config", func(c *fiber.Ctx) error {
		return c.SendFile("templates/server_config.html")
	})

	app.Get("/api/serverlist", handlers.GetServerList)
	app.Post("/api/serverlist", handlers.UpdateServerList)

	app.Get("/api/lastserver", handlers.GetLastServer)
	app.Post("/api/lastserver", handlers.UpdateLastServer)

	app.Get("/api/serverinfo", handlers.GetServerInfo)
	app.Post("/api/serverinfo", handlers.UpdateServerInfo)
	// end 服务器配置接口

	// 使用命令行参数设置端口
	app.Listen(fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port))
}

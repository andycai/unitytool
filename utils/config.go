package utils

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server    ServerConfig   `toml:"server"`
	Modules   ModulesConfig  `toml:"modules"`
	Database  DatabaseConfig `toml:"database"`
	JSONPaths JSONPathConfig `toml:"json_paths"`
	FTP       FTPConfig      `toml:"ftp"`
}

type ServerConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	Output       string `toml:"output"`
	ScriptPath   string `toml:"script_path"`
	StaticPath   string `toml:"static_path"`
	UserDataPath string `toml:"user_data_path"`
}

type DatabaseConfig struct {
	Driver string `toml:"driver"`
	DSN    string `toml:"dsn"`
}

type FTPConfig struct {
	Host       string `toml:"host"`
	Port       string `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	APKPath    string `toml:"apk_path"`
	ZIPPath    string `toml:"zip_path"`
	LogDir     string `toml:"log_dir"`
	MaxLogSize int64  `toml:"max_log_size"`
}

type JSONPathConfig struct {
	ServerList string `toml:"server_list"`
	LastServer string `toml:"last_server"`
	ServerInfo string `toml:"server_info"`
}

type ModulesConfig struct {
	Logs       bool `toml:"logs"`
	Stats      bool `toml:"stats"`
	Browse     bool `toml:"browse"`
	FTP        bool `toml:"ftp"`
	ServerConf bool `toml:"serverconf"`
	Cmd        bool `toml:"cmd"`
	Pack       bool `toml:"pack"`
}

var config Config

func LoadConfig() error {
	// 定义配置文件路径参数
	configPath := flag.String("config", "conf.toml", "配置文件路径")
	host := flag.String("host", "", "主机地址")
	port := flag.Int("port", 0, "端口号")
	output := flag.String("output", "", "输出目录")
	scriptPath := flag.String("script_path", "", "脚本路径")
	userDataPath := flag.String("user_data_path", "", "用户数据路径")
	dbPath := flag.String("db", "", "数据库路径")
	ftpHost := flag.String("ftp_host", "", "FTP主机地址")
	ftpPort := flag.String("ftp_port", "", "FTP端口")
	ftpUser := flag.String("ftp_user", "", "FTP用户名")
	ftpPass := flag.String("ftp_pass", "", "FTP密码")
	ftpApkPath := flag.String("ftp_apk_path", "", "FTP APK上传路径")
	ftpZipPath := flag.String("ftp_zip_path", "", "FTP ZIP上传路径")

	flag.Parse()

	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		return err
	}

	if *host != "" {
		config.Server.Host = *host
	}
	if *port != 0 {
		config.Server.Port = *port
	}
	if *output != "" {
		config.Server.Output = *output
	}
	if *scriptPath != "" {
		config.Server.ScriptPath = *scriptPath
	}
	if *userDataPath != "" {
		config.Server.UserDataPath = *userDataPath
	}
	if *dbPath != "" {
		config.Database.DSN = *dbPath
	}
	if *ftpHost != "" {
		config.FTP.Host = *ftpHost
	}
	if *ftpPort != "" {
		config.FTP.Port = *ftpPort
	}
	if *ftpUser != "" {
		config.FTP.User = *ftpUser
	}
	if *ftpPass != "" {
		config.FTP.Password = *ftpPass
	}
	if *ftpApkPath != "" {
		config.FTP.APKPath = *ftpApkPath
	}
	if *ftpZipPath != "" {
		config.FTP.ZIPPath = *ftpZipPath
	}

	return nil
}

func GetConfig() Config {
	return config
}

func GetServerConfig() ServerConfig {
	return config.Server
}

func GetDatabaseConfig() DatabaseConfig {
	return config.Database
}

func GetFTPConfig() FTPConfig {
	return config.FTP
}

func GetJSONPathConfig() JSONPathConfig {
	return config.JSONPaths
}

func UpdateServerConfig(newConfig ServerConfig) {
	config.Server = newConfig
}

type ModuleConfig struct {
	enabled bool
}

func (c ModuleConfig) IsEnabled() bool {
	return c.enabled
}

func GetModuleConfig(name string) ModuleConfig {
	switch name {
	case "logs":
		return ModuleConfig{enabled: config.Modules.Logs}
	case "stats":
		return ModuleConfig{enabled: config.Modules.Stats}
	case "browse":
		return ModuleConfig{enabled: config.Modules.Browse}
	case "ftp":
		return ModuleConfig{enabled: config.Modules.FTP}
	case "serverconf":
		return ModuleConfig{enabled: config.Modules.ServerConf}
	case "cmd":
		return ModuleConfig{enabled: config.Modules.Cmd}
	case "pack":
		return ModuleConfig{enabled: config.Modules.Pack}
	default:
		return ModuleConfig{enabled: false}
	}
}

func UpdateDatabaseConfig(newConfig DatabaseConfig) {
	config.Database = newConfig
}

func UpdateFTPConfig(newConfig FTPConfig) {
	config.FTP = newConfig
}

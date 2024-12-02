package utils

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server    ServerConfig   `toml:"server"`
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

var config Config

func LoadConfig(path string) error {
	_, err := toml.DecodeFile(path, &config)
	return err
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

package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"mind.com/log/utils"
)

// ServerList 结构体
type ServerList struct {
	ServerList []ServerListItem `json:"serverlist"`
}

type ServerListItem struct {
	ServerID     string `json:"server_id"`
	Name         string `json:"name"`
	ServerStatus string `json:"server_status"`
	Available    string `json:"available"`
	MergeID      string `json:"mergeid"`
	Online       string `json:"online"`
}

// LastServer 结构体
type LastServer struct {
	LastServer LastServerInfo `json:"lastserver"`
	Params     string         `json:"params"`
	SDKParams  string         `json:"sdkParams"`
}

type LastServerInfo struct {
	DefaultServer ServerInfo   `json:"deafult_server"`
	ListServer    []ServerInfo `json:"list_server"`
}

type ServerInfo struct {
	ServerID     string `json:"server_id"`
	Name         string `json:"name"`
	ServerStatus string `json:"server_status"`
}

// ServerInfo 结构体
type ServerInfoConfig struct {
	PFID          int    `json:"pfid"`
	PFName        string `json:"pfname"`
	Child         int    `json:"child"`
	AdKey         string `json:"adKey"`
	EntryURL      string `json:"entryURL"`
	CDNURL        string `json:"cdnURL"`
	CDNVersion    string `json:"cdnVersion"`
	LoginAPI      string `json:"loginAPI"`
	LoginURL      string `json:"loginURL"`
	ServerListURL string `json:"serverListURL"`
	Version       string `json:"version"`
	Time          string `json:"time"`
}

// 添加在文件开头的包级变量部分
var jsonPaths utils.JSONPathConfig

// 添加初始化函数
func InitJSONPaths(paths utils.JSONPathConfig) {
	jsonPaths = paths
}

// 读取 JSON 文件
func readJSONFile(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("解析 JSON 失败: %v", err)
	}

	return nil
}

// 写入 JSON 文件
func writeJSONFile(path string, v interface{}) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return fmt.Errorf("序列化 JSON 失败: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// 获取服务器列表
func GetServerList(c *fiber.Ctx) error {
	var serverList ServerList
	if err := readJSONFile(jsonPaths.ServerList, &serverList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(serverList)
}

// 更新服务器列表
func UpdateServerList(c *fiber.Ctx) error {
	var serverList ServerList
	if err := c.BodyParser(&serverList); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	if err := writeJSONFile(jsonPaths.ServerList, serverList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "更新成功"})
}

// 获取最后登录服务器
func GetLastServer(c *fiber.Ctx) error {
	var lastServer LastServer
	if err := readJSONFile(jsonPaths.LastServer, &lastServer); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(lastServer)
}

// 更新最后登录服务器
func UpdateLastServer(c *fiber.Ctx) error {
	var lastServer LastServer
	if err := c.BodyParser(&lastServer); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	if err := writeJSONFile(jsonPaths.LastServer, lastServer); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "更新成功"})
}

// 获取服务器信息
func GetServerInfo(c *fiber.Ctx) error {
	var serverInfo ServerInfoConfig
	if err := readJSONFile(jsonPaths.ServerInfo, &serverInfo); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(serverInfo)
}

// 更新服务器信息
func UpdateServerInfo(c *fiber.Ctx) error {
	var serverInfo ServerInfoConfig
	if err := c.BodyParser(&serverInfo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	if err := writeJSONFile(jsonPaths.ServerInfo, serverInfo); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "更新成功"})
}

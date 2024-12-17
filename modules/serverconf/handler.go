package serverconf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/andycai/unitool/modules/adminlog"
	"github.com/gofiber/fiber/v2"
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
	ServerPort   string `json:"server_port"`
	ServerIP     string `json:"server_ip"`
}

// LastServer 结构体
type LastServer struct {
	LastServer LastServerInfo `json:"lastserver"`
	Params     string         `json:"params"`
	SDKParams  string         `json:"sdkParams"`
}

type LastServerInfo struct {
	DefaultServer ServerInfo   `json:"default_server"`
	LastServer    []ServerInfo `json:"last_server"`
}

type ServerInfo struct {
	ServerID     string `json:"server_id"`
	Name         string `json:"name"`
	ServerStatus string `json:"server_status"`
	ServerPort   string `json:"server_port"`
	ServerIP     string `json:"server_ip"`
}

// ServerInfoConfig 结构体
type ServerInfoConfig struct {
	Fields []Field `json:"fields,omitempty"` // 字段列表，不保存到文件
}

// Field 字段结构
type Field struct {
	Key   string      `json:"key"`   // 字段名
	Value interface{} `json:"value"` // 字段值
	Type  string      `json:"type"`  // 字段类型
}

// NewServerInfoConfig 创建新的服务器配置
func NewServerInfoConfig() *ServerInfoConfig {
	return &ServerInfoConfig{
		Fields: make([]Field, 0),
	}
}

// MigrateFromOld 从旧配置迁移
func (c *ServerInfoConfig) MigrateFromOld(old map[string]interface{}) {
	// 处理所有字段
	for key, value := range old {
		if key != "fields" { // 跳过 fields 字段本身
			fieldType := "string"
			switch value.(type) {
			case float64:
				fieldType = "number"
			case bool:
				fieldType = "boolean"
			}
			c.Fields = append(c.Fields, Field{
				Key:   key,
				Value: value,
				Type:  fieldType,
			})
		}
	}
}

// ToMap 转换为map
func (c *ServerInfoConfig) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	// 添加所有字段
	for _, field := range c.Fields {
		var value interface{}
		switch field.Type {
		case "number":
			if v, ok := field.Value.(float64); ok {
				value = v
			} else if str, ok := field.Value.(string); ok {
				if v, err := strconv.ParseFloat(str, 64); err == nil {
					value = v
				}
			}
		case "boolean":
			if v, ok := field.Value.(bool); ok {
				value = v
			} else if str, ok := field.Value.(string); ok {
				value = str == "true"
			}
		default:
			value = fmt.Sprintf("%v", field.Value)
		}
		result[field.Key] = value
	}

	// 添加字段列表
	// result["fields"] = c.Fields

	return result
}

// 修改公告列表相关结构
type NoticeList []NoticeItem // 改为切片类型

type NoticeItem struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 添加公告数量相关结构
type NoticeNum struct {
	NoticeNum int `json:"noticenum"`
	Eject     int `json:"eject"`
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
func getServerList(c *fiber.Ctx) error {
	var serverList ServerList
	if err := readJSONFile(app.Config.JSONPaths.ServerList, &serverList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(serverList)
}

// 更新服务器列表
func updateServerList(c *fiber.Ctx) error {
	var serverList ServerList
	if err := c.BodyParser(&serverList); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	if err := writeJSONFile(app.Config.JSONPaths.ServerList, serverList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "serverconf_list", 0, "更新服务器信息配置")

	return c.JSON(fiber.Map{"message": "更新成功"})
}

// 获取最后登录服务器
func getLastServer(c *fiber.Ctx) error {
	var lastServer LastServer
	if err := readJSONFile(app.Config.JSONPaths.LastServer, &lastServer); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(lastServer)
}

// 更新最后登录服务器
func updateLastServer(c *fiber.Ctx) error {
	var lastServer LastServer
	if err := c.BodyParser(&lastServer); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	if err := writeJSONFile(app.Config.JSONPaths.LastServer, lastServer); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "serverconf_last", 0, "更新服务器信息配置")

	return c.JSON(fiber.Map{"message": "更新成功"})
}

// 获取服务器信息
func getServerInfo(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := readJSONFile(app.Config.JSONPaths.ServerInfo, &data); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	config := NewServerInfoConfig()
	config.MigrateFromOld(data)

	return c.JSON(config.ToMap())
}

// 更新服务器信息
func updateServerInfo(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	config := NewServerInfoConfig()

	// 处理字段列表
	if fields, ok := data["fields"].([]interface{}); ok {
		for _, field := range fields {
			if fieldMap, ok := field.(map[string]interface{}); ok {
				key, _ := fieldMap["key"].(string)
				value := fieldMap["value"]
				fieldType, _ := fieldMap["type"].(string)
				if key != "" {
					config.Fields = append(config.Fields, Field{
						Key:   key,
						Value: value,
						Type:  fieldType,
					})
				}
			}
		}
	}

	// 保存时只保存实际的字段值，不保存 fields 数组
	if err := writeJSONFile(app.Config.JSONPaths.ServerInfo, config.ToMap()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "serverconf_info", 0, "更新服务器信息配置")

	return c.JSON(fiber.Map{"message": "更新成功"})
}

// 获取公告列表
func getNoticeList(c *fiber.Ctx) error {
	var noticeList NoticeList
	if err := readJSONFile(app.Config.JSONPaths.NoticeList, &noticeList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(noticeList)
}

// 更新公告列表
func updateNoticeList(c *fiber.Ctx) error {
	var noticeList NoticeList
	if err := c.BodyParser(&noticeList); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	if err := writeJSONFile(app.Config.JSONPaths.NoticeList, noticeList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "serverconf_notice", 0, "更新公告列表配置")

	return c.JSON(fiber.Map{"message": "更新成功"})
}

// 获取公告数量
func getNoticeNum(c *fiber.Ctx) error {
	var noticeNum NoticeNum
	if err := readJSONFile(app.Config.JSONPaths.NoticeNum, &noticeNum); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(noticeNum)
}

// 更新公告数量
func updateNoticeNum(c *fiber.Ctx) error {
	var noticeNum NoticeNum
	if err := c.BodyParser(&noticeNum); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 确保 eject 是 int 类型
	if noticeNum.Eject < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "eject 必须是非负整数"})
	}

	if err := writeJSONFile(app.Config.JSONPaths.NoticeNum, noticeNum); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "serverconf_notice_num", 0, "更新公告数量配置")

	return c.JSON(fiber.Map{"message": "更新成功"})
}

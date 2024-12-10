package citask

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/andycai/unitool/admin/adminlog"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

// getTasks 获取任务列表
func getTasks(c *fiber.Ctx) error {
	var tasks []models.Task
	if err := db.Order("created_at desc").Find(&tasks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取任务列表失败: %v", err),
		})
	}
	return c.JSON(tasks)
}

// createTask 创建任务
func createTask(c *fiber.Ctx) error {
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	if err := db.Create(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建任务失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "create", "task", task.ID, fmt.Sprintf("创建任务：%s", task.Name))

	return c.JSON(task)
}

// getTask 获取任务详情
func getTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}
	return c.JSON(task)
}

// updateTask 更新任务
func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}

	// 解析请求体
	updates := new(models.Task)
	if err := c.BodyParser(updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 更新字段
	task.Name = updates.Name
	task.Description = updates.Description
	task.Type = updates.Type
	task.Script = updates.Script
	task.URL = updates.URL
	task.Method = updates.Method
	task.Headers = updates.Headers
	task.Body = updates.Body
	task.Timeout = updates.Timeout
	task.Status = updates.Status

	if err := db.Save(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("更新任务失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "update", "task", task.ID, fmt.Sprintf("更新任务：%s", task.Name))

	return c.JSON(task)
}

// deleteTask 删除任务
func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}

	if err := db.Delete(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除任务失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "task", task.ID, fmt.Sprintf("删除任务：%s", task.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// runTask 执行任务
func runTask(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}

	// 创建任务日志
	taskLog := models.TaskLog{
		TaskID:    task.ID,
		Status:    "running",
		StartTime: time.Now(),
	}
	if err := db.Create(&taskLog).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建任务日志失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "run", "task", task.ID, fmt.Sprintf("执行任务：%s", task.Name))

	// 异步执行任务
	go executeTask(&task, &taskLog)

	return c.JSON(taskLog)
}

// getTaskLogs 获取任务日志
func getTaskLogs(c *fiber.Ctx) error {
	taskID := c.Params("id")
	var logs []models.TaskLog
	if err := db.Where("task_id = ?", taskID).Order("created_at desc").Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取任务日志失败: %v", err),
		})
	}
	return c.JSON(logs)
}

// getTaskStatus 获取任务状态
func getTaskStatus(c *fiber.Ctx) error {
	logID := c.Query("log_id")
	var log models.TaskLog
	if err := db.First(&log, logID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("日志不存在: %v", err),
		})
	}
	return c.JSON(log)
}

// executeTask 执行任务
func executeTask(task *models.Task, log *models.TaskLog) {
	defer func() {
		log.EndTime = time.Now()
		log.Duration = int(log.EndTime.Sub(log.StartTime).Seconds())
		db.Save(log)
	}()

	switch task.Type {
	case "script":
		executeScriptTask(task, log)
	case "http":
		executeHTTPTask(task, log)
	default:
		log.Status = "failed"
		log.Error = "未知的任务类型"
	}
}

// executeScriptTask 执行脚本任务
func executeScriptTask(task *models.Task, log *models.TaskLog) {
	// 创建临时脚本文件
	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".bat"
	}

	tmpFile, err := os.CreateTemp("", "task_*"+ext)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建临时文件失败: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(task.Script); err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("写入脚本失败: %v", err)
		return
	}
	tmpFile.Close()

	// 设置脚本可执行权限
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("设置脚本权限失败: %v", err)
			return
		}
	}

	// 执行脚本
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", tmpFile.Name())
	} else {
		cmd = exec.Command("/bin/bash", tmpFile.Name())
	}

	// 设置工作目录
	cmd.Dir = filepath.Dir(tmpFile.Name())

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("执行失败: %v\n%s", err, string(output))
		return
	}

	log.Status = "success"
	log.Output = string(output)
}

// executeHTTPTask 执行HTTP任务
func executeHTTPTask(task *models.Task, log *models.TaskLog) {
	// 解析请求头
	headers := make(map[string]string)
	if task.Headers != "" {
		if err := json.Unmarshal([]byte(task.Headers), &headers); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("解析请求头失败: %v", err)
			return
		}
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(task.Timeout) * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest(task.Method, task.URL, strings.NewReader(task.Body))
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建请求失败: %v", err)
		return
	}

	// 添加请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("发送请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("读取响应失败: %v", err)
		return
	}

	log.Output = string(body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Status = "success"
	} else {
		log.Status = "failed"
		log.Error = fmt.Sprintf("请求返回错误状态码: %d", resp.StatusCode)
	}
}

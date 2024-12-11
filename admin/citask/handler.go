package citask

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andycai/unitool/admin/adminlog"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

// TaskProgress 任务进度
type TaskProgress struct {
	ID        uint      `json:"id"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  int       `json:"duration"`
	Progress  int       `json:"progress"` // 0-100
}

var (
	taskProgressMap = make(map[uint]*TaskProgress)
	progressMutex   sync.RWMutex
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

	// 初始化进度信息
	progress := &TaskProgress{
		ID:        taskLog.ID,
		Status:    "running",
		StartTime: taskLog.StartTime,
		Progress:  0,
	}
	progressMutex.Lock()
	taskProgressMap[taskLog.ID] = progress
	progressMutex.Unlock()

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

// getTaskProgress 获取任务进度
func getTaskProgress(c *fiber.Ctx) error {
	logId := c.Params("logId")
	if logId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少logId参数",
		})
	}

	// 将字符串转换为uint
	id, err := strconv.ParseUint(logId, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的logId参数",
		})
	}

	progressMutex.RLock()
	progress, exists := taskProgressMap[uint(id)]
	progressMutex.RUnlock()

	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "找不到任务进度信息",
		})
	}

	return c.JSON(progress)
}

// executeTask 执行任务
func executeTask(task *models.Task, log *models.TaskLog) {
	progress := taskProgressMap[log.ID]
	defer func() {
		log.EndTime = time.Now()
		log.Duration = int(log.EndTime.Sub(log.StartTime).Seconds())
		db.Save(log)

		// 更新并清理进度信息
		if progress != nil {
			progress.Status = log.Status
			progress.EndTime = log.EndTime
			progress.Duration = log.Duration
			progress.Progress = 100

			// 延迟删除进度信息
			time.AfterFunc(time.Hour, func() {
				progressMutex.Lock()
				delete(taskProgressMap, log.ID)
				progressMutex.Unlock()
			})
		}
	}()

	switch task.Type {
	case "script":
		executeScriptTask(task, log, progress)
	case "http":
		executeHTTPTask(task, log, progress)
	default:
		log.Status = "failed"
		log.Error = "未知的任务类型"
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
	}
}

// isUnsafeCommand 检查命令是否不安全
func isUnsafeCommand(script string) (bool, string) {
	// 转换为小写以进行大小写不敏感的检查
	lowerScript := strings.ToLower(script)

	// 定义不安全的命令列表
	unsafeCommands := []string{
		"rm -rf", "rm -r", // 递归删除
		"mkfs",                         // 格式化
		":(){:|:&};:", ":(){ :|:& };:", // Fork炸弹
		"dd",                // 磁盘操作
		"> /dev/", ">/dev/", // 设备文件操作
		"wget", "curl", // 外部下载
		"chmod 777", "chmod -R 777", // 危险的权限设置
		"sudo", "su", // 提权命令
		"nc", "netcat", // 网络工具
		"telnet",          // 远程连接
		"|mail", "|email", // 邮件命令
		"tcpdump",        // 网络抓包
		"chown -R",       // 递归改变所有者
		"mv /* ", "mv /", // 移动根目录
		"cp /* ", "cp /", // 复制根目录
		"shutdown", "reboot", "halt", "poweroff", // 系统关机重启
		"passwd",             // 修改密码
		"useradd", "userdel", // 用户管理
		"mkfs", "fdisk", "fsck", // 磁盘管理
		"iptables", "firewall", // 防火墙
		"nmap", // 端口扫描
		"eval", // 命令注入
	}

	// 检查危险的shell特殊字符和重定向
	dangerousPatterns := []string{
		"$(", "`", // 命令替换
		"&&", "||", // 命令链接
		"../",       // 目录遍历
		"/*",        // 根目录操作
		"> /", ">/", // 重定向到系统目录
		"2> /", "2>/", // 错误重定向到系统目录
		">> /", ">>/", // 追加重定向到系统目录
		"< /", "</", // 从系统目录读取
	}

	// 检查不安全的命令
	for _, cmd := range unsafeCommands {
		if strings.Contains(lowerScript, cmd) {
			return true, fmt.Sprintf("检测到不安全的命令: %s", cmd)
		}
	}

	// 检查危险的模式
	for _, pattern := range dangerousPatterns {
		if strings.Contains(script, pattern) {
			return true, fmt.Sprintf("检测到危险的命令模式: %s", pattern)
		}
	}

	// 检查环境变量操作
	if strings.Contains(lowerScript, "export") || strings.Contains(lowerScript, "env") {
		return true, "不允许修改环境变量"
	}

	return false, ""
}

// executeScriptTask 执行脚本任务
func executeScriptTask(task *models.Task, log *models.TaskLog, progress *TaskProgress) {
	// 首先检查脚本安全性
	if unsafe, reason := isUnsafeCommand(task.Script); unsafe {
		log.Status = "failed"
		log.Error = fmt.Sprintf("脚本包含不安全的命令: %s", reason)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}

	// 创建临时脚本文件
	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".bat"
	}

	tmpFile, err := os.CreateTemp("", "task_*"+ext)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建临时文件失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}
	defer os.Remove(tmpFile.Name())

	// 添加安全限制的shell选项（仅用于Unix系统）
	scriptContent := task.Script
	if runtime.GOOS != "windows" {
		scriptContent = "set -euo pipefail\ntrap 'exit 1' INT TERM\n" + scriptContent
	}

	if _, err := tmpFile.WriteString(scriptContent); err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("写入脚本失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}
	tmpFile.Close()

	// 设置脚本可执行权限
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("设置脚本权限失败: %v", err)
			if progress != nil {
				progress.Status = "failed"
				progress.Error = log.Error
			}
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

	// 设置工作目录为临时目录
	tmpDir := os.TempDir()
	cmd.Dir = tmpDir

	// 创建输出缓冲区
	var outputBuffer bytes.Buffer
	outputWriter := io.MultiWriter(&outputBuffer)

	// 创建管道用于实时读取输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建输出管道失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建错误输出管道失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}

	// 设置超时
	timeout := time.Duration(task.Timeout) * time.Second
	if timeout == 0 {
		timeout = 300 * time.Second // 默认5分钟超时
	}

	// 创建一个带有超时的context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("启动命令失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}

	// 实时读取输出
	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			line := scanner.Text() + "\n"
			outputWriter.Write([]byte(line))
			if progress != nil {
				progress.Output = outputBuffer.String()
			}
		}
	}()

	// 等待命令完成
	err = cmd.Wait()
	output := outputBuffer.String()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Status = "failed"
			log.Error = fmt.Sprintf("执行超时（%d秒）\n%s", task.Timeout, output)
		} else {
			log.Status = "failed"
			log.Error = fmt.Sprintf("执行失败: %v\n%s", err, output)
		}
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}

	log.Status = "success"
	log.Output = output
	if progress != nil {
		progress.Status = "success"
		progress.Output = output
		progress.Progress = 100
	}
}

// executeHTTPTask 执行HTTP任务
func executeHTTPTask(task *models.Task, log *models.TaskLog, progress *TaskProgress) {
	// 解析请求头
	headers := make(map[string]string)
	if task.Headers != "" {
		if err := json.Unmarshal([]byte(task.Headers), &headers); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("解析请求头失败: %v", err)
			if progress != nil {
				progress.Status = "failed"
				progress.Error = log.Error
			}
			return
		}
	}

	// 创建请求
	req, err := http.NewRequest(task.Method, task.URL, strings.NewReader(task.Body))
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建请求失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 设置超时
	timeout := time.Duration(task.Timeout) * time.Second
	if timeout == 0 {
		timeout = 300 * time.Second // 默认5分钟超时
	}
	client := &http.Client{
		Timeout: timeout,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("发送请求失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("���取响应失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		log.Status = "failed"
		log.Error = fmt.Sprintf("请求失败: HTTP %d\n%s", resp.StatusCode, string(body))
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		return
	}

	log.Status = "success"
	log.Output = string(body)
	if progress != nil {
		progress.Status = "success"
		progress.Output = log.Output
		progress.Progress = 100
	}
}

// stopTask 停止任务
func stopTask(c *fiber.Ctx) error {
	logId := c.Params("logId")
	if logId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少logId参数",
		})
	}

	// 将字符串转换为uint
	id, err := strconv.ParseUint(logId, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的logId参数",
		})
	}

	progressMutex.RLock()
	progress, exists := taskProgressMap[uint(id)]
	progressMutex.RUnlock()

	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "找不到任务进度信息",
		})
	}

	if progress.Status != "running" {
		return c.Status(400).JSON(fiber.Map{
			"error": "任务已经结束",
		})
	}

	// TODO: 实现任务停止逻辑

	return c.JSON(fiber.Map{
		"message": "任务已停止",
	})
}

package citask

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/andycai/unitool/admin/adminlog"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
)

// TaskProgress 任务进度
type TaskProgress struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	TaskName  string    `json:"task_name"`
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
	taskCmdMap      = make(map[uint]*exec.Cmd)
	cronEntries     = make(map[uint]cron.EntryID)
	progressMutex   sync.RWMutex
	cronScheduler   *cron.Cron
)

func initCron() {
	// 初始化定时任务调度器
	cronScheduler = cron.New()
	cronScheduler.Start()

	// 从数据库加载定时任务
	var tasks []models.Task
	if err := db.Where("enable_cron = ? AND status = ?", true, "active").Find(&tasks).Error; err != nil {
		fmt.Printf("加载定时任务失败: %v\n", err)
		return
	}

	for _, task := range tasks {
		if err := scheduleCronTask(&task); err != nil {
			fmt.Printf("调度任务失败 [%d]: %v\n", task.ID, err)
		} else {
			fmt.Printf("成功加载定时任务 [%d]: %s\n", task.ID, task.Name)
		}
	}
}

// 加载定时任务
func loadCronTasks() {
	var tasks []models.Task
	if err := db.Where("enable_cron = ? AND status = ?", 1, "active").Find(&tasks).Error; err != nil {
		fmt.Printf("加载定时任务失败: %v\n", err)
		return
	}

	for _, task := range tasks {
		if err := scheduleCronTask(&task); err != nil {
			fmt.Printf("调度任务失败 [%d]: %v\n", task.ID, err)
		}
	}
}

// 调度定时任务
func scheduleCronTask(task *models.Task) error {
	if task.EnableCron == 0 || task.CronExpr == "" {
		return nil
	}

	entryID, err := cronScheduler.AddFunc(task.CronExpr, func() {
		// 创建任务日志
		taskLog := &models.TaskLog{
			TaskID:    task.ID,
			StartTime: time.Now(),
			Status:    "running",
		}

		if err := db.Create(taskLog).Error; err != nil {
			fmt.Printf("创建任务日志失败: %v\n", err)
			return
		}

		// 创建进度信息
		progress := &TaskProgress{
			Status:    "running",
			StartTime: time.Now(),
		}

		progressMutex.Lock()
		taskProgressMap[taskLog.ID] = progress
		progressMutex.Unlock()

		// 执行任务
		go func() {
			if task.Type == "script" {
				executeScriptTask(task, taskLog, progress)
			} else {
				executeHTTPTask(task, taskLog, progress)
			}

			// 更新任务日志
			taskLog.EndTime = time.Now()
			taskLog.Status = progress.Status
			taskLog.Output = progress.Output
			taskLog.Error = progress.Error
			if err := db.Save(taskLog).Error; err != nil {
				fmt.Printf("更新任务日志失败: %v\n", err)
			}
		}()
	})

	if err != nil {
		return err
	}

	// 保存定时任务ID
	progressMutex.Lock()
	cronEntries[task.ID] = entryID
	progressMutex.Unlock()

	return nil
}

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
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 如果启用了定时执行，验证cron表达式
	if task.EnableCron == 1 {
		if task.CronExpr == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "启用定时执行时必须提供Cron表达式",
			})
		}
		if _, err := cron.ParseStandard(task.CronExpr); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("无效的Cron表达式: %v", err),
			})
		}
	}

	if err := db.Create(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建任务失败: %v", err),
		})
	}

	// 如果启用了定时执行，添加到调度器
	if task.EnableCron == 1 {
		if err := scheduleCronTask(&task); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("设置定时任务失败: %v", err),
			})
		}
	}

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
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少任务ID",
		})
	}

	var updates models.Task
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 获取原有任务信息
	var task models.Task
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}

	// 检查定时任务状态变化
	cronChanged := task.EnableCron != updates.EnableCron ||
		(updates.EnableCron == 1 && task.CronExpr != updates.CronExpr)

	// 如果启用了定时执行，验证cron表达式
	if updates.EnableCron == 1 {
		if updates.CronExpr == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "启用定时执行时必须提供Cron表达式",
			})
		}
		if _, err := cron.ParseStandard(updates.CronExpr); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("无效的Cron表达式: %v", err),
			})
		}
	}

	// 更新任务信息
	if err := db.Model(&task).Updates(updates).Update("enable_cron", updates.EnableCron).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("更新任务失败: %v", err),
		})
	}

	// 如果定时配置发生变化，处理定时任务
	if cronChanged {
		progressMutex.Lock()
		// 如果存在旧的定时任务，先移除
		if entryID, ok := cronEntries[task.ID]; ok {
			cronScheduler.Remove(entryID)
			delete(cronEntries, task.ID)
			fmt.Printf("已移除任务 [%d] 的定时配置\n", task.ID)
		}
		progressMutex.Unlock()

		// 如果启用了定时执行，添加新的定时任务
		if updates.EnableCron == 1 {
			if err := scheduleCronTask(&task); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": fmt.Sprintf("更新定时任务失败: %v", err),
				})
			}
			fmt.Printf("已为任务 [%d] 添加新的定时配置: %s\n", task.ID, updates.CronExpr)
		}
	}

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": task,
	})
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
		TaskID:    task.ID,
		TaskName:  task.Name,
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

		// 执行完成任务，保存任务日志到数据库
		db.Save(log)

		// 更新并清理进度信息
		if progress != nil {
			progress.Status = log.Status
			progress.EndTime = log.EndTime
			progress.Duration = log.Duration
			progress.Progress = 100

			// 延迟删除进度信息
			time.AfterFunc(time.Hour*2, func() {
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
	fmt.Printf("开始执行脚本任务: %s (ID: %d)\n", task.Name, task.ID)

	// 首先检查脚本安全性
	if unsafe, reason := isUnsafeCommand(task.Script); unsafe {
		log.Status = "failed"
		log.Error = fmt.Sprintf("脚本包含不安全的命令: %s", reason)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("脚本安全检查失败: %s\n", reason)
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
		fmt.Printf("创建临时文件失败: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	fmt.Printf("创建临时脚本文件: %s\n", tmpFile.Name())

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
		fmt.Printf("写入脚本内容失败: %v\n", err)
		return
	}
	tmpFile.Close()
	fmt.Printf("写入脚本内容:\n%s\n", scriptContent)

	// 设置脚本可执行权限
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("设置脚本权限失败: %v", err)
			if progress != nil {
				progress.Status = "failed"
				progress.Error = log.Error
			}
			fmt.Printf("设置脚本权限失败: %v\n", err)
			return
		}
		fmt.Println("设置脚本可执行权限成功")
	}

	// 执行脚本
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", tmpFile.Name())
		fmt.Printf("Windows 命令: cmd /c %s\n", tmpFile.Name())
	} else {
		cmd = exec.Command("/bin/bash", tmpFile.Name())
		fmt.Printf("Unix 命令: /bin/bash %s\n", tmpFile.Name())
	}

	// 设置工作目录为临时目录
	tmpDir := os.TempDir()
	cmd.Dir = tmpDir
	fmt.Printf("工作目录: %s\n", tmpDir)

	// 创建输出缓冲区
	var outputBuffer, errorBuffer bytes.Buffer

	// 设置命令的标准输出和错误输出
	cmd.Stdout = io.MultiWriter(&outputBuffer, os.Stdout)
	cmd.Stderr = io.MultiWriter(&errorBuffer, os.Stderr)

	// 设置超时
	timeout := time.Duration(task.Timeout) * time.Second
	if timeout == 0 {
		timeout = 300 * time.Second // 默认5分钟超时
	}
	fmt.Printf("设置超时时间: %v\n", timeout)

	// 创建一个带有超时的context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 使用context创建新的命令
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	cmd.Dir = tmpDir
	cmd.Stdout = io.MultiWriter(&outputBuffer, os.Stdout)
	cmd.Stderr = io.MultiWriter(&errorBuffer, os.Stderr)

	// 保存命令到映射中
	progressMutex.Lock()
	taskCmdMap[log.ID] = cmd
	progressMutex.Unlock()

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("启动命令失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("启动命令失败: %v\n", err)
		return
	}
	fmt.Println("命令启动成功")

	// 等待命令完成
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	// 清理函数
	cleanup := func() {
		progressMutex.Lock()
		delete(taskCmdMap, task.ID)
		progressMutex.Unlock()
	}
	defer cleanup()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				log.Status = "failed"
				log.Error = fmt.Sprintf("执行超时（%d秒）\n%s\n%s", task.Timeout, outputBuffer.String(), errorBuffer.String())
				if progress != nil {
					progress.Status = "failed"
					progress.Error = log.Error
					progress.Output = outputBuffer.String() + "\nError: " + errorBuffer.String()
				}
				fmt.Printf("命令执行超时: %v\n", ctx.Err())
				return
			}
		case err := <-done:
			output := outputBuffer.String()
			errorOutput := errorBuffer.String()

			fmt.Printf("命令执行完成，输出:\n%s\n", output)
			if errorOutput != "" {
				fmt.Printf("错误输出:\n%s\n", errorOutput)
			}

			if err != nil {
				log.Status = "failed"
				log.Error = fmt.Sprintf("执行失败: %v\n%s\n%s", err, output, errorOutput)
				if progress != nil {
					progress.Status = "failed"
					progress.Error = log.Error
					progress.Output = output + "\nError: " + errorOutput
				}
				fmt.Printf("命令执行失败: %v\n", err)
				return
			}

			log.Status = "success"
			log.Output = output
			if errorOutput != "" {
				log.Output += "\nError: " + errorOutput
			}

			if progress != nil {
				progress.Status = "success"
				progress.Output = log.Output
				progress.Progress = 100
			}
			fmt.Printf("任务执行成功完成: %s (ID: %d)\n", task.Name, task.ID)
			return
		case <-ticker.C:
			if progress != nil {
				progress.Output = outputBuffer.String()
				if errorBuffer.Len() > 0 {
					progress.Output += "\nError: " + errorBuffer.String()
				}
			}
		}
	}
}

// executeHTTPTask 执行HTTP任务
func executeHTTPTask(task *models.Task, log *models.TaskLog, progress *TaskProgress) {
	fmt.Printf("开始执行HTTP任务: %s (ID: %d)\n", task.Name, task.ID)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(task.Timeout) * time.Second,
	}

	// 创建请求
	var body io.Reader
	if task.Body != "" {
		body = strings.NewReader(task.Body)
	}
	req, err := http.NewRequest(task.Method, task.URL, body)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建请求失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("创建HTTP请求失败: %v\n", err)
		return
	}

	// 添加请求头
	if task.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(task.Headers), &headers); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("解析请求头失败: %v", err)
			if progress != nil {
				progress.Status = "failed"
				progress.Error = log.Error
			}
			fmt.Printf("解析请求头失败: %v\n", err)
			return
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
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
		fmt.Printf("发送HTTP请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("读取响应失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("读取HTTP响应失败: %v\n", err)
		return
	}

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		log.Status = "failed"
		log.Error = fmt.Sprintf("HTTP请求失败: %s\n响应内容: %s", resp.Status, string(respBody))
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("HTTP请求返回错误状态码: %d\n", resp.StatusCode)
		return
	}

	// 更新任务状态
	log.Status = "success"
	log.Output = string(respBody)
	if progress != nil {
		progress.Status = "success"
		progress.Output = log.Output
		progress.Progress = 100
	}
	fmt.Printf("HTTP任务执行成功完成: %s (ID: %d)\n", task.Name, task.ID)
}

// stopTask 停止正在执行的任务
func stopTask(c *fiber.Ctx) error {
	logId := c.Params("logId")
	if logId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少日志ID",
		})
	}

	// 将字符串ID转换为uint
	id, err := strconv.ParseUint(logId, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的日志ID",
		})
	}
	taskId := uint(id)

	// 获取进度信息
	progressMutex.Lock()
	progress, exists := taskProgressMap[taskId]
	progressMutex.Unlock()

	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "任务不存在或已结束",
		})
	}

	// 如果任务不是运行状态，返回错误
	if progress.Status != "running" {
		return c.Status(400).JSON(fiber.Map{
			"error": "任务不在运行状态",
		})
	}

	// 获取命令进程
	progressMutex.Lock()
	cmd := taskCmdMap[taskId]
	progressMutex.Unlock()

	if cmd == nil || cmd.Process == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "无法获取任务进程",
		})
	}

	// 停止进程
	if runtime.GOOS == "windows" {
		// Windows 下使用 taskkill 强制结束进程树
		exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
	} else {
		// Unix 系统下发送 SIGTERM 信号
		err = cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			// 如果 SIGTERM 失败，尝试 SIGKILL
			err = cmd.Process.Kill()
		}
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("停止任务失败: %v", err),
		})
	}

	// 更新任务状态
	progressMutex.Lock()
	progress.Status = "failed"
	progress.Error = "任务被手动停止"
	progress.EndTime = time.Now()
	progressMutex.Unlock()

	// 清理命令映射
	progressMutex.Lock()
	delete(taskCmdMap, taskId)
	progressMutex.Unlock()

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "任务已停止",
	})
}

// GetRunningTasks 获取正在执行的任务列表
func GetRunningTasks(c *fiber.Ctx) error {
	progressMutex.Lock()
	defer progressMutex.Unlock()

	// 从内存中获取所有正在执行的任务
	var runningTasks []fiber.Map
	for id, progress := range taskProgressMap {
		if progress.Status == "running" {
			// 查询任务信息
			var taskLog models.TaskLog
			if err := db.First(&taskLog, id).Error; err != nil {
				continue
			}

			// 构建返回数据
			runningTasks = append(runningTasks, fiber.Map{
				"id":         id,
				"name":       progress.TaskName,
				"status":     progress.Status,
				"progress":   progress.Progress,
				"output":     progress.Output,
				"error":      progress.Error,
				"start_time": progress.StartTime.Unix(),
			})
		}
	}

	// 按开始时间倒序排序
	sort.Slice(runningTasks, func(i, j int) bool {
		return runningTasks[i]["start_time"].(int64) > runningTasks[j]["start_time"].(int64)
	})

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": runningTasks,
	})
}

// getNextRunTime 计算下次执行时间
func getNextRunTime(c *fiber.Ctx) error {
	expr := c.Query("expr")
	if expr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少cron表达式",
		})
	}

	// 解析cron表达式
	schedule, err := cron.ParseStandard(expr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的cron表达式: %v", err),
		})
	}

	// 计算下次执行时间
	nextTime := schedule.Next(time.Now())

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": fiber.Map{
			"next_run":      nextTime.Unix(),
			"next_run_text": nextTime.Format("2006-01-02 15:04:05"),
		},
	})
}

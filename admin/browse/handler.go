package browse

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/andycai/unitool/admin/adminlog"
	"github.com/andycai/unitool/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jlaffaye/ftp"
)

var ftpConfig utils.FTPConfig

// initFTP 初始化 FTP 配置
func initFTP(config utils.FTPConfig) {
	ftpConfig = config
}

// FileEntry 存储文件信息的结构体
type FileEntry struct {
	Name     string    // 文件名
	Size     int64     // 文件大小
	ModTime  time.Time // 修改时间
	IsDir    bool      // 是否是目录
	FileType string    // 文件类型
}

// handleBrowseDirectory 处理目录浏览请求
func handleBrowseDirectory(c *fiber.Ctx, path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	var entries []FileEntry
	for _, file := range files {
		// 获取文件类型
		fileType := ""
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".apk" {
				fileType = "apk"
			} else if ext == ".zip" {
				fileType = "zip"
			}
		}

		entries = append(entries, FileEntry{
			Name:     file.Name(),
			Size:     file.Size(),
			ModTime:  file.ModTime(),
			IsDir:    file.IsDir(),
			FileType: fileType,
		})
	}

	// 按照文件夹在前，文件在后的顺序排序
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return entries[i].Name < entries[j].Name
	})

	// 获取相对于根目录的路径
	rootDir := utils.GetServerConfig().Output
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return c.Status(500).SendString("Error resolving root path")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return c.Status(500).SendString("Error resolving directory path")
	}

	relPath, err := filepath.Rel(absRootDir, absPath)
	if err != nil {
		return c.Status(500).SendString("Error calculating relative path")
	}

	if relPath == "." {
		relPath = ""
	}

	rootPath := "/admin/browse"

	return c.Render("admin/directory", fiber.Map{
		"Title":    "目录浏览",
		"RootPath": rootPath,
		"Path":     relPath,
		"Entries":  entries,
		"Scripts": []string{
			"/static/js/admin/directory.js",
		},
	}, "admin/layout")
}

// handleBrowseFile 处理文件内容显示请求
func handleBrowseFile(c *fiber.Ctx, path string) error {
	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c.Status(404).SendString("File not found")
	}

	// 读取文件内容
	content, err := os.ReadFile(path)
	if err != nil {
		return c.Status(500).SendString("Error reading file")
	}

	// 获取相对于根目录的路径
	rootDir := utils.GetServerConfig().Output
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return c.Status(500).SendString("Error resolving root path")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return c.Status(500).SendString("Error resolving file path")
	}

	relPath, err := filepath.Rel(absRootDir, absPath)
	if err != nil {
		return c.Status(500).SendString("Error calculating relative path")
	}

	// 获取目录路径
	dirPath := filepath.Dir(relPath)
	if dirPath == "." {
		dirPath = ""
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "view", "browse", 0, fmt.Sprintf("查看文件：%s", path))

	rootPath := "/admin/browse"

	return c.Render("admin/file", fiber.Map{
		"Title":    "文件内容",
		"Path":     filepath.Base(relPath),
		"DirPath":  dirPath,
		"Content":  string(content),
		"RootPath": rootPath,
	}, "admin/layout")
}

// handleBrowseDelete 处理文件删除请求
func handleBrowseDelete(c *fiber.Ctx, path string) error {
	err := os.Remove(path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "文件删除失败",
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "delete", "browse", 0, fmt.Sprintf("删除文件：%s", path))

	return c.JSON(fiber.Map{
		"message": "文件删除成功",
	})
}

// HandleFTPUpload 处理 FTP 上传请求
func HandleFTPUpload(c *fiber.Ctx, rootPath string) error {
	filePath := c.Query("file")
	fileType := c.Query("type")

	if filePath == "" || (fileType != "apk" && fileType != "zip") {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的文件路径或类型",
		})
	}

	// 解码文件路径
	decodedPath, err := url.QueryUnescape(filePath)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的路径编码",
		})
	}

	fullPath := filepath.Join(rootPath, decodedPath)

	// 上传到 FTP
	if err := uploadToFTP(fullPath, fileType); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 记录操作日志
	adminlog.CreateAdminLog(c, "upload", "browse", 0, fmt.Sprintf("上传文件：%s", fullPath))

	return c.JSON(fiber.Map{
		"success": true,
		"message": "文件上传成功",
	})
}

// 上传文件到 FTP
func uploadToFTP(localPath string, fileType string) error {
	// 连接 FTP
	conn, err := ftp.Dial(fmt.Sprintf("%s:%s", ftpConfig.Host, ftpConfig.Port))
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("FTP连接失败: %v", err))
		return fmt.Errorf("FTP连接失败: %v", err)
	}
	defer conn.Quit()

	username, password, err := utils.ReadFromBinaryFile(utils.GetServerConfig().UserDataPath)
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("读取用户数据失败: %v", err))
		return fmt.Errorf("读取用户数据失败: %v", err)
	}

	// 登录
	if err := conn.Login(username, password); err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("FTP登录失败: %v", err))
		return fmt.Errorf("FTP登录失败: %v", err)
	}

	// 根据文件类型选择上传路径
	remotePath := ftpConfig.APKPath
	if fileType == "zip" {
		remotePath = ftpConfig.ZIPPath
	}

	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("打开文件失败: %v", err))
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 上传文件
	fileName := filepath.Base(localPath)
	err = conn.Stor(filepath.Join(remotePath, fileName), file)
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("上传文件失败: %v", err))
		return fmt.Errorf("上传文件失败: %v", err)
	}

	// 记录成功日志
	writeUploadLog(localPath, fileType, true, "上传成功")
	return nil
}

// 添加日志写入函数
func writeUploadLog(localPath, fileType string, success bool, message string) error {
	// 确保日志目录存在
	if err := os.MkdirAll(ftpConfig.LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 准备日志内容
	logContent := fmt.Sprintf("[%s] File: %s, Type: %s, Success: %v, Message: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		localPath,
		fileType,
		success,
		message,
	)

	// 获取当前日志文件路径
	logFile := filepath.Join(ftpConfig.LogDir, fmt.Sprintf("ftpupload_%s.log", time.Now().Format("20060102150405")))

	// 追加写入日志
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(logContent); err != nil {
		return fmt.Errorf("写入日志失败: %v", err)
	}

	return nil
}

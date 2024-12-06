package handlers

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/andycai/unitool/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jlaffaye/ftp"
)

// FileInfo 存储文件信息的结构体
type FileInfo struct {
	Name       string
	Size       int64
	FormatSize string
	Mode       os.FileMode
	ModTime    string
	IsDir      bool
	Path       string
	FileType   string
}

// BreadcrumbPath 存储面包屑路径的结构体
type BreadcrumbPath struct {
	Name string
	Path string
}

// FTP配置
type FTPConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	APKPath    string
	ZIPPath    string
	LogDir     string
	MaxLogSize int64
}

var outputPath string
var ftpConfig FTPConfig

const logTimeFormat = "20060102150405"

var currentLogFile string

// 初始化 FTP 配置
func InitFTP(config FTPConfig) {
	ftpConfig = config
}

// 添加日志写入函数
func writeUploadLog(localPath, fileType string, success bool, message string) error {
	// 确保日志目录存在
	if err := os.MkdirAll(ftpConfig.LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 检查当前日志文件
	if currentLogFile == "" || shouldCreateNewLog() {
		newLogFile := fmt.Sprintf("ftpupload_%s.log", time.Now().Format(logTimeFormat))
		currentLogFile = filepath.Join(ftpConfig.LogDir, newLogFile)
	}

	// 准备日志内容
	logContent := fmt.Sprintf("[%s] File: %s, Type: %s, Success: %v, Message: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		localPath,
		fileType,
		success,
		message,
	)

	// 追加写入日志
	f, err := os.OpenFile(currentLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(logContent); err != nil {
		return fmt.Errorf("写入日志失败: %v", err)
	}

	return nil
}

// 检查是否需要创建新的日志文件
func shouldCreateNewLog() bool {
	if currentLogFile == "" {
		return true
	}

	info, err := os.Stat(currentLogFile)
	if err != nil {
		return true
	}

	return info.Size() >= ftpConfig.MaxLogSize
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

// 添加处理 FTP 上传的路由处理函数
func HandleFTPUpload(c *fiber.Ctx, output string) error {
	filePath := c.Query("file")
	fileType := c.Query("type")

	if filePath == "" || (fileType != "apk" && fileType != "zip") {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid file path or type",
		})
	}

	// 解码文件路径
	decodedPath, err := url.QueryUnescape(filePath)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid path encoding",
		})
	}

	fullPath := filepath.Join(output, decodedPath)

	// 上传到 FTP
	if err := uploadToFTP(fullPath, fileType); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
	})
}

func HandleBrowseDirectory(c *fiber.Ctx, path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return c.Status(500).SendString("Error reading directory")
	}

	var fileInfos []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// 获取文件扩展名并判断类型
		fileType := ""
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".apk" {
				fileType = "apk"
			} else if ext == ".zip" {
				fileType = "zip"
			}
		}

		relativePath := trimPrefix(filepath.Join(path, entry.Name()))
		encodedPath := url.QueryEscape(relativePath)

		fileInfos = append(fileInfos, FileInfo{
			Name:       entry.Name(),
			Size:       info.Size(),
			FormatSize: formatSize(info.Size()),
			Mode:       info.Mode(),
			ModTime:    info.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:      entry.IsDir(),
			Path:       encodedPath,
			FileType:   fileType,
		})
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		if fileInfos[i].IsDir != fileInfos[j].IsDir {
			return fileInfos[i].IsDir
		}
		timeI, _ := time.Parse("2006-01-02 15:04:05", fileInfos[i].ModTime)
		timeJ, _ := time.Parse("2006-01-02 15:04:05", fileInfos[j].ModTime)
		return timeI.After(timeJ)
	})

	// 处理显示路径
	displayPath := filepath.ToSlash(path)
	if displayPath == filepath.ToSlash(outputPath) {
		displayPath = "."
	} else {
		displayPath = strings.TrimPrefix(strings.TrimPrefix(displayPath, filepath.ToSlash(outputPath)), "/")
		if displayPath == "" {
			displayPath = "."
		}
	}

	// 处理父目录路径
	parentPath := filepath.ToSlash(filepath.Dir(displayPath))
	if parentPath == "." || parentPath == displayPath {
		parentPath = "."
	}

	data := struct {
		Path            string
		ParentPath      string
		Files           []FileInfo
		BreadcrumbPaths []BreadcrumbPath
	}{
		Path:            displayPath,
		ParentPath:      parentPath,
		Files:           fileInfos,
		BreadcrumbPaths: generateBreadcrumbs(path),
	}

	return c.Render("admin/directory", fiber.Map{
		"Path":            data.Path,
		"ParentPath":      data.ParentPath,
		"Files":           data.Files,
		"BreadcrumbPaths": data.BreadcrumbPaths,
	})
}

func HandleBrowseFile(c *fiber.Ctx, path string) error {
	ext := strings.ToLower(filepath.Ext(path))

	textExts := map[string]bool{
		".txt": true, ".md": true, ".json": true,
		".go": true, ".js": true, ".css": true,
		".html": true, ".xml": true, ".yml": true,
		".yaml": true, ".conf": true, ".log": true,
	}

	if textExts[ext] {
		content, err := os.ReadFile(path)
		if err != nil {
			return c.Status(500).SendString("Error reading file")
		}

		relativePath := trimPrefix(path)
		relativeDirPath := trimPrefix(filepath.Dir(path))

		data := struct {
			Path    string
			DirPath string
			Content string
		}{
			Path:    relativePath,
			DirPath: relativeDirPath,
			Content: string(content),
		}

		return c.Render("admin/file", fiber.Map{
			"Path":    data.Path,
			"DirPath": data.DirPath,
			"Content": data.Content,
		})
	}

	return c.SendFile(path)
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMG"[exp])
}

func trimPrefix(path string) string {
	// 统一使用斜杠作为分隔符
	path = filepath.ToSlash(path)
	prefix := filepath.ToSlash(outputPath) + "/"
	trimmed := strings.TrimPrefix(path, prefix)
	if trimmed == "" {
		return "."
	}
	return trimmed
}

// 添加处理删除请求的函数
func HandleBrowseDelete(c *fiber.Ctx, path string) error {
	// 添加路径日志，帮助调试
	fmt.Printf("Attempting to delete file: %s\n", path)

	// 检查文件是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error stating file: %v\n", err)
		return c.Status(404).SendString(fmt.Sprintf("File not found: %s", path))
	}

	// 只允许删除文件，不允许删除目录
	if fileInfo.IsDir() {
		return c.Status(400).SendString("Cannot delete directories")
	}

	// 删除文件
	if err := os.Remove(path); err != nil {
		fmt.Printf("Error deleting file: %v\n", err)
		return c.Status(500).SendString(fmt.Sprintf("Failed to delete file: %v", err))
	}

	fmt.Printf("Successfully deleted file: %s\n", path)
	return c.SendString("File deleted successfully")
}

// 添加生成面包屑路径的函数
func generateBreadcrumbs(path string) []BreadcrumbPath {
	if path == "." || path == outputPath {
		return []BreadcrumbPath{}
	}

	// 统一使用斜杠作为分隔符
	path = filepath.ToSlash(path)
	outputPathSlash := filepath.ToSlash(outputPath)

	// 移除 outputPath 前缀
	path = strings.TrimPrefix(path, outputPathSlash)
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		return []BreadcrumbPath{}
	}

	parts := strings.Split(path, "/")
	breadcrumbs := make([]BreadcrumbPath, len(parts))

	for i := 0; i < len(parts); i++ {
		// 构建相对路径时也使用斜杠
		relativePath := strings.Join(parts[:i+1], "/")
		breadcrumbs[i] = BreadcrumbPath{
			Name: parts[i],
			Path: url.QueryEscape(relativePath),
		}
	}

	return breadcrumbs
}

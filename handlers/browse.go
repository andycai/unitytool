package handlers

import (
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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
}

var outputPath string
var templates *template.Template

func init() {
	// 初始化模板
	templates = template.Must(template.ParseFiles(
		"templates/directory.html",
		"templates/file.html",
	))
}

// HandleFileServer 处理文件服务器请求
func HandleFileServer(c *fiber.Ctx, output string) error {
	requestPath := c.Params("*")
	if requestPath == "" {
		requestPath = "."
	}

	// URL 解码路径
	decodedPath, err := url.QueryUnescape(requestPath)
	if err != nil {
		return c.Status(400).SendString("Invalid path encoding")
	}

	outputPath = output

	// 处理删除请求
	if c.Method() == "DELETE" {
		fullPath := filepath.Join(outputPath, decodedPath)
		return handleDelete(c, fullPath)
	}

	// 获取文件信息
	fullPath := filepath.Join(outputPath, decodedPath)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return c.Status(404).SendString(fmt.Sprintf("File not found: %s", fullPath))
	}

	// 如果是目录，显示目录内容
	if fileInfo.IsDir() {
		return handleDirectory(c, fullPath)
	}

	// 如果是文件，显示文件内容
	return handleFile(c, fullPath)
}

func handleDirectory(c *fiber.Ctx, path string) error {
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

		// URL 编码文件名
		relativePath := trimPrefix(filepath.Join(path, entry.Name()))
		encodedPath := url.QueryEscape(relativePath)

		fileInfos = append(fileInfos, FileInfo{
			Name:       entry.Name(),
			Size:       info.Size(),
			FormatSize: formatSize(info.Size()),
			Mode:       info.Mode(),
			ModTime:    info.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:      entry.IsDir(),
			Path:       encodedPath, // 使用编码后的路径
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

	relativePath := trimPrefix(path)
	parentPath := trimPrefix(filepath.Dir(relativePath))
	if parentPath == relativePath {
		parentPath = "."
	}

	// URL 编码父目录路径
	encodedParentPath := url.QueryEscape(parentPath)

	data := struct {
		Path       string
		ParentPath string
		Files      []FileInfo
	}{
		Path:       relativePath,
		ParentPath: encodedParentPath,
		Files:      fileInfos,
	}

	var buf strings.Builder
	if err := templates.ExecuteTemplate(&buf, "directory.html", data); err != nil {
		return c.Status(500).SendString("Error rendering template")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(buf.String())
}

func handleFile(c *fiber.Ctx, path string) error {
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

		var buf strings.Builder
		if err := templates.ExecuteTemplate(&buf, "file.html", data); err != nil {
			return c.Status(500).SendString("Error rendering template")
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(buf.String())
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
	path = filepath.ToSlash(path)
	prefix := filepath.ToSlash(outputPath) + "/"
	return strings.TrimPrefix(path, prefix)
}

// 添加处理删除请求的函数
func handleDelete(c *fiber.Ctx, path string) error {
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

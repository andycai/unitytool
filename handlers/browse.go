package handlers

import (
	"fmt"
	"html/template"
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

// HandleFileServer 处理文件服务器请求
func HandleFileServer(c *fiber.Ctx, output string) error {
	requestPath := c.Params("*")
	if requestPath == "" {
		requestPath = "."
	}
	outputPath = output
	requestPath = filepath.Join(outputPath, requestPath)

	// 获取文件信息
	fileInfo, err := os.Stat(requestPath)
	if err != nil {
		return c.Status(404).SendString("File or directory not found")
	}

	// 如果是目录，显示目录内容
	if fileInfo.IsDir() {
		return handleDirectory(c, requestPath)
	}

	// 如果是文件，显示文件内容
	return handleFile(c, requestPath)
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

		// 获取相对路径，需要考虑 windows 和 linux 路径分隔符不同的情况
		relativePath := trimPrefix(filepath.Join(path, entry.Name()))

		fileInfos = append(fileInfos, FileInfo{
			Name:       entry.Name(),
			Size:       info.Size(),
			FormatSize: formatSize(info.Size()),
			Mode:       info.Mode(),
			ModTime:    info.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:      entry.IsDir(),
			Path:       relativePath,
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

	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Directory listing for {{.Path}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .list { list-style: none; padding: 0; }
        .list li { padding: 8px; border-bottom: 1px solid #eee; }
        .list li:hover { background-color: #f5f5f5; }
        .dir { color: #2c3e50; font-weight: bold; }
        .file { color: #34495e; }
        .size { color: #7f8c8d; margin-left: 10px; }
        .time { color: #95a5a6; margin-left: 10px; }
        .empty-message { color: #666; text-align: center; padding: 20px; }
        .pagination { display: none; }
        .pagination.has-content { display: flex; justify-content: space-between; align-items: center; margin-top: 20px; }
    </style>
</head>
<body>
    <h1>Directory listing for {{.Path}}</h1>
    {{if .Files}}
    <ul class="list">
        {{if ne .Path "."}}
        <li>
            <a href="/browse/{{.ParentPath}}" class="dir">..</a>
        </li>
        {{end}}
        {{range .Files}}
        <li>
            {{if .IsDir}}
            <a href="/browse/{{.Path}}" class="dir">{{.Name}}/</a>
            {{else}}
            <a href="/browse/{{.Path}}" class="file">{{.Name}}</a>
            <span class="size">{{.FormatSize}}</span>
            {{end}}
            <span class="time">{{.ModTime}}</span>
        </li>
        {{end}}
    </ul>
    {{else}}
    <p class="empty-message">This folder is empty.</p>
    {{end}}
</body>
</html>`

	tmpl := template.Must(template.New("directory").Parse(html))

	relativePath := trimPrefix(path)
	parentPath := trimPrefix(filepath.Dir(relativePath))
	if parentPath == relativePath {
		parentPath = "."
	}

	data := struct {
		Path       string
		ParentPath string
		Files      []FileInfo
	}{
		Path:       relativePath,
		ParentPath: parentPath,
		Files:      fileInfos,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
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

		html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>File: {{.Path}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        pre { background-color: #f5f5f5; padding: 15px; border-radius: 5px; }
        .back { margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="back"><a href="/browse/{{.DirPath}}">← Back to directory</a></div>
    <h2>File: {{.Path}}</h2>
    <pre>{{.Content}}</pre>
</body>
</html>`

		tmpl := template.Must(template.New("file").Parse(html))

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
		if err := tmpl.Execute(&buf, data); err != nil {
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
	return strings.TrimPrefix(filepath.ToSlash(path), outputPath+"/")
}

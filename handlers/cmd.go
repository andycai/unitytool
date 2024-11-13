package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ScriptConfig struct {
	Path string            // 脚本文件路径
	Args []string          // 命令行参数
	Env  map[string]string // 环境变量
}

type ScriptForm struct {
	Name        string              `json:"name"`
	Repository  string              `json:"repository"`
	Platform    string              `json:"platform"`
	PublishType string              `json:"publishType"`
	Params      string              `json:"params"`
	Ext         []map[string]string `json:"ext"`
}

// 执行 shell 脚本
func ExecShell(c *fiber.Ctx, dir string) error {
	var form ScriptForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	scriptPath := filepath.Join(dir, form.Name)

	config := ScriptConfig{
		Path: scriptPath,
		Args: []string{form.Name, form.Repository, form.Platform, form.PublishType, form.Params},
		Env:  map[string]string{},
	}

	config.Env["repository"] = form.Repository
	config.Env["platform"] = form.Platform
	config.Env["publishType"] = form.PublishType
	config.Env["params"] = form.Params
	ext, err := json.Marshal(form.Ext)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// 对 ext 进行转移
	// ext = bytes.ReplaceAll(ext, []byte("\\"), []byte("\\\\"))
	str := strings.ReplaceAll(string(ext), "\"", "\\\"")
	config.Env["ext"] = str

	// 检查文件是否存在
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		return fmt.Errorf("script file not found: %s", config.Path)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(config.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		ext := filepath.Ext(config.Path)
		switch ext {
		case ".bat", ".cmd":
			args := append([]string{"/C", absPath}, config.Args...)
			cmd = exec.Command("cmd", args...)
		case ".ps1":
			args := append([]string{"-File", absPath}, config.Args...)
			cmd = exec.Command("powershell", args...)
		default:
			return fmt.Errorf("unsupported script type for Windows: %s", ext)
		}
	case "linux", "darwin":
		// 检查文件是否有执行权限
		info, err := os.Stat(absPath)
		if err != nil {
			return fmt.Errorf("failed to get file info: %v", err)
		}

		mode := info.Mode()
		if mode&0111 == 0 {
			if err := os.Chmod(absPath, mode|0111); err != nil {
				return fmt.Errorf("failed to set execute permission: %v", err)
			}
		}

		ext := filepath.Ext(config.Path)
		switch ext {
		case ".sh":
			args := append([]string{absPath}, config.Args...)
			cmd = exec.Command("bash", args...)
		default:
			// 直接执行文件，将参数传递给脚本
			args := append([]string{absPath}, config.Args...)
			cmd = exec.Command(absPath, args[1:]...)
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// 设置工作目录为脚本所在目录
	cmd.Dir = filepath.Dir(absPath)

	// 设置环境变量
	if len(config.Env) > 0 {
		env := os.Environ()
		for k, v := range config.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = env
	}

	// 将标准输出和错误输出设置为程序的标准输出
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute script: %v", err)
	}

	return c.SendString(string(output))
}

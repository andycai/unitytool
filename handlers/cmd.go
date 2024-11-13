package handlers

import (
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

// 执行 shell 脚本
func ExecShell(c *fiber.Ctx, dir string) error {
	// 获取 get 参数，包括 name、type、platform、ext
	name := c.Query("name")
	typ := c.Query("type")
	platform := c.Query("platform")
	ext := c.Query("ext")
	params := c.Query("params")

	// params 是使用 ｜ 连接的字符串，使用 ｜ 分割 作为动态参数，解析出来放到 args 中
	args := []string{}
	if params != "" {
		arr := strings.Split(params, ",")
		args = append(args, arr...)
	}

	scriptPath := fmt.Sprintf("%s/%s", dir, name)

	config := ScriptConfig{
		Path: scriptPath,
		Args: []string{name, typ, platform, ext},
		Env:  map[string]string{},
	}

	// args 加到 config.Args 中
	config.Args = append(config.Args, args...)

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

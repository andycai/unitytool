package unibuild

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

// UnityBuildConfig Unity构建配置
type UnityBuildConfig struct {
	UnityPath    string
	ProjectPath  string
	BuildMethod  string
	OutputPath   string
	BuildTarget  string
	BuildOptions string
	LogFilePath  string
}

// executeSVNOperations 执行SVN操作
func executeSVNOperations(projectPath string) error {
	// SVN revert
	revertCmd := exec.Command("svn", "revert", "-R", projectPath)
	output, err := revertCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("svn revert failed: %v\nOutput: %s", err, string(output))
	}

	// SVN update
	updateCmd := exec.Command("svn", "update", projectPath)
	output, err = updateCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("svn update failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// buildResources 处理AssetBundle打包请求
func buildResources(c *fiber.Ctx) error {
	config := UnityBuildConfig{
		ProjectPath:  c.Query("projectPath", ""),
		OutputPath:   c.Query("outputPath", ""),
		BuildMethod:  c.Query("method", "BuildAssetBundles"),
		BuildTarget:  c.Query("target", "Android"),
		BuildOptions: c.Query("options", ""),
	}

	if config.ProjectPath == "" || config.OutputPath == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Project path and output path are required",
		})
	}

	// 执行SVN操作
	if err := executeSVNOperations(config.ProjectPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("SVN operation failed: %v", err),
		})
	}

	config.UnityPath = getUnityEditorPath()
	if config.UnityPath == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "Unity editor path not found",
		})
	}

	err := executeUnityBuild(config)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("Build failed: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "AssetBundle build completed",
		"output":  config.OutputPath,
	})
}

// buildApp 处理APK打包请求
func buildApp(c *fiber.Ctx) error {
	config := UnityBuildConfig{
		ProjectPath:  c.Query("projectPath", ""),
		OutputPath:   c.Query("outputPath", ""),
		LogFilePath:  c.Query("logFilePath", ""),
		BuildMethod:  c.Query("method", "BuildAndroid"),
		BuildTarget:  "Android",
		BuildOptions: c.Query("options", ""),
	}

	if config.ProjectPath == "" || config.OutputPath == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Project path and output path are required",
		})
	}

	// 执行SVN操作
	if err := executeSVNOperations(config.ProjectPath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("SVN operation failed: %v", err),
		})
	}

	config.UnityPath = getUnityEditorPath()
	if config.UnityPath == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "Unity editor path not found",
		})
	}

	err := executeUnityBuild(config)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("Build failed: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "APK build completed",
		"output":  config.OutputPath,
	})
}

// executeUnityBuild 执行Unity命令行构建
func executeUnityBuild(config UnityBuildConfig) error {
	args := []string{
		"-quit",
		"-batchmode",
		"-nographics",
		"silent-crashes",
		"-disable-assembly-updater",
		"-projectPath", config.ProjectPath,
		"-executeMethod", config.BuildMethod,
		"-buildTarget", config.BuildTarget,
		"-outputPath", config.OutputPath,
		"-logFile", config.LogFilePath,
	}

	if config.BuildOptions != "" {
		args = append(args, "-buildOptions", config.BuildOptions)
	}

	cmd := exec.Command(config.UnityPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unity build failed: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// getUnityEditorPath 获取Unity编辑器路径
func getUnityEditorPath() string {
	switch runtime.GOOS {
	case "darwin":
		return "/Applications/Unity/Hub/Editor/2021.3.21f1/Unity.app/Contents/MacOS/Unity"
	case "windows":
		return filepath.Join("C:", "Program Files", "Unity", "Hub", "Editor", "2021.3.21f1", "Editor", "Unity.exe")
	default:
		return ""
	}
}

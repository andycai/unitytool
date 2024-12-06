package utils

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"path/filepath"
	"time"
)

// FileEntry 存储文件信息的结构体
type FileEntry struct {
	Name    string    // 文件名
	Size    int64     // 文件大小
	ModTime time.Time // 修改时间
	IsDir   bool      // 是否是目录
}

// FormatSize 格式化文件大小
func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return formatBytes(size, "B")
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return formatBytes(size/div, "KMGTPE"[exp:exp+1]+"B")
}

// formatBytes 格式化字节数
func formatBytes(size int64, unit string) string {
	if size == 0 {
		return "0" + unit
	}
	return string(rune(size)) + unit
}

// IsFile 判断路径是否为文件
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// GetParentDirectory 获取父目录
func GetParentDirectory(path string) string {
	return filepath.Dir(path)
}

func WriteToBinaryFile(filename, username, password string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将用户名和密码进行 Base64 编码
	usernameBase64 := base64.StdEncoding.EncodeToString([]byte(username))
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))

	// 写入用户名长度和内容
	usernameLength := uint32(len(usernameBase64))
	if err := binary.Write(file, binary.LittleEndian, usernameLength); err != nil {
		return err
	}
	if _, err := file.Write([]byte(usernameBase64)); err != nil {
		return err
	}

	// 写入密码长度和内容
	passwordLength := uint32(len(passwordBase64))
	if err := binary.Write(file, binary.LittleEndian, passwordLength); err != nil {
		return err
	}
	if _, err := file.Write([]byte(passwordBase64)); err != nil {
		return err
	}

	return nil
}

func ReadFromBinaryFile(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// 读取用户名
	var usernameLength uint32
	if err := binary.Read(file, binary.LittleEndian, &usernameLength); err != nil {
		return "", "", err
	}
	usernameBase64 := make([]byte, usernameLength)
	if _, err := file.Read(usernameBase64); err != nil {
		return "", "", err
	}
	usernameBytes, err := base64.StdEncoding.DecodeString(string(usernameBase64))
	if err != nil {
		return "", "", err
	}
	username := string(usernameBytes)

	// 读取密码
	var passwordLength uint32
	if err := binary.Read(file, binary.LittleEndian, &passwordLength); err != nil {
		return "", "", err
	}
	passwordBase64 := make([]byte, passwordLength)
	if _, err := file.Read(passwordBase64); err != nil {
		return "", "", err
	}
	passwordBytes, err := base64.StdEncoding.DecodeString(string(passwordBase64))
	if err != nil {
		return "", "", err
	}
	password := string(passwordBytes)

	return username, password, nil
}

// func test() {
// 	filename := "user_data.bin"
// 	username := "testuser"
// 	password := "mypassword123"

// 	// 写入
// 	if err := WriteToBinaryFile(filename, username, password); err != nil {
// 		fmt.Println("写入错误:", err)
// 		return
// 	}

// 	// 读取
// 	readUsername, readPassword, err := ReadFromBinaryFile(filename)
// 	if err != nil {
// 		fmt.Println("读取错误:", err)
// 		return
// 	}

// 	fmt.Printf("读取到的用户名: %s, 密码: %s\n", readUsername, readPassword)
// }

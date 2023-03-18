package cli

import (
	"os"
)

// 获取本地配置
func GetLocalConfig() string {
	file := GetLocalConfigPath()
	content, err := os.ReadFile(file)
	if err == nil {
		return string(content)
	}
	return "未配置本地地址"
}

// 获取本地缓存目录
func GetLocalConfigPath() string {
	dir, _ := os.UserHomeDir()
	file := dir + "/.chaxelib_local"
	return file
}

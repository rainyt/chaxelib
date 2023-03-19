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
	return ""
}

// 获得授权码
func GetAccestCode() string {
	file := GetLocalConfigPwdPath()
	content, err := os.ReadFile(file)
	if err == nil {
		return string(content)
	}
	return ""
}

// 获得本地授权码储存
func GetLocalConfigPwdPath() string {
	return GetLocalConfigPath() + "_pwd"
}

// 获取本地缓存目录
func GetLocalConfigPath() string {
	dir, _ := os.UserHomeDir()
	file := dir + "/.chaxelib_local"
	return file
}

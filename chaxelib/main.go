package main

import (
	"bufio"
	"fmt"
	"haxelib/v2/chaxelib/cli"
	"os"
	"strings"
)

var VERSION = "1.0.0"

type CommandParams struct {
	list map[string][]string
}

func main() {
	if len(os.Args) >= 2 {
		var command = os.Args[1]
		switch command {
		case "local":
			// 本地化配置
			dir, _ := os.UserHomeDir()
			file := dir + "/.chaxelib_local"
			fmt.Println("当前配置：", cli.GetLocalConfig())
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("请输入本地化haxelib储存服务器IP地址(127.0.0.1:5000):")
			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			fmt.Println("服务器配置成功：", text)
			writeErr := os.WriteFile(file, []byte(strings.ReplaceAll(text, "\n", "")), 0666)
			if writeErr != nil {
				panic(writeErr)
			}
		case "clone":
			// 镜像
			lib := os.Args[2]
			if strings.Contains(lib, ":") {
				params := strings.Split(lib, ":")
				cli.CloneHaxelib(params[0], params[1])
			} else {
				if len(os.Args) >= 4 {
					cli.CloneHaxelib(os.Args[2], os.Args[3])
				} else {
					cli.CloneHaxelib(os.Args[2], "")
				}
			}
		case "install":
			// 安装
			lib := os.Args[2]
			if strings.Contains(lib, ":") {
				params := strings.Split(lib, ":")
				cli.InstallHaxelib(params[0], params[1])
			} else {
				cli.InstallHaxelib(os.Args[2], "")
			}
		default:
			fmt.Println("不支持" + command + "命令")
		}
	} else {
		fmt.Println("CHaxelib (CN) Tools version:", VERSION)
		fmt.Println("  Usage: haxelib [command] [options]")
		params := CommandParams{
			list: map[string][]string{
				"基础":  {"install#通过库名安装第三方库"},
				"镜像":  {"clone#通过库名进行镜像，也可以通过该命令查询镜像情况"},
				"本地化": {"local#配置本地服务器IP，可绑定haxelib服务器，会优先从本地安装，配置之后，才允许使用upload命令"},
				"上传":  {"upload#上传haxelib库到本地服务器"},
				"更新":  {"update#会优先从本地进行检查，不存在时，将从线上进行更新"},
			},
		}
		maxlen := 0
		for _, v := range params.list {
			for _, v2 := range v {
				p := strings.Split(v2, "#")
				if len(p[0]) > maxlen {
					maxlen = len(p[0])
				}
			}
		}
		for k, v := range params.list {
			fmt.Println("  " + k + ":")
			for _, v2 := range v {
				p := strings.Split(v2, "#")
				space := ""
				l := maxlen - len(p[0])
				for i := 0; i < l; i++ {
					space += " "
				}
				fmt.Println("    " + p[0] + space + ":" + p[1])
			}
		}
	}
}

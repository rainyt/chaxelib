package main

import (
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
	if len(os.Args) >= 3 {
		var command = os.Args[1]
		switch command {
		case "clone":
			// 镜像
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
				"基础": {"install#通过库名安装第三方库"},
				"镜像": {"clone#通过库名进行镜像，也可以通过该命令查询镜像情况"},
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

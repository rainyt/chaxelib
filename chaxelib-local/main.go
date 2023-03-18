package main

import (
	"archive/zip"
	"fmt"
	"haxelib/v2/chaxelib/cli"
	"net/http"
	"net/rpc"
	"os"
	"strings"
)

// Haxelib RPC实现
type Haxelib struct{}

// 获取Haxelib库的下载地址
func (h *Haxelib) GetHaxelibUrl(haxelibname string, ret *string) error {
	args := strings.Split(haxelibname, ":")
	version := ""
	fmt.Println(args)
	if len(args) == 1 {
		haxelibname = args[0]
		// 读取配置last配置
		last := "haxelib/" + haxelibname + "/last"
		_, e2 := os.Stat(last)
		if e2 != nil {
			return e2
		}
		b, _ := os.ReadFile(last)
		version = string(b)
	} else {
		haxelibname = args[0]
		version = args[1]
	}
	p, e := FindHaxelib(haxelibname, version)
	if e == nil {
		*ret = p
	}
	return e
}

func (h *Haxelib) UploadHaxelib(bytes []byte, ret *int) error {
	fmt.Println("接收到二进制数据", len(bytes))
	temp := "haxelib/temp.zip"
	defer os.Remove(temp)
	defer os.RemoveAll("haxelib/temp")
	err := os.WriteFile(temp, bytes, 0777)
	z, readerr := zip.OpenReader(temp)
	if readerr != nil {
		return readerr
	}
	defer z.Close()
	ziperr := cli.Unzip("haxelib/temp", &z.Reader)
	if ziperr != nil {
		return ziperr
	}
	// 识别是否有效的haxelib.json
	err3 := SaveHaxelib("haxelib/temp", bytes)
	if err3 != nil {
		return err3
	}
	return err
}

func InitConfig() {
	// 检测本地haxelib目录是否存在，如不存在，则需要创建
	os.Mkdir("haxelib", 0777)
}

func main() {
	InitConfig()
	haxelib := &Haxelib{}
	rpc.Register(haxelib)
	rpc.HandleHTTP()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		bytes, err := os.ReadFile("." + r.URL.Path)
		if err == nil {
			w.Write(bytes)
		} else {
			fmt.Println(err)
			w.Write([]byte(r.URL.Path + " is not found"))
		}
	})
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"archive/zip"
	"fmt"
	"haxelib/v2/chaxelib/cli"
	"net/http"
	"net/rpc"
	"os"
)

// Haxelib RPC实现
type Haxelib struct{}

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
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		panic(err)
	}
}

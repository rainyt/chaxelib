package cli

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
)

// 上传haxelib本地库支持
func UploadHaxeLibDir(dir string) {
	_, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	os.Chdir(dir + "/../")
	wd, _ := os.Getwd()
	baseFile := filepath.Base(dir)
	tempFile := "temp"
	fmt.Println("当前目录：", wd)
	// 压缩文件
	defer os.Remove(wd + "/" + tempFile + ".zip")
	zip, _ := os.Create(wd + "/" + tempFile + ".zip")
	defer zip.Close()
	Zip(zip, wd+"/"+baseFile)
	// 再将文件上传
	// 1.连接远程rpc服务
	conn, err := rpc.DialHTTP("tcp", GetLocalConfig())
	if err != nil {
		log.Fatal(err)
	}
	ret := 0
	bytes, _ := os.ReadFile(wd + "/" + tempFile + ".zip")
	err2 := conn.Call("Haxelib.UploadHaxelib", bytes, &ret)
	if err2 != nil {
		log.Fatal(err2)
	}
}

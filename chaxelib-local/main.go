package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"haxelib/v2/chaxelib/cli"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"sync"
)

// Haxelib RPC实现
type Haxelib struct{}

var (
	Passworld = flag.String("pwd", "", "设定授权码，当存在授权码时，需要请求时提供授权码进行登录下载，否则会被拒绝")
	Port      = flag.String("port", "5555", "设定端口，默认5555")
)

// 获取Haxelib库的下载地址
func (h *Haxelib) GetHaxelibUrl(haxelibname string, ret *string) error {
	args := strings.Split(haxelibname, ":")
	version := ""
	fmt.Println(args)
	if len(args) == 1 {
		haxelibname = args[0]
		// 读取配置last配置
		version = GetLastVersion(haxelibname)
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

var upload_lock sync.Mutex

func (h *Haxelib) UploadHaxelib(bytes []byte, ret *int) error {
	// Passworld检查授权码
	if *Passworld != "" {
		l := len(*Passworld)
		pwd := string(bytes[len(bytes)-l:])
		if pwd != *Passworld {
			return fmt.Errorf("授权码不正确")
		}
		bytes = bytes[0 : len(bytes)-l]
	}
	// 上传haxelib进行安全锁处理
	upload_lock.Lock()
	defer upload_lock.Unlock()
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
	flag.Parse()
	InitConfig()
	haxelib := &Haxelib{}
	rpc.Register(haxelib)
	rpc.HandleHTTP()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if *Passworld != "" {
			pwd := r.URL.Query().Get("pwd")
			if pwd != *Passworld {
				w.Write([]byte(r.URL.Path + " accect is error."))
				return
			}
		}
		fmt.Println(r.URL.Path)
		bytes, err := os.ReadFile("." + r.URL.Path)
		if err == nil {
			w.Write(bytes)
		} else {
			fmt.Println(err)
			w.Write([]byte(r.URL.Path + " is not found"))
		}
	})
	err := http.ListenAndServe(":"+*Port, nil)
	if err != nil {
		panic(err)
	}
}

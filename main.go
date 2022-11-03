package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	OssId       = flag.String("ossid", "", "阿里云OSS_ACCESSKEY_ID")
	OssSecret   = flag.String("osssecret", "", "阿里云OSS_ACCESSKEY_SECRET")
	OssEndpoint = flag.String("endpoint", "", "阿里云Endpoint")
	OssBucket   = flag.String("bucket", "none", "阿里云Bucket")
	Port        = flag.Int("port", 80, "阿里云Bucket")
)

// pathExists 判断一个文件或文件夹是否存在
// 输入文件路径，根据返回的bool值来判断文件或文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 上传到OSS服务器
func uploadOSS(filename string, data []byte) {
	// 移除第一个/
	filename = strings.Join(strings.Split(filename, "")[1:], "")
	fmt.Println("开始镜像到OSS", filename)
	exist, _ := pathExists(filename)
	if exist {
		fmt.Println("已经在镜像进行中")
		return
	}
	client, err := oss.New(*OssEndpoint, *OssId, *OssSecret)
	if err != nil {
		fmt.Println("无法链接Oss服务器", err.Error())
		return
	}
	bucket, err := client.Bucket(*OssBucket)
	if err != nil {
		fmt.Println("无法链接Bucket:"+*OssBucket, err.Error())
		return
	}
	// 判断是否已经镜像好了
	existObject, _ := bucket.IsObjectExist(filename)
	if existObject {
		fmt.Println(filename, "已镜像")
		return
	}
	werr := ioutil.WriteFile(filename, data, 0666)
	if werr != nil {
		fmt.Println("WriteFile Error:", werr.Error())
		return
	}
	err = bucket.PutObjectFromFile(filename, filename)
	if err != nil {
		fmt.Println("文件"+filename+"无法上传到OSS", err.Error())
	} else {
		fmt.Println("镜像成功", filename)
	}
	os.Remove(filename)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("请求地址：", r.URL.Path, r.Method)
		index := strings.Index(r.URL.Path, "files/3.0/")
		if index != -1 {
			rep, err := http.Get("https://lib.haxe.org" + r.URL.Path)
			if err != nil {
				fmt.Println("请求错误：", err.Error())
			} else {
				defer rep.Body.Close()
				bytes, b := ioutil.ReadAll(rep.Body)
				if b == nil {
					// 开始上传到OSS
					go uploadOSS(r.URL.Path, bytes)
					w.Write(bytes)
				} else {
					fmt.Println("请求错误：", b.Error())
				}
			}
		} else {
			w.Write([]byte("Not support the Path"))
		}
	})
	dir := "files/3.0/"
	direrr := os.MkdirAll(dir, 0777)
	if direrr != nil {
		panic(direrr)
	}
	fmt.Println("服务器启动：" + fmt.Sprint(*Port))
	err := http.ListenAndServe(":"+fmt.Sprint(*Port), nil)
	if err != nil {
		fmt.Println("服务器错误：", err.Error())
		panic(err)
	}
}

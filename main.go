package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
	OssUrl      = flag.String("ossurl", "", "阿里云下载地址")
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

// 判断OSS是否存在镜像库
func existOSS(filename string) bool {
	filename = strings.Join(strings.Split(filename, "")[1:], "")
	fmt.Println("查询镜像库：", filename)
	client, err := oss.New(*OssEndpoint, *OssId, *OssSecret)
	if err != nil {
		fmt.Println("无法链接Oss服务器", err.Error())
		return false
	}
	bucket, err := client.Bucket(*OssBucket)
	if err != nil {
		fmt.Println("无法链接Bucket:"+*OssBucket, err.Error())
		return false
	}
	// 判断是否已经镜像好了
	existObject, _ := bucket.IsObjectExist(filename)
	if existObject {
		fmt.Println(filename, "已镜像")
		return true
	}
	return false
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
	werr := os.WriteFile(filename, data, 0666)
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

type RetData struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

func sendData(w http.ResponseWriter, data RetData) {
	content, _ := json.Marshal(data)
	w.Write([]byte(content))
}

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ossIndex := strings.Index(r.URL.Path, "/oss")
		if ossIndex == 0 {
			// 查询oss镜像库
			queryUrl := strings.ReplaceAll(r.URL.Path, "/oss", "")
			if existOSS(queryUrl) {
				sendData(w, RetData{
					Code: 0,
					Data: map[string]any{
						"url": *OssUrl + queryUrl,
					},
				})
			} else {
				sendData(w, RetData{
					Code: -1,
					Data: "不存在镜像",
				})
			}
			return
		}
		cloneIndex := strings.Index(r.URL.Path, "/clone")
		if cloneIndex == 0 {
			// 查询克隆结果
			// 查询oss镜像库
			queryUrl := strings.ReplaceAll(r.URL.Path, "/clone", "")
			if existOSS(queryUrl) {
				sendData(w, RetData{
					Code: 0,
					Data: "已镜像完成",
				})
			} else {
				// 尝试镜像
				fmt.Println("请求地址：", queryUrl)
				bdata, err := readHaxelib(queryUrl)
				if err != nil {
					fmt.Println(err.Error())
					sendData(w, RetData{
						Code: -1,
						Data: err.Error(),
					})
				} else {
					// 开始上传到OSS
					go uploadOSS(queryUrl, bdata)
					sendData(w, RetData{
						Code: -1,
						Data: "正在镜像中...",
					})
				}
			}
			return
		}
		fmt.Println("请求地址：", r.URL.Path, r.Method)
		bdata, err := readHaxelib(r.URL.Path)
		if err != nil {
			fmt.Println(err.Error())
			sendData(w, RetData{
				Code: -1,
				Data: err.Error(),
			})
		} else {
			// 开始上传到OSS
			go uploadOSS(r.URL.Path, bdata)
			w.Write(bdata)
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

// 读取Haxelib库
func readHaxelib(path string) ([]byte, error) {
	index := strings.Index(path, "files/3.0/")
	if index != -1 {
		rep, err := http.Get("https://lib.haxe.org" + path)
		if err != nil {
			fmt.Println("请求错误：", err.Error())
			return nil, err
		} else {
			if rep.StatusCode == 200 {
				defer rep.Body.Close()
				bytes, b := io.ReadAll(rep.Body)
				if b == nil {
					return bytes, nil
				} else {
					return nil, fmt.Errorf("StatuCode error:" + fmt.Sprint(rep.StatusCode) + b.Error())
				}
			} else {
				return nil, fmt.Errorf("Not exist the Haxelib:" + path)
			}
		}
	} else {
		return nil, fmt.Errorf("Not support the Path" + path)
	}
}

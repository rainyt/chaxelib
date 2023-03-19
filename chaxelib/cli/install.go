package cli

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var Haxelib_path = "https://haxelib.zygame.cc/"

// 带加载进度的
type Reader struct {
	io.Reader
	Currnet int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	r.Currnet += int64(n)
	fmt.Printf("\r %0.2fmb", float64(r.Currnet)/1024./1024.)
	return
}

// Haxelib.json的格式
type HaxelibData struct {
	Dependencies map[string]any
}

// 检测依赖是否已存在
func existHaxelib(libname string) bool {
	c := exec.Command("haxelib", "path", libname)
	c.Start()
	err := c.Wait()
	return err == nil
}

// 读取haxelib.json配置
func readHaxelibJson(z []*zip.File) *HaxelibData {
	for _, f := range z {
		if f.FileInfo().Name() == "haxelib.json" {
			rw, _ := f.Open()
			bytes, bytesErr := io.ReadAll(rw)
			if bytesErr != nil {
				panic(bytesErr)
			}
			haxelibJson := &HaxelibData{}
			json.Unmarshal(bytes, haxelibJson)
			return haxelibJson
		}
	}
	return nil
}

// 获取项目版本号
func getProjectVersion(libname string) []string {
	var versions = []string{}
	var query = "https://lib.haxe.org/p/" + libname + "/versions/"
	h, e := http.Get(query)
	if e != nil {
		panic(e)
	} else {
		defer h.Body.Close()
		b, _ := io.ReadAll(h.Body)
		content := string(b)
		// 正则条件，获取所有支持的版本号
		// fmt.Println(content)
		re := regexp.MustCompile(">[0-9.]+</a")
		ret := re.FindAllStringSubmatch(content, -1)
		for _, v := range ret {
			version := v[0]
			version = strings.ReplaceAll(version, ">", "")
			version = strings.ReplaceAll(version, "</a", "")
			versions = append(versions, version)
		}
	}
	// fmt.Println(versions)
	return versions
}

// 检测版本是否有效
func CheckVersion(libname string, version string) (string, error) {
	versions := getProjectVersion(libname)
	if len(versions) == 0 {
		return "", fmt.Errorf("库" + libname + "不存在")
	}
	if version != "" {
		hasVersion := false
		for _, v := range versions {
			if v == version {
				hasVersion = true
				break
			}
		}
		if !hasVersion {
			return "", fmt.Errorf("库" + libname + "不存在" + version + "版本")
		}
	} else {
		version = versions[0]
	}
	version = strings.ReplaceAll(version, ".", ",")
	return version, nil
}

// 通过本地化服务器更新库
func UpdateHaxelib(libname string, version string) error {
	fmt.Println("update lib:", libname, version)
	conn, err := rpc.DialHTTP("tcp", GetLocalConfig())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	path := ""
	if version != "" {
		libname += ":" + version
	}
	err2 := conn.Call("Haxelib.GetHaxelibUrl", libname, &path)
	if err2 != nil {
		fmt.Println(err2.Error())
		return err2
	}
	fmt.Println("查询到的路径", path)
	baseName := filepath.Base(path)
	// 开始安装
	downloadUrl := "http://" + GetLocalConfig() + "/" + path
	pwd := GetAccestCode()
	if pwd != "" {
		downloadUrl += "?pwd=" + pwd
	}
	downloadPath(baseName, downloadUrl)
	return nil
}

func InstallHaxelib(libname string, version string) {
	// 远程服务器
	version, err := CheckVersion(libname, version)
	if err != nil {
		panic(err)
	}
	libzipfile := libname + "-" + version + ".zip"
	// 做一个检测
	ossurl := Haxelib_path + "oss/files/3.0/" + libzipfile
	ossret, e := http.Get(ossurl)
	if e != nil {
		panic(e)
	} else {
		defer ossret.Body.Close()
		data, _ := io.ReadAll(ossret.Body)
		var jsonData map[string]any
		json.Unmarshal(data, &jsonData)
		println("镜像结果", string(data), jsonData["code"].(float64))
		if jsonData["code"].(float64) == 0 {
			downloadPath(libzipfile, jsonData["data"].(map[string]any)["url"].(string))
			return
		}
	}

	liburl := Haxelib_path + "files/3.0/" + libzipfile
	downloadPath(libzipfile, liburl)
}

func downloadPath(libzipfile string, liburl string) {
	fmt.Println("正在下载：" + liburl)
	h, e := http.Get(liburl)
	if e != nil {
		panic(e)
	} else {
		defer h.Body.Close()
		file, err := os.Create(libzipfile)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := &Reader{
			Reader: h.Body,
		}
		io.Copy(file, reader)
		fmt.Println()
		installLocalZip(libzipfile)
	}
}

func installLocalZip(zipfile string) {
	// 安装完成后，将压缩包删除
	defer os.Remove(zipfile)
	// 先安装依赖，避免haxelib检测依赖
	curZip, zipErr := zip.OpenReader(zipfile)
	fmt.Println("检测依赖：", zipfile)
	if zipErr != nil {
		panic(zipErr)
	} else {
		// 解析haxelib.json
		haxelibJson := readHaxelibJson(curZip.File)
		if haxelibJson != nil {
			fmt.Println("检测依赖...")
			if haxelibJson.Dependencies != nil {
				// 检查依赖
				for k, v := range haxelibJson.Dependencies {
					if !existHaxelib(k) {
						InstallHaxelib(k, v.(string))
					} else {
						fmt.Println("依赖" + k + "已安装")
					}
				}
			}
		}
	}

	fmt.Println("开始安装：" + zipfile)
	c := exec.Command("haxelib", "install", zipfile)
	stdout, _ := c.StdoutPipe()
	e := c.Start()
	if e != nil {
		panic(e)
	}
	output, _ := io.ReadAll(stdout)
	fmt.Println(string(output))
}

// 从hxml配置中进行更新haxelib
func UpdateHaxelibByHxml(hxmlpath string) {
	b, err := os.ReadFile(hxmlpath)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("haxelib", "list")
	output, _ := cmd.StdoutPipe()
	cmd.Start()
	all, _ := io.ReadAll(output)
	haxelibReg, _ := regexp.Compile("[^[dev\n][a-zA-Z-_0-9]+:")
	allHaxelib := haxelibReg.FindAllString(string(all), -1)
	hxml := string(b)
	reg, _ := regexp.Compile("-D [a-zA-Z-_0-9]+")
	uarr := reg.FindAllString(hxml, -1)
	haxelib := []string{}
	for _, v := range uarr {
		lib := strings.ReplaceAll(v, "-D ", "")
		for _, v2 := range allHaxelib {
			if v2 == lib+":" {
				haxelib = append(haxelib, lib)
				break
			}
		}
	}
	haxelibPath := CmtToString("haxelib config")
	haxelibPath = strings.ReplaceAll(haxelibPath, "\n", "")
	fmt.Println("haxelibPath = ", haxelibPath)
	for _, v := range haxelib {
		// 先从本地获取库的环境配置，如果是dev版本，则不会更新，如果是普通版本，则会从线上更新
		dev := IsHaxelibDev(haxelibPath, v)
		if !dev {
			UpdateHaxelib(v, GetHaxelibCurrentVersion(haxelibPath, v)+"+")
		}
	}
}

// 获取Haxelib当前版本
func GetHaxelibCurrentVersion(config string, libname string) string {
	path := config + libname + "/.current"
	_, err := os.Stat(path)
	if err == nil {
		r, _ := os.ReadFile(path)
		return string(r)
	}
	return ""
}

// 该haxelib是否dev环境
func IsHaxelibDev(config string, libname string) bool {
	path := config + libname + "/.dev"
	_, err := os.Stat(path)
	return err == nil
}

// 执行命令并转成string类型
func CmtToString(cmd string) string {
	args := strings.Split(cmd, " ")
	pathCmd := exec.Command(args[0], args[1:]...)
	pathOutput, _ := pathCmd.StdoutPipe()
	pathCmd.Start()
	pathString, _ := io.ReadAll(pathOutput)
	return string(pathString)
}

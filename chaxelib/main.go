package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var haxelib_path = "https://haxelib.zygame.cc/"

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

func main() {
	fmt.Println(os.Args)
	if len(os.Args) >= 3 {
		var command = os.Args[1]
		switch command {
		case "install":
			lib := os.Args[2]
			if strings.Contains(lib, ":") {
				params := strings.Split(lib, ":")
				installHaxelib(params[0], params[1])
			} else {
				installHaxelib(os.Args[2], "")
			}
		default:
			fmt.Println("不支持" + command + "命令")
		}
	}
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
		b, _ := ioutil.ReadAll(h.Body)
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

func installHaxelib(libname string, version string) {
	versions := getProjectVersion(libname)
	if len(versions) == 0 {
		panic("库" + libname + "不存在")
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
			panic("库" + libname + "不存在" + version + "版本")
		}
	} else {
		version = versions[0]
	}
	version = strings.ReplaceAll(version, ".", ",")
	libzipfile := libname + "-" + version + ".zip"
	liburl := haxelib_path + "files/3.0/" + libzipfile
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
	fmt.Println("开始安装：" + zipfile)
	c := exec.Command("haxelib", "install", zipfile)
	stdout, _ := c.StdoutPipe()
	e := c.Start()
	if e != nil {
		panic(e)
	}
	output, _ := ioutil.ReadAll(stdout)
	fmt.Println(string(output))
}

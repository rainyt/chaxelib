package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	url := "https://lib.haxe.org/all/"
	fmt.Println("Clone", url)
	r, e := http.Get(url)
	if e == nil {
		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		} else {
			reg := regexp.MustCompile("/p/[a-z]+/")
			content := string(data)
			arrays := reg.FindAllStringSubmatch(content, -1)
			allCounts := len(arrays)
			for _, v := range arrays {
				libname := v[0]
				libname = strings.ReplaceAll(libname, "/p/", "")
				libname = strings.ReplaceAll(libname, "/", "")
				fmt.Printf("开始克隆%s", libname)
				c := exec.Command("chaxelib", "clone", libname)
				out, _ := c.StdoutPipe()
				err := c.Start()
				if err != nil {
					panic(err)
				}
				data, _ := io.ReadAll(out)
				fmt.Println(string(data))
			}
			fmt.Println("镜像总数：" + fmt.Sprint(allCounts))
		}
	}
}

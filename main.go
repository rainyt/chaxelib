package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
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
					w.Write(bytes)
				} else {
					fmt.Println("请求错误：", b.Error())
				}
			}
		} else {
			w.Write([]byte("Not support the Path"))
		}
	})
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("服务器错误：", err.Error())
		panic(err)
	} else {
		fmt.Println("服务器已启动")
	}
}

package cli

import (
	"fmt"
	"io"
	"net/http"
)

// 镜像克隆库
func CloneHaxelib(name string, version string) {
	// 检查库版本是否有效
	version, err := CheckVersion(name, version)
	if err != nil {
		panic(err)
	}
	libzipfile := name + "-" + version + ".zip"
	// 开始克隆
	ossurl := Haxelib_path + "clone/files/3.0/" + libzipfile
	ossret, e := http.Get(ossurl)
	if e != nil {
		panic(e)
	} else {
		defer ossret.Body.Close()
		content, err := io.ReadAll(ossret.Body)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(string(content))
		}
	}
}

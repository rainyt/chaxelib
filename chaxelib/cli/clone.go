package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func CloneHaxelib(name string, version string) {
	version, err := CheckVersion(name, version)
	if err != nil {
		panic(err)
	}
	libzipfile := name + "-" + version + ".zip"
	// 做一个检测
	ossurl := Haxelib_path + "clone/files/3.0/" + libzipfile
	ossret, e := http.Get(ossurl)
	if e != nil {
		panic(e)
	} else {
		defer ossret.Body.Close()
		content, err := ioutil.ReadAll(ossret.Body)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(string(content))
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// 获得库最新的版本
func GetLastVersion(haxelib string) string {
	last := "haxelib/" + haxelib + "/last"
	_, e2 := os.Stat(last)
	if e2 != nil {
		return ""
	}
	b, _ := os.ReadFile(last)
	return string(b)
}

// 查找haxelib库
func FindHaxelib(name string, version string) (string, error) {
	fmt.Println("查询", name, version)
	isUpdata := strings.Contains(version, "+")
	if isUpdata {
		// 需要比对版本
		last := GetLastVersion(name)
		if last == version {
			return "", fmt.Errorf(name + ":" + version + "版本一致，跳过更新")
		}
	}
	savePath := "haxelib/" + name + "/" + version + ".zip"
	_, err := os.Stat(savePath)
	if err == nil {
		return savePath, nil
	}
	return "", fmt.Errorf(name + ":" + version + "不存在")
}

// 识别haxelib是否有效，并缓存版本号库文件
func SaveHaxelib(path string, bytes []byte) error {
	var haxelibJsonPath = findHaxelibJson(path)
	if haxelibJsonPath == "" {
		return fmt.Errorf("不是有效的haxelib库")
	}
	fmt.Println("识别haxelib.json路径：", haxelibJsonPath)
	jsonData, err := os.ReadFile(haxelibJsonPath)
	if err != nil {
		return err
	}
	var _map map[string]any = map[string]any{}
	json.Unmarshal(jsonData, &_map)
	fmt.Println("_map=", _map)
	version := _map["version"].(string)
	haxelibname := _map["name"].(string)
	// 库路径名称/版本号.zip作为存库
	saveName := version + ".zip"
	// 判断存档是否已存在，如果已存在，则必须提交新版本
	savePath := "haxelib/" + haxelibname + "/" + saveName
	os.MkdirAll("haxelib/"+haxelibname, 0777)
	// 并储存最后一个上传的版本记录
	os.WriteFile("haxelib/"+haxelibname+"/last", []byte(version), 0777)
	_, err2 := os.Stat(savePath)
	if err2 == nil {
		return fmt.Errorf(saveName + "存库已经存在，请更新版本号重新更新")
	}
	// 进行存档
	os.WriteFile(savePath, bytes, 0777)
	return nil
}

func findHaxelibJson(dir string) string {
	f, err := os.Open(dir)
	if err != nil {
		return ""
	}
	list, err2 := f.ReadDir(0)
	if err2 != nil {
		return ""
	}
	for _, de := range list {
		fmt.Println(de.Name())
		name := de.Name()
		if name == "haxelib.json" {
			return dir + "/" + name
		} else if de.Type().IsDir() {
			return findHaxelibJson(dir + "/" + name)
		}
	}
	return ""
}

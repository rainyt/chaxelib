package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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

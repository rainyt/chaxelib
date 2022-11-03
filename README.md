## chaxelib
中国（国内）使用`haxelib`安装库的时候，会经常遇到`Timeout`的问题，为了解决这个问题，实现了一个中转服务器，和一个`chaxelib`命令行工具。

## Go环境
请使用`go1.19+`版本

## 编译chaxelib命令行工具
#### Mac
```shell
cd chaxelib
make build-mac
```
#### Window
```shell
cd chaxelib
make build-window
```

## 命令行使用
与`haxelib`保持一致，但只支持`install`命令：
```shell
chaxelib install openfl
chaxelib install openfl:8.0.0
```
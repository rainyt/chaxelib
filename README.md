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

## 安装库
与`haxelib`保持一致，但只支持`install`命令：
```shell
chaxelib install oname:version
```

## 镜像克隆库
如果install太慢，可以等待镜像完成后，再重新install，可通过下述命令，确认克隆情况：
```shell
chaxelib clone name:version
```

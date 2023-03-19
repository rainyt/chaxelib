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

## 本地化服务器
允许构造一个本地化服务器（chaxelib-local）来管理haxelib库，可使用命令：

#### 配置本地化地址
当运行chaxelib-local服务器后，可将启动的端口配置：
```shell
chaxelib local
192.168.1.8:5555
授权码
```

#### 提交库目录
将一个haxelib库上传到本地化服务器中
```shell
chaxelib upload 库目录
```

#### 更新库
将一个haxelib从本地化服务器更新
```shell
chaxelib update 库名:版本号
```
运行将一个hxml文件进行批量更新流程
```shell
chaxelib update ./build.hxml
```

#### 授权码模式
当服务器希望只给拥有授权码的客户端访问，可提供`--pwd`参数：
```shell
chaxelib-local --pwd=ABC123
```
# GMC 框架简介

# 快速开始

# 控制器

# 路由

# 模版引擎

## 模版函数

# Web HTTP 服务器

# API HTTP 服务器

# 数据库

## MySQL 数据库

## SQLITE3 数据库

# 缓存

## Redis 缓存

## File 缓存

## 内存缓存

# I18n国际化

# 中间件

## HTTP 服务器

API 和 Web HTTP服务器工作流程架构图如下，此图可以直观的帮助你快速掌握中间件的使用。

<img src="/doc/images/http-and-api-server-architecture.png" width="960" height="auto"/>  

# 工具包

## gpool

## gmcmap

## sizeutil

## timeutil

## gmclog

## cast

# GMCT 工具链

## 工具链简介

GMC框架，为了降低使用者学习成本和加速开发，提供了开源GMCT工具链，此工具目前可以完成如下功能：

1. 一键初始化一个基于GMC的新Web项目，含有控制器，路由，缓存，数据库，视图等完整的项目代码。

1. 一键初始化一个基于GMC的新API项目，含有控制器，路由，缓存，等完整的项目代码。
    API项目比Web项目轻量化，GMC针对API类型的项目单独定义了API服务器。
    
1. 一键初始化一个基于GMC新API轻量化项目，含有建立一个API服务的基础代码。
    API轻量化项目，适合嵌入到任何现有项目代码中，由几个调用轻松完成一个API服务。

1. 一键打包视图文件为go文件，GMC集成了把试图打包进入二进制的功能，编译项目代码之前，
    只需要在相关目录执行 `gmct tpl --dir ../views` 命令，即可打包视图目录所有视
    图文件为一个go文件。然后正常编译项目，视图就会被打包进入项目编译后的二进制文件中。
    
1. 一键打包网站静态文件为go文件，GMC集成了把网站静态文件，诸如：css，js，font等等，
    打包进入二进制的功能，编译项目代码之前，只需要在相关目录执行 `gmct static --dir ../static` 
    命令，即可打包视图目录所有视图文件为一个go文件。然后正常编译项目，视图就会被打包进入项目编译
    后的二进制文件中。
1. 热编译项目，在项目开发过程中，我们会不断的修改go文件或者视图文件，然后需要手动重新编译，运行，
   才能看到修改后的效果，这占用了同学们不少的时间，为了解决此问题，只需要在项目编译目录执行`gmct run`
   即可，GMCT工具可以侦测到你对项目做的修改，会自动重新编译并运行项目，你修改了项目文件，只要刷新浏览器
   就可以看见最新修改效果。
   
## 工具链使用指南

工具链的使用，需要本机安装好git，配置好了GO环境，go版本1.12及以上即可，并设置 `GOPATH`环境变量。
`PATH`环境变量包含`$GOPATH/bin`目录。环境变量配置不熟悉的同学可以先搜索学习一下配置系统环境变量。

### 安装工具链
gmct工具链的安装有两种方式。一种是直接从源码编译。一种是下载编译好的二进制。

#### 1、从源码编译
此方法需要确保你的网络可以正常的`go get`和`go mod`到各种依赖包，由于特殊原因，如不能下载依赖包，
可以使用通过设置GOPROXY环境变量走代理下载依赖包。设置请参考 [设置GOPROXY环境变量](https://goproxy.io/) 。

然后打开一个命令行，依次执行下面的命令，首次安装，需要下载比较多的依赖包，请同学确保网络正常，耐心等待。

Linux系统：

```shell
export GO111MODULE=on 
git clone https://github.com/snail007/gmct.git
cd gmct && go mod tidy
go install
gmct --help
```

Windows系统：

```shell
set GO111MODULE=on
git clone https://github.com/snail007/gmct.git
cd gmct && go mod tidy
go install
gmct --help
```

#### 2、下载二进制
下载地址：[GMCT工具链](https://github.com/snail007/gmct/releases) ，需要根据你的操作系统平台，下载对应的二进制文件压缩包即可，然后解压得到gmct或者gmct.exe
把它放在`$GOPATH/bin`目录即可,然后打开一个命令行，执行`gmct --help`如果有显示`gmct`的帮助信息，说明安装成功。

### 初始化一个Web项目

GMCT初始化的项目默认使用`go mod`管理依赖,项目路径开始开始于：`$GOPATH/src` 。
初始化项目只需要一个参数`--pkg`就是项目路径。

操作步骤，顺序执行下面命令：

```shell
gmct new web --pkg foo.com/foo/myweb
cd $GOPATH/src/foo.com/foo/myweb
gmct run
```

打开浏览器访问：http://127.0.0.1:7080 , 就可以看见新建的web项目运行效果。

### 初始化一个API项目

GMCT初始化的项目默认使用`go mod`管理依赖,项目路径开始开始于：`$GOPATH/src` 。
初始化项目只需要一个参数`--pkg`就是项目路径。

操作步骤，顺序执行下面命令：

```shell
gmct new api --pkg foo.com/foo/myapi
cd $GOPATH/src/foo.com/foo/myapi
gmct run
```

打开浏览器访问：http://127.0.0.1:7081 , 就可以看见新建的API项目运行效果。

### 初始化一个API轻量级项目

GMCT初始化的项目默认使用`go mod`管理依赖,项目路径开始开始于：`$GOPATH/src` 。
初始化项目只需要一个参数`--pkg`就是项目路径。

操作步骤，顺序执行下面命令：

```shell
gmct new api-simple --pkg foo.com/foo/myapi0
cd $GOPATH/src/foo.com/foo/myapi0
gmct run
```

打开浏览器访问：http://127.0.0.1:7082 , 就可以看见新建的API轻量级项目运行效果。
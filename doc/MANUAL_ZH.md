# GMC 框架简介

# 快速开始

# 控制器

# 路由

# 模版引擎

## 模版函数

# Web 服务器

# API 服务器

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

# 平滑重启/热升级

此功能只适用于Linux平台系统，不适用于Windows系统。

当我们的程序部署到线上的时候，就面临一个平滑重启/平滑升级的问题，也就是在不中断当前已有连接的保证服务一直可用的情况下，进行程序的重启升级。

通过gmc.APP启动的Web和API服务都支持平滑重启/平滑升级，使用非常简单，当你需要重启的时候，使用pkill或者kill命令给程序发送`USR2`信号即可。

比如：

```shell
pkill -USR2 website
kill -USR2 11297
```

本示例中 `website` 是程序名称，`11297`是程序的`pid`。两种方式都可以，自己看习惯选择。

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
    命令，即可打包静态资源目录所有文件为一个go文件。然后正常编译项目，静态资源文件就会被打包进入项目编译
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

### 打包视图文件

GMC视图模块支持打包视图文件到编译的二进制程序中。由于打包功能和项目目录结构有关，所以这里假设目录结构是gmct生成的web项目目录结构。

目录结构如下：

```text
new_web/
├── conf
├── controller
├── initialize
├── router
├── static
└── views
```

完成此功能需要下面几个步骤：

```shell
cd initialize
gmct tpl --dir ../views
```

执行了命令后，会发现目录initialize中多了一个前缀是`gmc_templates_bindata_`的go文件，
比如：`gmc_templates_bindata_2630881503983182670.go`。此文件中有init方法，
会在`initialize`包被引用的时候自动执行，把视图文件二进制数据注入GMC视图模块中。

当你go build编译了你的项目后，避免后面开发，运行代码一直使用这个go文件里面的视图数据，
可以在目录initialize执行：`gmct tpl --clean` 可以安全的清理上面生成的go文件。

### 打包静态文件

GMC的HTTP静态文件模块支持打包静态文件到编译的二进制程序中。由于打包功能和项目目录结构有关，
所以这里假设目录结构是gmct生成的web项目目录结构。

目录结构如下：

```text
new_web/
├── conf
├── controller
├── initialize
├── router
├── static
└── views
```

完成此功能需要下面几个步骤：

```shell
cd initialize
gmct static --dir ../static
```

执行了命令后，会发现目录initialize中多了一个前缀是`gmc_static_bindata_`的go文件，
比如：`gmc_static_bindata_1780615241186372497.go`。此文件中有init方法，
会在`initialize`包被引用的时候自动执行，把视图文件二进制数据注入GMC的HTTP静态文件服务模块中。

当你go build编译了你的项目后，避免后面开发，运行代码一直使用这个go文件里面的静态文件数据，
可以在目录initialize执行：`gmct static --clean` 可以安全的清理上面生成的go文件。

另外，静态文件打包进入二进制以后，查找静态文件的顺序是：  
1、查找二进制数据里面是否有该文件。  
2、查找静态目录static是否有该文件。  

### 热编译项目

在项目开发过程中，我们会不断的修改go文件或者视图文件，然后需要手动重新编译，运行，
才能看到修改后的效果，这占用了同学们不少的时间，为了解决此问题，只需要在项目编译目录执行`gmct run`
即可，GMCT工具可以侦测到你对项目做的修改，会自动重新编译并运行项目，你修改了项目文件，只要刷新浏览器
就可以看见最新修改效果。

执行`gmct run`后会在当前目录生成一个名称为`gmcrun.toml`的配置文件，你可以通过修改该文件，定制`gmct run`的编译行为。

默认情况下配置如下：

```toml
[build]
# ${DIR} is a placeholder presents current dir absolute path, no slash in the end.
# you can using it in monitor_dirs, include_files, exclude_files, exclude_dirs.
monitor_dirs=["."]
args=["-ldflags","-s -w"]
env=["CGO_ENABLED=1","GO111MODULE=on"]
include_exts=[".go",".html",".htm",".tpl",".toml",".ini",".conf",".yaml"]
include_files=[]
exclude_files=["gmcrun.toml"]
exclude_dirs=["vendor"]
```

配置说明：

1. `monitor_dirs` 监控目录，默认是当前目录，可以指定多个，数组形式。
1. `args` 额外传递给 `go build` 的参数，数组，多个参数分开写。
1. `env` 设置`go build`执行时的环境变量，可以指定多个，数组形式。
1. `include_exts` 监控的文件后缀，只监控此后缀文件的改动。可以指定多个，数组形式。
1. `include_files` 设置额外监控的文件，支持相对路径和绝对路径，可以指定多个，数组形式。
1. `exclude_files` 设置额外不监控的文件，支持相对路径和绝对路径，可以指定多个，数组形式。
1. `exclude_dirs`  设置额外不监控的目录，支持相对路径和绝对路径，可以指定多个，数组形式。

在 `monitor_dirs`, `include_files`, `exclude_files`, `exclude_dirs`中可以使用
变量`{DIR}`代表当前目录的绝对路径, 末尾没有`/`。

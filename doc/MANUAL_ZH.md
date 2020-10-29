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

## 缓存介绍

GMC缓存Cache支持Redis、File、内存缓存三种类型，为适应不同的业务场景，开发者也可以自己实现`gmccore.Cache`接口，
然后通过`gmccachehelper.AddCacheU(id,cache)`注册自己的缓存，然后就可以通过`gmc.Cache.Cache(id)`获取自己注册的缓存对象。

如果要在`gmc`项目里面使用缓存，需要修改配置文件`app.toml`里面的`[cache]`部分，首先设置默认缓存类型,比如使用redis `default="redis"`，
然后需要修改对应缓存驱动`[[redis]]`部分的启用`enable=true`。每个驱动类型的缓存都可以配置多个，每个的id必须唯一，id是"default"的将作为默认使用。

比如：

如果redis配置了多个，那么`gmc.Cache.Redis()`获取的就是id为default的那个。

## 缓存配置说明

```shell
[cache]
default="redis" //设置默认生效缓存配置项，比如项目默认为redis缓存生效
[[cache.redis]] //redis配置项
[[cache.file]]  //file配置项
[[cache.memory]]//内存缓存配置项，其中cleanupinterval为自动垃圾收集时间单位是second
```

通过gmc.APP启动的API或者Web服务，使用配置文件配置缓存，当你在配置文件app.toml启用了缓存，
那么可以通过gmc.Cache.Cache()使用缓存。

### Redis缓存

基于redigo@v2.0.0实现，支持redis官方主流方法调用，可以适用绝大部分业务场景

```shell
[[cache.redis]]
debug=true       //是否启用调试
enable=true      //开启redis缓存
id="default"     //缓存池ID
address=":6379"  //redis客户端链接地址
prefix=""
password=""
timeout=10       //等待连接池分配连接的最大时长（毫秒），超过这个时长还没可用的连接则发生
dbnum=0          //连接Redis时的 DB 编号，默认是0.
maxidle=10       //连接池中最多可空闲maxIdle个连接
maxactive=30     //连接池支持的最大连接数
idletimeout=300  //一个连接idle状态的最大时长（毫秒），超时则被释放
maxconnlifetime=3600 //一个连接的生命时长（毫秒），超时而且没被使用则被释放
wait=true
```

### Memory缓存

cache.go是轻量级的go缓存实现，shard.go没有使用 go 的”hash/fnv”中的 hash.Hash 函数，使用的是djb3算法，在大块文件存储效率比标准cache提升约1倍
配置信息如下所示，cleanupinterval信息表示 GC 的时间，表示每隔 30s 会进行一次过期清理,id是默认连接池id,enable表示是否开启缓存，默认关闭。

```shell
[[cache.memory]]
enable=false
id="default"
cleanupinterval=30
```

### File缓存

配置信息如下所示，配置 dir 表示缓存的文件目录，cleanupinterval信息表示 GC 的时间，表示每隔 30s 会进行一次过期清理,id是默认连接池id,enable表示是否开启缓存，默认关闭。

```shell
[[cache.file]]
enable=false
id="default"
dir="{tmp}"
cleanupinterval=30
```

## 单独使用缓存模块

当然`gmccache`包也可以单独使用，不依赖gmc框架,自己实例化配置对象，初始化缓存配置，使用方法示例如下：

```golang
package main

import (
	"time"

	"github.com/snail007/gmc"
)

func main() {
	cfg := gmc.New.Config()
	cfg.SetConfigFile("../../app/app.toml")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// Init only using [cache] section in app.toml
	gmc.Cache.Init(cfg)

	// cache default is redis in app.toml
	// so gmc.Cache() equal to  gmc.Redis()
	// we can connect to multiple cache drivers at same time, id is the unique name of driver
	// gmc.Cache(id) to load `id` named default driver.
	c := gmc.Cache.Cache()
	c.Set("test", "aaa", time.Second)
	c.Get("test")
}
```

# I18n国际化

GMC国际化文件内容是由多行的key=value组成。文件名称使用标准的HTTP头部`Accept-Language`中的格式，比如：zh-CN，en-US。
后缀是`.toml`。

value，可以理解为，value是fmt.Printf里第一个参数，里面可以写支持的占位符。
模版输出的时候，可以格式化输出。

比如：

{{printf (trs .Lang "foo2") 100 }}

printf，tr，string都是模版函数。
1. printf是格式化字符串作用和fmt.Printf一样。
1. tr是翻译函数，第一个参数是固定的，它对应控制器里面的this.Lang。
1. 由于tr返回的是template.HTML类型数据，可以用string函数转为string类型。

比如中文翻译文件 `zh-CN.toml` 内容：

```text
hello="你好"
foo="任意内容"
foo1="你的年龄是%d岁"
foo2="还支持html哟，<b>加粗</b>"
```

比如英文翻译文件 `en-US.toml` 内容：

```text
hello="Hello"
foo="Bar"
```

hello，foo 其实就是一个key，用来寻找对应的value，也就对应的内容。

以上就是GMC的国际化的规则和原理。

国际化功能，在使用之前需要在配置app.toml里的`[i18n]`模块启用国际化功能，设置`enable=true`，
设置国际化文件目录`dir="i18n"`，默认是程序平级目录`i18n`，里面存放上面说的国际化翻译文件。
当有请求的时候，控制器方法在被调用之前，控制器的`this.Lang`成员变量，已经根据访问者浏览器的
Accept-Language被初始化好。在控制器中可以通过this.Tr(key , hint string)，寻找多个国际
化语言文件中最优匹配的语言，然后在该语言文件中进行寻找匹配key翻译结果。

国际化相关配置如下：

```toml
[i18n]
enable=false
dir="i18n"
default="zh-cn"
```

`default`是默认语言，如果用户HTTP请求头部的语言，国际化模块找不到与之匹配的语言，那么就使用
这个`default`设置的语言进行翻译。

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

# GMC 框架简介

# 快速开始

GMC快速开始的最好方式是通过GMCT工具链，所以首先同学们要安装好`GMCT工具链`，可以参考本手册[GMCT工具链](#GMCT-工具链)章节进行安装。

快速建立一个Web项目，请参考：[初始化一个Web项目](#初始化一个Web项目)

假设我们用`GMCT工具链`新建了一个项目，项目路径是：`$GOPATH/src/gmcdemo`，默认情况下，项目目录文件结构如下：

```text
./gmcdemo
├── conf
│   └── app.toml
├── controller
│   └── demo.go
├── gmcrun.toml
├── go.mod
├── go.sum
├── initialize
│   └── initialize.go
├── main.go
├── router
│   └── router.go
├── static
│   └── jquery.js
└── views
    └── welcome.html
```

项目文件说明：
1. `conf/app.toml` 是项目配置文件，几乎所有Web功能都在这里配置，比如：web监听设置，session，缓存，数据库等等，是最常用的文件。
1. `controller/demo.go` demo控制器文件，默认为同学们建立好了一个示例控制器，快打开它看看吧。
1. `gmcrun.toml`这个文件不是gmc项目的，是gmct工具的配置文件，因为我用了gmct run运行项目，所以会有此文件，具体使用可以参考[GMCT工具链](#GMCT-工具链)。
1. `go.mod, go.sum` 这两个是项目go mod包依赖管理文件。
1. `initialize/initialize.go` 是项目启动时，初始化的操作都单独在这一个文件里面，而不是所有东西都写在入口件里面，这种方法是GMC推荐的项目结构组织方式。
    示例项目默认情况下，这个文件的初始化方法里面调用了初始化了路由的配置。
1. `router/router.go` 路由的配置，项目单独最为一个文件独立出来，这样可以很清晰的管理我们的路由。
1. `main.go` 这个是项目的入口文件，主要是通过新建的`GMC APP`对象启动了`Web服务`，然后在服务初始化的时候，调用了上面`initialize`初始化包的初始化方法。
1. `static/jquery.js` 是项目默认主页引入的js文件。
1. `views/welcome.html` 是项目默认主页使用的视图文件，在上面说的`Demo控制器`里面你可以看见它的身影哟，我们访问 `http://127.0.0.1:7080/` 看见的页面，就是这个视图渲染的结果。

# 控制器

## 编写控制器

编写一个自己的GMC控制器，需要导入包`github.com/snail007/gmc`,定义自己的`struct`继承`gmc.Controller`即可。

下面示例代码，实现了一个简单的控制器，定义了一个`Hello`方法。

```go
package controller

import (
	"github.com/snail007/gmc"
)

type Demo struct {
	gmc.Controller
}

func (this *Demo) Hello() {
	this.Write("fmt.Println(\"Hello GMC!\")")
}
```

## 控制器规则

1. 控制器名称无限制。
1. 后缀是两个连续英文半角下划线`__`的控制器方法，在路由里面绑定控制器的时候会被忽略。
1. 控制器方法名称不能包含以下名称，它们是GMC完成框架功能的方法，同学们的控制器里面不能使用这些名称。
   `MethodCallPre__()`，`MethodCallPost__()`，`Stop()`，`Die()`，`Tr()`，`SessionStart()`，
   `SessionDestroy()`，`Write()`，`StopE()`。
1. 方法名称是`Before__()`的方法，是控制器的构造方法，不需要可以不定义，在被访问控制器方法执行之前被调用，
    可以调用`this.Stop()`阻止被访问`控制器方法`的调用，但不能阻止`After__()`控制器析构方法的调用。
    可以通过`this.Die()`阻止被访问`控制器方法`和`After__()`调用。
1. 方法名称是`After__()`的方法，是控制器的析构方法，不需要可以不定义，在被访问控制器方法执行之后被调用。
1. 控制器成员不能包含以下名称，它们是GMC完成框架功能用的，同学们的控制器里面的成员名称不能使用这些名称。
   `Response`，`Request`，`Param`，`Session`，`Tpl`，`SessionStore`，`Router`，`Config`，
   `Cookie`，`Ctx`，`View`，`Lang`，`Logger`，这些成员是十分有用的，我们经常会使用到它们，下面会对它们一一介绍。

## 获取输入

获取输入可以通过`this.Request`对象，它是原生标准的`*http.Request`对象。通过可以访问`GET，POST，COOKIE，上传文件`等数据。
另外还能通过`this.Ctx`进行一些便捷的输入操作。

## 输出内容

输出可以通过`this.Response`对象，它是原生标准的`http.ResponseWriter`对象。另外你还可以通过`this.Write()`，输出绝大部分数据内容，
它能自动识别各种数据类型，自动转换后输出到浏览器。

## Session操作

GMC的session和其它同类实现有很大区别，其它同类实现，要么你全局开启，要么全局关闭，这样带来的缺点是，无论你访问哪个控制器的方法，
框架都会进行初始化session的各种操作，造成不必要的，很大的性能消耗。GMC避免了这个缺点，借鉴了PHP的SESSION实现机制，在需要操作session数据的地方，
手动调用this.SessionStart()开启session后，才可以访问session数据，这样把性能的消耗降低到最低。另外需要销毁session数据的时候，
调用this.SessionDestroy()即可。

控制里面可以通过：
- `this.SessionStart()`开启session。
- `this.SessionDestroy()`销毁session数据。
- `this.Session`访问或者设置session数据，必须`this.SessionStart()`开启session后才能使用这个对象，不然这个对象是nil的。

## Cookie操作

操作Cookie，不仅可以通过标准的`this.Request`和`http.ResponseWriter`，
还能通过GMC提供的`this.Cookie`进行Cookie的`Set，Get，Remove`操作。

# 路由

GMC的路由支持以下功能：

- 绑定控制器，自动识别方法名称，方法名称在URL中一律使用小写。
- 绑定控制器方法。
- 绑定`http.Handler`。
- 自定义控制器在URL里面的路径。
- 路由分组。
- 每个控制器可以自定义自己的URL后缀。
- 每个路由组可以自定义自己的URL后缀。

## 获取路由

通过示例项目，同学们看到，通过新建的Web服务器或者API服务器对象的Router()方法就可以获取对应的路由对象，然后就可以它进行路由的各种配置了。

示例项目的router包里面，通过路由对象进行了简单路由绑定。

## 绑定路由

路由绑定有以下几种方式。

### 绑定控制器

```go
//...
r.Controller("/demo", new(controller.Demo))
//...
```

- 第一个参数`/demo`是控制器在URL中的路径。
- 访问控制器方法的全路径是：控制器的路径+全部小写的方法名称。
- URL访问`Demo`控制器的`Hello`方法的URL路径就是：`/demo/hello`。
- 还能传递第三个参数，作为此控制器可以在URL中访问的方法的URL后缀。

### 绑定控制器方法

```go
//...
r.ControllerMethod("/",new(controller.Demo),"Index__")
//...
```

`ControllerMethod`，可以绑定控制器任意的公有方法，即使后缀是`__`也能绑定。

### 绑定`http.Handler`

下面示例中使用了各种方法绑定一个`http.Handler` `myHanlder` 到URL路径 `/hello` 上，
`HandlerAny`相当于一次性绑定`GET,POST,DELETE,OPTIONS,PUT,PATCH,HEAD`。

```go
//...
r.Handler("GET","/hello",myHanlder)
r.Handler("POST","/hello",myHanlder)
r.Handler("DELETE","/hello",myHanlder)
r.Handler("OPTIONS","/hello",myHanlder)
r.Handler("PUT","/hello",myHanlder)
r.Handler("PATCH","/hello",myHanlder)
r.Handler("HEAD","/hello",myHanlder)
r.HandlerAny("/hello",myHanlder)
//...
```

### 绑定`gmc.Handle`

下面示例中使用了各种方法绑定一个`gmc.Handle` `Hello` 到URL路径 `/hello` 上，
`HandleAny`相当于一次性绑定`GET,POST,DELETE,OPTIONS,PUT,PATCH,HEAD`。

```go
//...
func Hello(w gmc.W, r gmc.R, ps gmc.P) {
    w.Write([]byte("Hello"))
}
r.GET("/hello",Hello)
r.POST("/hello",Hello)
r.DELETE("/hello",Hello)
r.OPTIONS("/hello",Hello)
r.PUT("/hello",Hello)
r.PATCH("/hello",Hello)
r.HEAD("/hello",Hello)
r.Handle("GET","/hello",Hello)
r.HandleAny("/hello",Hello)
//...
```

## 定义参数

示例代码如下：
         
 ```go
 //...
r.GET("/user/:name",Hello)
 //...
 ```
`:name`是一个占位符参数，此位置的值：
- 在控制器里面可以通过 `this.Param.ByName("name")` 获取。
- 在绑定的Handle里面可以通过 `gmc.P.ByName("name")` 获取。
- 在 `http.Handler` 里面可以通过：
    ```go
    params := gmc.ParamsFromContext(r.Context())
    params.ByName("name")
    ```

## 参数规则

定义一个或多个参数：

```text
URL绑定路径: /user/:user

 /user/gordon              匹配
 /user/you                 匹配
 /user/gordon/profile      不匹配
 /user/                    不匹配
```

匹配所有路径：

```text
URL绑定路径: /src/*filepath

 /src/                     匹配
 /src/somefile.go          匹配
 /src/subdir/somefile.go   匹配
```


# 模版引擎

GMC模版引擎，是对`text/template`的包装，增强了其功能。模版的基本使用方式，语法都和`text/template`一样。

## 模版配置

模版的配置在app.toml中，默认如下：

```toml
[template]
dir="views"
ext=".html"
delimiterleft="{{"
delimiterright="}}"
```

- dir 是存放视图模版文件的目录。
- ext 是模版文件的后缀，只有此后缀的文件才会被解析。
- delimiterleft 模版中语法块的左符号。
- delimiterright 模版中语法块的右符号。

## 包含模版

假设模版目录是views，它里面的文件结构如下：

```text
views/
├── common
│   └── head.html
└── user
    └── list.html
```

user/list.html 内容如下：

```text
{{template "common/head" . }}
```
这样就在`user/list.html`中包含了模版文件`common/head.html`，不需要写`.html`后缀，`.` 的作用是把当前模版的数据
传递给`common/head.html`。如果不需要传递数据，可以省略`.`，即 `{{template "common/head"}}`。

## 模版布局

布局模版是开发网站经常用到的，我们渲染很多页面的时候，它们的基础布局框架是一样的，只是某一部分不同，这个时候就需要用到布局模版支持，
原理是渲染某个模版的的时候，把它渲染的结果放到"布局模版"的指定位置。

假设视图文件夹是views，它的文件结构如下：

```text
views/
├── layout
│   └── page.html
└── user
    └── profile.html
```

`page.html`内容如下：

```text
`{{.GMC_LAYOUT_CONTENT}}`
```

当我们在控制器里面渲染`profile.html`，代码如下：

```go
this.View.Layout("layout/page").Render("profile")
```

渲染过程是，首先使用视图数据渲染`profile.html`,然后把渲染的结果给模版变量 `GMC_LAYOUT_CONTENT`，
然后用这个数据渲染`layout/page.html`，把渲染结果输出到浏览器。

## 模版函数

### 内置函数

`text/template` 内置函数如下：

`and`, `call`, `html`, `index`, `slice`, `js`, `len`, `not`, `or`, `print`, `printf`, `println`, `urlquery`

### 内置比较函数

`text/template` 内置比较函数如下：

`eq`, `ne`, `lt`, `le`, `gt`, `ge`

### GMC 模版函数

GMC为了方面模版的使用，定义了很多有用的函数。

1. `tr` 用在国际化中，第一个参数固定是 `.Lang`, 第二个参数是在国际化配置文件中定义的key。第三个参数是提示信息，这个信息不会输出。返回的数据类型是：template.HTML
    
    示例:
    
    `{{tr .Lang "key001" "提示信息，这个是什么"}}`

1. `trs` 和`tr`功能一样，但是返回的数据类型是string。

1. `string` 只有一个参数，转换数据为字符串。 

1. `tohtml` 只有一个参数，转换数据为 template.HTML 类型.

1. `val` 获取模版变量内容，如果变量不存在，返回空""，而不是"&lt;no value&gt;，可以用于安全输出模版变量。
        两个参数，第一个固定是 `.` 第二个是变量名称。
            
    示例:
    
    `{{val . "name"}}`

另外基于引入了 [sprig](https://github.com/masterminds/sprig) 定义的大量有用模版方法，
gmc做了精简，精简后的全部介绍在：[template/sprig/docs](/http/template/sprig/docs)，文档里的方法在gmc的模版里面可以直接使用。

### GMC 模版变量

在模版中，GMC 为模版准备好了以下变量，可以直接在任何模版中使用，它们的含义如下：

1. `.Lang` 它是客户端发送的http头部 Accept-Language 被解析后的标准格式，比如：zh-CN。
    模版国际化函数`tr`会用到它。
1. `.G` GET数据，是一个 `map[string][string]`，`{{.G.key}}` key是GET数据的中的键名称。
1. `.P` POST数据，是一个 `map[string][string]`，`{{.P.key}}` key是POST数据的中的键名称。
1. `.S` SESION数据，是一个 `map[string][string]`，`{{.S.key}}` key是SESSION数据的中的键名称。
    提示：只有在控制器里面调用 `SessionStart()` 开启了session，才有效。
1. `.C` POST数据，是一个 `map[string][string]`，`{{.C.key}}` key是COOKIE数据的中的键名称。
1. `.U` 当前请求URL的信息， 是一个 `map[string]string`，它的详细内容如下：
    
    ```golang
    //u0 是一个url对象。
    u["HOST"] = u0.Host
    u["HOSTNAME"]=u0.Hostname()
    u["PORT"]=u0.Port()
    u["PATH"] = u0.Path
    u["FRAGMENT"] = u0.Fragment
    u["OPAQUE"] = u0.Opaque
    u["RAW_PATH"] = u0.RawPath
    u["RAW_QUERY"] = u0.RawQuery
    u["SCHEME"] = u0.Scheme
    u["USER"] = u0.User.Username()
    u["PASSWORD"],_ = u0.User.Password()
    u["URI"]=u0.RequestURI()
    u["URL"]=u0.String()
    ```
   
   提示:
   
   URL的完整格式： `[scheme:][//[userinfo@]host][/]path[?query][#fragment]`  
 
## 获取渲染内容

默认情况下，控制器里面使用`this.View,Render()`渲染模版，结果会被输出到浏览器。
如果我们想获取渲染结果，而不是输出到浏览器，可以使用`this.View.RenderR()`渲染模版，它会把渲染结果返回。

# Web 服务器

GMC把Web开发有网页页面的网站，比如新闻站，管理后台这种归为Web服务， 对于这种项目，GMC专门设计了`Web服务`。

GMC Web服务默认包含了模版，静态文件，数据库，session，缓存，国际化这些模块，可以在配置里面开启关闭，开箱即用。

通过GMCT工具链生成的Web项目，默认配置文件位于：conf/app.toml，app.toml是项目的核心，几乎所有gmc功能都是这里配置。

Web服务器在gmc中对应的是：`gmc.HTTPWebServer`，我们也可以看见，在项目的main文件，里面通过`gmc.APP`对象启动一个Web服务。

还可以在服务初始化前后执行一些自己的初始化等操作。

gmc生成的web项目，主文件内容如下：

```go
package main

import (
	"github.com/snail007/gmc"
	"mygmcweb/initialize"
)

func main() {

	// 1. create an default app to run.
	app := gmc.New.AppDefault()

	// 2. add a http server service to app.
	app.AddService(gmc.ServiceItem{
		Service: gmc.New.HTTPServer(),
		AfterInit: func(s *gmc.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(s.Service.(*gmc.HTTPServer))
			return
		},
	})

	// 3. run the app
	if e := gmc.StackE(app.Run());e!=""{
		app.Logger().Panic(e)
	}
}

```

主要做了三件事情：

1. 创建一个默认gmc.APP对象，用来管理程序和服务的整个生命周期，默认APP对象会使用`conf/app.toml`作为配置文件。
1. 创建一个gmc.HTTPWebServer服务对象，添加到APP的服务列表中，并定义了服务初始化后执行的一些自己的初始化。
1.  启动APP，并捕获错误信息，然后输出错误，APP启动后Run()会阻塞，直到手动关闭APP或者发生异常。

# API 服务器

同学们经常会为APP或者第三方写一些数据的HTTP API接口，就是对外提供数据接口没有页面的这种API服务，GMD专门设计了`API服务`。

API服务相对于Web服务更加精简，性能更好，去掉了模版，静态文件，国际化，session这些模块。

## API项目

通过GMCT工具链生成的API项目，默认配置文件位于：conf/app.toml，app.toml是项目的核心，几乎所有gmc功能都是这里配置，
此项目适合最为独立完整的一个API项目。

项目主文件如下：

```go
package main

import (
	"github.com/snail007/gmc"
	"mygmcapi/handlers"
)

func main() {
	// 1. create app
	app := gmc.New.App()
	// 2. parse config file
	cfg, err := gmc.New.ConfigFile("conf/app.toml")
	if err != nil {
		app.Logger().Error(err)
	}
	// 3. create api server
	api, err := gmc.New.APIServerDefault(cfg)
	if err != nil {
		app.Logger().Error(err)
	}
	//4. init db, cache, handlers
	// int db
	gmc.DB.Init(cfg)
	// init cache
	gmc.Cache.Init(cfg)
	// init api handlers
	handlers.Init(api)
	// 5. add service
	app.AddService(gmc.ServiceItem{
		Service: api,
	})
	// 6. run app
	if e := gmc.StackE(app.Run());e!=""{
		app.Logger().Panic(e)
	}
}

```
主要做了以下事情：

1. 创建一个默认gmc.APP对象，用来管理程序和服务的整个生命周期。
1. 创建一个配置对象，设置使用`conf/app.toml`作为配置文件。
1. 创建一个gmc.APIServer服务对象，并使用上面的配置对象。
1. 初始化后执行的一些用到的功能，数据库，缓存，路由。
1. 把API对象添加到APP管理服务列表中。
1. 启动APP，并捕获错误信息，然后输出错误，APP启动后Run()会阻塞，直到手动关闭APP或者发生异常。


## 极简API

如果你就只想对外提供一个或几个简单接口，什么多余的功能都不需要，甚至不需要任何配置文件，
仅仅是监听端口管理路由，你只需要注册路由处理函数即可。通过GMC生成的`simple API`是一个很好示例。

Simple API示例：

```go
package main

import (
	"github.com/snail007/gmc"
)

func main() {

	api := gmc.New.APIServer(":7082")
	api.API("/", func(c gmc.C) {
 		c.Write(gmc.M{
 			"code":0,
 			"message":"Hello GMC!",
 			"data":nil,
		})
	})

	app := gmc.New.App()
	app.AddService(gmc.ServiceItem{
		Service: api,
		BeforeInit: func(s gmc.Service, cfg *gmc.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	if e := gmc.StackE(app.Run());e!=""{
		app.Logger().Panic(e)
	}
}
```

上面的示例使用APP对API服务进行管理，我们甚至APP都可以不用。

代码如下：

```go
package main

import (
	"github.com/snail007/gmc"
)

func main() {

	api := gmc.New.APIServer(":7082").Ext(".json")
	api.API("/hello", func(c gmc.C) {
 		c.Write(gmc.M{
 			"code":0,
 			"message":"Hello GMC!",
 			"data":nil,
		})
	})
	if e := gmc.StackE(api.Run());e!=""{
		panic(e)
	}
}
```

主要做了以下事情：

1. 创建一个gmc.APIServer服务对象，监听7082端口，设置url后缀为`.json`。
1. 注册了一个url路径`/hello.json`，输出一个json数据。
1. 启动API，并捕获错误信息，然后输出错误，APP启动后Run()会阻塞，直到手动关闭APP或者发生异常。

# 数据库

GMC默认提供MYSQL和SQLite3数据库的便捷访问支持，其它类型数据库请同学们使用成熟第三方数据库操作库。

GMC 数据库操作支持：
1. 多数据源支持.
1. SQLite3支持加密.
1. 方便的链式查询和更新。
1. 映射结果集到结构体。
1. 事务支持。

## MySQL 数据库

MySQL数据库支持连接池，各种超时配置，详细配置信息在app.toml中的`[database]`部分。

```toml
[database]
default="mysql"

[[database.mysql]]
enable=true
id="default"
host="127.0.0.1"
port="3306"
username="root"
password="admin"
database="test"
prefix=""
prefix_sql_holder="__PREFIX__"
charset="utf8"
collate="utf8_general_ci"
maxidle=30
maxconns=200
timeout=15000
readtimeout=15000
writetimeout=15000
maxlifetimeseconds=1800

[[database.mysql]]
enable=false
id="news"
host="127.0.0.1"
port="3306"
username="root"
password="admin"
database="test"
prefix=""
prefix_sql_holder="__PREFIX__"
charset="utf8"
collate="utf8_general_ci"
maxidle=30
maxconns=200
timeout=3000
readtimeout=5000
writetimeout=5000
maxlifetimeseconds=1800
```

可以发现`database.mysql`是一个数组，同学们可以配置个，只要保证id唯一就可以。
`id`是在代码里面使用这个配置数据库对象的标识。

### 使用数据库

在API或者Web项目里面，app启动之后，也就是数据库被初始化成功后，在代码里面可以通过
`gmc.DB.DB(id ...string)`或`gmc.DB.MySQL(id ...string)`获取数据库操作对象。

示例代码：

```go

package main

import (
	"github.com/snail007/gmc"
)

func main() {
	cfg := gmc.New.Config()
	cfg.SetConfigFile("../../app/app.toml")
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// Init only using [database] section in app.toml
	gmc.DB.Init(cfg)

	// database default is mysql in app.toml
	// so gmc.DB.DB() equal to  gmc.DB.MySQL()
	// we can connect to multiple cache drivers at same time, id is the unique name of driver
	// gmc.DB.DB(id) to load `id` named default driver.
	db := gmc.DB.DB().(*gmc.MySQL)
	//do something with db
	db.AR()
}
```

由于数据库说明比较多，详细的请 [参考这里](/db/mysql/README.md)

## SQLITE3 数据库

SQLITE3 数据库各种配置和加密，详细配置信息在app.toml中的`[database]`部分。

```toml
[[database.sqlite3]]
enable=false
id="default"
database="test.db"
# if password is not empty , database will be encrypted.
password=""
prefix=""
prefix_sql_holder="__PREFIX__"
# syncmode 0:OFF, 1:NORMAL, 2:FULL, 3:EXTRA
syncmode=0
# openmode ro,rw,rwc,memory
openmode="rw"
# shared,private
cachemode="shared"
```

### 使用数据库

在API或者Web项目里面，app启动之后，也就是数据库被初始化成功后，在代码里面可以通过
`gmc.DB.DB(id ...string)`或`gmc.DB.SQLite3(id ...string)`获取数据库操作对象。


由于数据库说明比较多，详细的请 [参考这里](/db/sqlite3/README.md)

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

GMC的Web和API服务器都支持中间件，当现有功能无法完成你的需求，你可以通过注册中间件，完成各种功能，比如：`权限认证`，`日志记录`，`数据埋点`，`修改请求`等等。
从GMC收到用户请求开始，到控制器方法或者handler被真正被执行这期间，按着执行的顺序和时机，中间件分为四种
名字分别为：`Middleware0`，`Middleware1`，`Middleware2`，`Middleware3`，每种中间件都可以注册多个，
它们会按添加的顺序被依次执行，但是当某种的某个中间件返回FALSE的时候，位于这个中间件后的此类中间件都将被跳过执行。

一个中间件就是一个function。

API服务中它的定义是`func(ctx *gmcrouter.Ctx, server *APIServer) (isStop bool)`。

Web服务中它的定义是`func(ctx *gmcrouter.Ctx, server *HTTPServer) (isStop bool)`。

API 和 Web HTTP服务器工作流程架构图如下，它们执行的顺序和时机，此图可以直观的帮助你快速掌握中间件的使用。

下图是适用于Web和API服务，把下面的controller理解为API中注册的 handler，就是API服务中间件的流程图。

图中的STOP对应的就是中间件function返回的值是true的时候。

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

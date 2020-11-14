# GMC INTRODUCTION

<img align="right" src="https://github.com/snail007/gmc/blob/master/doc/images/logo2.png?raw=true" width="200" height="auto"/>  

GMC provides easy and fast Web and API development, and aims at high performance, high productivity, and lightweight. 

GMC provides a powerful supporting tool chain, let you one-command to generate a base project, get rid of the manual 
creation of various folders, files worry.

Integrate hot compilations, smooth updates/restarts, package static files to binaries, package views files to binaries, 
package i18n files to binaries, and more to keep your Web applications highly portable and maintainable.
Rich documentation and detailed explanations make it easy to use.

# QUICK START

The best way to start GMC quickly is through the GMCT tool chain, so first students to install`GMCT tool chain`, 
you can refer to this manual [GMCT Tool Chain](#gmct-tool-chain) chapter to install.

To quickly create a Web project, see: [NEW WEB PROJECT](#new-web-project)

We create a new project using the`GMCT tool chain`. The project path is:`$GOPATH/src/gmcdemo`. 

By default, the project directory file structure is as follows:

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

Project structure:
1.`conf/app.toml`is a project configuration file, where almost all Web functions are configured, such as: Web Server 
settings, session, cache, database, etc. It is the most commonly used file.

1.`controller/demo.go`demo controller file, default for students to set up a sample controller, open it to have a look.

1.`gmcrun.toml`is not a GMC project, but a configuration file of GMCT tool. Since I used GMCT Run to run the project, 
there is this file, please refer to [GMCT tool chain](#gmct-tool-chain) for specific use.

1.`go.mod, go.sum`are the project go mod package dependency management files.

1.`initialize/initialize.go`is the method recommended by GMC to organize the project structure by initializing the 
operation in this file instead of everything written in the entry.
By default, the initialization method for this file calls the configuration that initializes the routing.

1. Router/Router. go`configuration, the project is separate as a single file, so that we can manage our routes clearly.

1.`main.go`is the entry file for the project. It starts`Web service`with the newly created`GMC APP`object, and then
 calls the initialize`initialize`package above when the service is initialized.

1.`static/jquery.js`is the js file introduced into the project's default home page.

1.`views/welcome.html`is the view file used by the default home page of the project. You can see its figure in the
`Demo Controller`mentioned above.

# CONTROLLER

## WRITE A CONTROLLER

To write your own GMC Controller, you need to import the package`github.com/snail007/gmc`and define your own`struct`inheriting`gmc.Controller`.

The following example code implements a simple controller that defines a`Hello`method.

```go
package controller

import(
	"github.com/snail007/gmc"
)

type Demo struct {
	gmc.Controller
}

func(this *Demo) Hello() {
	this.Write("fmt.Println(\"Hello GMC!\")")
}
```

## RULES IOF CONTROLLER

1. Unlimited controller name.

1. Suffixes are two or one underscore `__` , `_` controller method, in the routing binding controller will be ignored.

1. The controller method names cannot contain the following names. They are GMC methods to complete the framework function, 
and these names cannot be used in the controllers of students.`MethodCallPre()`,`MethodCallPost()`,`Stop()`,`Die()`,`Tr()`,`SessionStart()`,`SessionDestroy()`,`Write()`,`StopE()`.

1. The method with the name`Before()`is the constructor of the controller, need not be undefined, and is called before 
the execution of the accessed controller method, You can call`this.stop()`to prevent calls to the controller method being accessed, 
but not to block calls to the`After()`controller destructor method. The`this.die()`controller method and`After()`calls can be prevent by`this.die()`.

1. A method whose name is`After()`is the controller's destructor, which is not required to be undefined and is called after the execution of the accessor controller method.

1. The controller members cannot contain the following names, which are used by GMC to complete the framework function. 
The names in the controller cannot be these:
`Response`,`Request`,`Param`,`Session`,`Tpl`,`SessionStore`,`Router`,`Config`,
`Cookie`,`Ctx`,`View`,`Lang`,`Logger`, these members are very useful and we often use them.

## INPUT

The input can be obtained through the `this.request` object, which is the native standard `*http.Request` object. By accessing `GET, POST, COOKIE, upload file` and other data.
In addition, `this.ctx` can be used for some convenient input operations.

## OUTPUT

The output can be passed through the `this.response` object, which is the native standard `http.responsewriter` object. In addition, you can also output most of the data content through `this.write()`,
It can automatically recognize various data types, automatically converted after output to the browser.

## SESSION

GMC's Session is very different from other implementations of this kind, where you either turn it on globally or you turn it off globally, and the downside of this is that no matter which method of the controller you`re accessing,
The framework will perform all kinds of operations to initialize session, causing unnecessary and large performance consumption. GMC avoids this shortcoming by referring to PHP's SESSION implementation mechanism. Where SESSION data needs to be manipulated,
Session data can only be accessed after the session is opened by manually calling this.sessionStart(), which minimizes performance consumption. In addition, if you need to destroy the session data,
Call this.sessiondestroy().

The control can be passed through:
- `this.sessionStart()` opens the session.
- `this.sessiondestroy()` destroys session data.
- `this.session` to access or set Session data, it must be `this.sessionStart()` before it can be used. Otherwise, the object is nil.

## COOKIE

Cookies can be operated not only through the standard `this.Request` and `http.ResponseWriter`, but also through the `this.Cookie` provided by GMC for `Set, Get, Remove` operations.

# HTTP ROUTER

GMC routing supports the following functions:

- Bind the controller to automatically identify the method name, and use lowercase method names in the URL.
- Bind the controller method.
- Bind `http.Handler`.
- Customize the path of the controller in the URL.
- Routing packet.
- Each controller can customize its own URL suffix.
- Each routing group can customize its own URL suffix.

## ACCESS ROUTER

Through the sample project, we call `Router()` method of the newly created Web server or API server object was used to obtain the corresponding routing object, and then it was ready to perform various configurations of routing.

In the Router package for the sample project, simple routing bindings are made through the routing object.

## ROUTER BINDING

Routing binding can be done in the following ways.

### BINDING CONTROLLER

```go
//...
r.Controller("/demo", new(controller.Demo))
//...
```

- The first parameter `/demo` is the path of the controller in the URL.
- The full path to access controller methods is: controller path + all lowercase method names.
- URL access `Demo` controller `Hello` method URL path is: `/demo/hello`.
- A third parameter can also be passed as the URL suffix for the method that this controller can access in the URL.

### BINDING METHOD

```go
//...
r.ControllerMethod("/",new(controller.Demo),"Index__")
//...
```

`ControllerMethod`, can bind to any public method of the controller, even if the suffix is `__`.

### BINDING `http.Handler`

The following example USES various methods to bind a `http.Handler` `myHanlder` to the URL path `/hello`,
` HandlerAny ` is equivalent to a one-time binding ` GET, POST, DELETE, OPTIONS, PUT, PATCH, HEAD `.

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

### BINDING `gmc.Handle`

The following example USES the ways of binding a GMC Handle `Hello` to the URL path `/hello`, ` HandleAny ` is equivalent to a one-time binding ` GET, POST, DELETE, OPTIONS, PUT, PATCH, HEAD `.

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

## ROUTER PARAMETER

Example：
         
```go
 //...
r.GET("/user/:name",Hello)
 //...
```

`:name` is a placeholder parameter, the value of this location:

- In the controller can be obtained by `this.Param.ByName("name")`.
- Can be obtained by `gmc.P.ByName("name")` in the binding Handle.
- In `http.handler` you can pass:
  ```go
    params := gmc.ParamsFromContext(r.Context())
    params.ByName("name")
  ```

## PARAMETER RULES

Using parameters:

```text
URL binding: /user/:user

 /user/gordon              match
 /user/you                 match
 /user/gordon/profile      not match
 /user/                    not match
```

match all:

```text
URL binding: /src/*filepath

 /src/                     match
 /src/somefile.go          match
 /src/subdir/somefile.go   match
```


# TEMPLATE

The GMC template engine, a wrapper around `text/template`, enhances its functionality. The basic usage of templates is the same syntax as `text/template `.

## CONFIGURATION

The template configuration is in app.toml, and the default is as follows:

```toml
[template]
dir="views"
ext=".html"
delimiterleft="{{"
delimiterright="}}"
```

- dir: the directory where the view template files are stored.
- ext: the suffix of a template file. Only files with this suffix will be parsed.
- delimiterleft: left symbol of the syntax block in the template.
- delimiterRight: right symbol of the syntax block in the template.

## INCLUDE

The template directory is views and the file structure is as follows:

```text
views/
├── common
│   └── head.html
└── user
    └── list.html
```

user/list.html contents:

```text
{{template "common/head" . }}
```

In this way, the template file `common/head.html` is included in `user/list.html`. There is no need to write `.html `suffix,`. `is used to put the data of the current template
Pass to `common/head.html`. If you do not need to pass data, you can omit `.`, that is, `{{template "common/head"}}`.

## LAYOUT

Layout template is often used in the development of web sites. When we render many pages, their basic layout is the same, but some parts are different. This time, we need to use the layout template support.
The principle is that when rendering a template, the result of its rendering is placed in the location specified in the layout template.

The View folder is views and its file structure is as follows:

```text
views/
├── layout
│   └── page.html
└── user
    └── profile.html
```

`page.html`content is as follows:

```text
{{.GMC_LAYOUT_CONTENT}}
```

When we render `profile.html` in the controller, the code is as follows:

```go
this.View.Layout("layout/page").Render("profile")
```

The rendering process is to first render `profile.html` using the view data and then render the result to the template variable `GMC_LAYOUT_CONTENT`,
Then use the data to render `layout/page.html` and output the rendering results to the browser.

## INTERNAL FUNCTIONS

`text/template` functions are as follows:

`and`,`call`,`html`,`index`,`slice`,`js`,`len`,`not`,`or`,`print`,`printf`,`println`,`urlquery`

## COMPARE FUNCTIONS

`text/template` comparison functions is as follows:

`eq`,`ne`,`lt`,`le`,`gt`,`ge`

## GMC TEMPLATE FUNCTIONS

GMC defines a number of useful functions for the use of aspect templates.

1. `tr` is used in internationalization, the first parameter is fixed as `.Lang`, and the second parameter is the key defined in the internationalization profile. The third parameter is the prompt message, which will not be output. The data type returned is: `template.HTML`

Example:

`{{tr.lang "key001" "prompt message, what is this"}}`

1. `trs` has the same function as `tr`, but the data type returned is String.

1. `string` has only one argument. Convert the data to a string.

1. `tohtml` has only one parameter, which converts the data to template.html type.

1. `val` gets the template variable content. If the variable does not exist, return empty string instead of `<no value>` Can be used to safely output template variables.
Two arguments, the first fixed is`. `and the second variable name.

Example:

`{{val . "name"}}`

Based on the introduction of the other [sprig](https://github.com/masterminds/sprig) definition of a large number of useful template method,
GMC to streamline, leaner, in all the [template/sprig/docs](https://github.com/snail007/gmc/tree/master/http/template/sprig/docs) , the method in the document template which can be used directly in GMC.

## GMC TEMPLATE VARIABLE

In the template, GMC has the following variables ready for the template that can be used directly in any template. They have the following meanings:

1. `.Lang `is the standard format after the HTTP header accept-language is parsed by the client, such as: zh-CN. it is used by function `tr` .
1. `.G `GET data, is a` map[string][string] `, `{{.G.key}}` key is the key name in GET data.
1. `.P `POST data, is a` map[string][string] `, `{{.P.key}}` key is the key name in the POST data.
1. `.S `SESION data, is a` map[string][string] `, `{{.S.key}}` key is the name of the key in the SESSION data.
Tip: The session is only enabled by calling `SessionStart()` in the controller.
1. `.C `POST data, is a` map[string][string] `, `{{.C.key}}` key is the key name in COOKIE data.
1. `.U `is a` map[string] String `, details of `.U` are as follows:
   
  ```go
    //u0 is a url object
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
   
   Tips:
   
   Full format of URL:`[scheme:][//[userinfo@]host][/]path[?query][#fragment]`
 
## GET RENDER CONTENTS

By default, the `this.View,Render()` Render template is used inside the controller and the result is output to the browser.
If we want to get the render results instead of the output to the browser, we can use the `this.View.RenderR()` render template, which will return the render results.

# Web SERVER

GMC classifies Web development of Web sites with Web pages, such as news stations, and management background as Web services. For such projects, GMC specially designs` Web services`.

By default, GMC Web services contain modules such as templates, static files, databases, Sessions, caching, and internationalization, which can be turned on and off in the configuration, out of the box.

For the Web project generated through the GMCT tool chain, the default configuration file is located at: `conf/app.toml`. App.toml is the core of the project, and almost all GMC functions are configured here.

The corresponding Web server in GMC is: `gmc.HTTPWebServer`. We can also see that in the main file of the project, a Web service is launched through the object of `gmc.APP`.

You can also perform some of your own initialization and other operations before and after the service is initialized.

The main file of the Web project generated by GMCT is as follows:

```go
package main

import(
	"github.com/snail007/gmc"
	"mygmcweb/initialize"
)

func main() {

	// 1. create an default app to run.
	app := gmc.New.AppDefault()

	// 2. add a http server service to app.
	app.AddService(gmc.ServiceItem{
		Service: gmc.New.HTTPServer(),
		AfterInit: func(s *gmc.ServiceItem)(err error) {
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

It mainly does three things:

1. Create a default GMC. APP object to manage the entire life cycle of programs and services. The default APP object will use `conf/app.toml` as the configuration file.
1. Create a GMC.httpWebServer service object, add it to the service list of APP, and define some of its own initialization after the service initialization.
1. Start the APP, catch the error message, and then output the error. After the APP started, `Run()` will block until the APP is closed manually or an exception occurs.

# API SERVER

We often write some HTTP API interfaces for APP or a third party, which is to provide external data interface without pages of this API service. GMD has specially designed `API service`.

API services are more streamlined and perform better than Web services, removing modules such as templates, static files, internationalization, and sessions.

## API PROJECT

For the API project generated by GMCT tool chain, the default configuration file is located in: `conf/app.toml`. App.toml is the core of the project, and almost all GMC functions are configured here.
This project is suitable for the most self-contained API project.

The main file of the project is as follows:

```go
package main

import(
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

Mainly did the following things:

1. Create a default GMC. APP object to manage the entire life cycle of programs and services.
1. Create a configuration object and set `conf/app.toml` as the configuration file.
1. Create a GMC.APIServer service object and use the configuration object above.
1. Some functions implemented after initialization, such as database, cache and routing.
1. Add API objects to the list of APP management services.
1. Start the APP, catch the error message, and then output the error. After the APP started, `Run()` will block until the APP is closed manually or an exception occurs.

## SIMPLE API

If all you want to do is provide one or a few simple interfaces and you don't need any extra functionality, you don't even need any configuration files,
Just listen on the port to manage the routing, you just need to register the routing handler. The `Simple API` generated by GMC is a good example.

Simple API Example：

```go
package main

import(
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
		BeforeInit: func(s gmc.Service, cfg *gmc.Config)(err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	if e := gmc.StackE(app.Run());e!=""{
		app.Logger().Panic(e)
	}
}
```

The example above USES APP to manage API services, and we don't even use APP.

The code is as follows:

```go
package main

import(
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

Mainly did the following things:

1. Create a GMC.APIServer service object, listen to port 7082, and set the URL suffix as `. json`.
1. Registered a URL path `/hello.json`, output a JSON data.
1. Start the API, catch the error message, and then output the error. After the APP is started, Run() will block until the APP is closed manually or an exception occurs.

# DATABASE

GMC provides convenient access support for MYSQL and SQLite3 databases by default. For other types of databases, you are should to use third-party database libraries.

GMC database operation support:
1. Support for multiple data sources.
1. SQLite3 supports encryption.
1. Convenient chain query and update.
1. Map result set to structure.
1. Transaction support.

## MySQL

MySQL database supports connection pooling, various timeout configurations, and detailed configuration information in the `[database]` section of app.toml.

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

It can be found that `database.mysql` is an array, students can configure one, as long as the `id` should be unique.
`ID` is the id used in the code to configure the database object.

### ACCESS

In an API or Web project, after the app is launched, that is, after the database has been initialized successfully, it can be passed in the code
`gmc.DB.DB(id ... String) ` or ` GMC.DB.MySQL(id... String) `gets the database action object.

Example code:

```go

package main

import(
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
	// gmc.DB.DB(id) to load`id`named default driver.
	db := gmc.DB.DB().(*gmc.MySQL)
	//do something with db
	db.AR()
}
```

Due to a database that is more, detailed please reference [here](https://github.com/snail007/gmc/blob/master/db/mysql/README.md)

## SQLITE3

SQLITE3 database various configuration and encryption, detailed configuration information in the `[database]` section of app.toml.

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

### ACCESS

In an API or Web project, after the app is launched, that is, after the database has been initialized successfully, it can be passed in the code
`gmc.DB.DB(id ... String) ` or ` GMC.DB.SQLite3(id... String) `gets the database action object.

Due to a database that is more, detailed please reference [here](https://github.com/snail007/gmc/blob/master/db/sqlite3/README.md)

# CACHE

## INTRODUCTION

GMC Cache supports Redis, File and memory Cache. In order to adapt to different business scenarios, developers can also implement the `gmccore.Cache` interface by themselves.
Then register your cache with `gmccacheheller.AddCacheU(id,cache)`. You can then get your registered cache object through `gmc.Cache.Cache(id)`.

If you want to use the cache in the `GMC` project, you need to modify the `[cache]` section of the configuration file `app.toml`. First, set the default cache type, such as redis `default=redis`.
You then need to modify the `enable=true` part of the corresponding cache driver `[[cache.redis]]`. Multiple caches can be configured for each driver type, each id must be unique, and an `id` of `default` will be used as the default.

Such as:

If redis is configured with more than one, `gmc.Cache.Redis()` retrieves the one with the id default.

## CONFIGURATION

```shell
[cache]
default="redis" //Set default cache validation configuration items, such as redIS cache validation by default
[[cache.redis]] //redis configuration
[[cache.file]]  //file cache configuration
[[cache.memory]]//Memory cache configuration item, where cleanupInterval is second for automatic garbage collection time
```

Configure the cache of API or Web service, which is stared by `gmc.APP`. When you enable caching in the profile app.toml,
The Cache can then be used via `gmc.Cache.Cache()`.

### REDIS

Based on redigo@v2.0.0 implementation, support redis official mainstream method call, can be applied to most business scenarios.

```shell
[[cache.redis]]
debug =true    // Whether debugging is enabled
enable =true   // Enable redis cache
id ="default"  // Cache unique name
address =":6379" // Redis client link address
prefix=""
password=""
timeout =10    // Wait for the maximum length of connection pool allocation(milliseconds), beyond which connections that are not available occur
dbnum =0       // When connecting to Redis, the DB number is 0 by default.
maxidle =10    // Maxidle connections in the connection pool are allowed at most
maxactive =30  // The maximum number of connections supported by the connection pool
idletimeout =300       // the maximum length(milliseconds) for a connection to the IDLE state, and the timeout is released
maxconnlifetime =3600  // the lifetime of a connection(milliseconds), timed out and unused is released
wait=true // when no idle connection to using, to wait or not.
```

### MEMORY

GMC memory cache is a lightweight GO cache implementation that does not use the `hash.Hash` function of go's `hash/ FNV`. Instead, it USES the `djb3` algorithm, which improves the efficiency of bulk file storage by about one time compared to the standard cache
The configuration information is shown below. `cleanupInterval` is the time for GC, indicating that an expiration cleanup will be conducted every 30s,id is the default connection pool ID,enable indicates whether caching is enabled, and default is off.

```shell
[[cache.memory]]
enable=false
id="default"
cleanupinterval=30
```

### FILE

The configuration information is shown below, dir represents the file directory of the cache, `cleanupInterval` information represents the time of GC, indicates that the expiration will be cleaned every 30s, ID is the default connection pool ID,enable indicates whether the cache is turned on and turned off by default.

```shell
[[cache.file]]
enable=false
id="default"
dir="{tmp}"
cleanupinterval=30
```

## PACKAGE CACHE

Of course, `gmccache` package can also be used separately. It does not depend on the GMC framework. It instantiates the configuration object and initializes the cache configuration by itself.

```go
package main

import(
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
	// gmc.Cache(id) to load`id`named default driver.
	c := gmc.Cache.Cache()
	c.Set("test", "aaa", time.Second)
	c.Get("test")
}
```

# I18n

The content of the GMC internationalization file is composed of multiple lines of key=value. File names use the format of the standard HTTP header `Accept-Language`, such as: zh-CN, en-US.
The suffix is `.toml`.

Value is the first parameter in FMT.Printf, where the supported placeholder can be written.
Template output, you can format the output.

Such as:

{{printf(trs .Lang "foo2") 100 }}

Printf, tr, and String are template functions.
1. Printf is a formatted string, just like FMT.Printf.
1. Tr is the translation function. The first parameter is fixed, which corresponds to this.lang in the controller.
1. Since tr returns template.html data, you can use the String function to convert it to string.

For example, the content of Chinese translation document `zh-CN.toml` :

```text
hello="你好"
foo="任意内容"
foo1="你的年龄是%d岁"
foo2="还支持html哟，<b>加粗</b>"
```

For example, English translation file `en-US.toml` contents:

```text
hello="Hello"
foo="Bar"
```

Hello, foo is really just a key to find the corresponding value, the corresponding content.

These are the rules and principles of GMC internationalization.

Internationalization function. Before use, the `[i18n]` module in app.toml should be configured to enable internationalization function, and set `enable=true`.
Set internationalization file directory `dir="i18n"`, default is the program level directory `i18n`, where the above mentioned internationalization translation files are stored.

The controller's `this.Lang` member variable, which is used by the visitor's browser before the controller method is called when a request is made
`Accept-language` is initialized well. In the controller, you can find multiple internationals by `this.Tr(key, hint string)`
Change the language that matches the best in the language file, and then search for the translation result of matching key in the language file.

Internationalization-related configurations are as follows:

```toml
[i18n]
enable=false
dir="i18n"
default="zh-cn"
```

`default` is the default language, used if the language of the user's HTTP request header and the internationalization module cannot find a matching language` default `is set to the language for translation.

# MIDDLEWARE

Both THE Web and API servers of GMC support middleware. When existing functions cannot meet your requirements, you can complete various functions by registering middleware, such as: `authorization`, `logging`, `modification request` and so on.
From the time the GMC receives a user request to the time the controller method or handler is actually executed, the middleware is divided into four types according to the order and timing of execution.

With the names` Middleware0 `, `Middleware1`, `Middleware2`, `Middleware3`, each type of middleware can be registered with more than one.
They are executed in the order in which they were added, but when one of the middleware returns FALSE, the middleware behind that middleware is skipped.

A middleware is a function.

In the API service it is defined as` func(ctx *gmcrouter.Ctx, Server *APIServer)(isStop bool) `.

In the Web service it is defined as` func(ctx *gmcrouter.Ctx, Server *HTTPServer)(isStop bool) `.

The diagram of the API and Web HTTP server workflow architecture is shown below. The sequence and timing of their execution can help you quickly master the use of middleware.

The following figure is for Web and API services, think of the controller below as a handler registered in the API, which is the flow chart of API service middleware.

The STOP in the figure corresponds to when the value returned by the middleware function is true.

<img src="https://github.com/snail007/gmc/blob/master/doc/images/http-and-api-server-architecture.png?raw=true" width="960" height="auto"/>  

# HOT UPDATE & RESTART

This feature is only available for Linux platforms, not Windows systems.

When our application is deployed online, we face a smooth restart/smooth upgrade problem, which is to restart and upgrade the application without interrupting the service that is currently connected to ensure that the service is always available.

Both Web and API services launched through GMC.app support smooth restart/smooth upgrade, which is very simple to use. When you need to restart, use pkill or kill command to send `USR2` signal to the program.

Such as:

```shell
pkill -USR2 website
kill -USR2 11297
```

In this example, `website` is the program name, and `11297` is the program's `PID`. Either way, you can choose your habits.

# USEFUL PACKAGE

## gpool

## gmcmap

## sizeutil

## timeutil

## gmclog

## cast

# GMCT TOOL CHAIN

## ABOUT

GMC framework, in order to reduce user learning costs and accelerate development, provides an open source GMCT tool chain, which can currently accomplish the following functions:

1. One-command to initialize a new GMC based Web project, including the complete project code of controller, routing, cache, database, view, etc.

1. One-command initialization of a new GMC-based API project, including the complete project code of controller, routing, cache, etc.
API projects are lighter than Web projects, and GMC defines a separate API server for API type projects.

1. One-command initialization of a new GMC API lightweight project, containing the basic code to build an API service.
API lightweight projects, suitable for embedding into any existing project code, easily complete an API service with several calls.

1. One-command pack view files to go file, GMC integrates the function of trying to package into binary, before compiling the project code,
Simply execute `gmct tpl --dir ../views` command, you can package all views in the view directory
The diagram file is a GO file. The project is then compiled normally, and the view is packaged into the compiled binaries of the project.

1. One command pack the website static files to go files, GMC integrated the website static files, such as: CSS, JS, font and so on, 
Command `gmct static --dir ../static` to package all files in the static resource directory as a go file. The project is then compiled normally, and the static resource files are packaged into the project binary file.

1. Hot compile the project. During the development process of the project, we will constantly modify the GO file or view file, and then need to manually recompile and run it.
In order to see the modified effect, which took a lot of time, in order to solve this problem, just need to execute `gmct run` in the project directory
The GMCT tool detects changes you make to the project, automatically re-compiles and runs the project, you change the project file, and just refresh the browser, you can see the latest changes.

## INSTALL

To use the tool chain, you need to install Git locally, configure the Golang environment, Golang version 1.12 or above, and set the `GOPATH` environment variable.
The `PATH` environment variable contains the `$GOPATH/bin` directory. If you are not familiar with the configuration of environment variables, you can search and learn to configure system environment variables.

There are two ways to install the GMCT tool chain. One is to compile directly from source code. One is to download the compiled binary.

### 1. COMPLIES FROM SOURCE

This method needs to ensure that your network can properly `go get` and `go mod` to various dependent packages. For special reasons, such as unable to download the dependent packages,
You can use the proxy to download dependency packages by setting the `GOPROXY` environment variable via the proxy. Please refer to [setting the GOPROXY environment variable](https://goproxy.io/).

Then open a command line and execute the following commands in turn. For the first installation, more dependent packages need to be downloaded. Please make sure that the network is normal and wait patiently.

Linux:

```shell
export GO111MODULE=on 
git clone https://github.com/snail007/gmct.git
cd gmct && go mod tidy
go install
gmct --help
```

Windows:

```shell
set GO111MODULE=on
git clone https://github.com/snail007/gmct.git
cd gmct && go mod tidy
go install
gmct --help
```

### 2. DOWNLOAD BINARY

Download address: [GMCT tool chain](https://github.com/snail007/gmct/releases), need according to your operating system platform, download the corresponding binary package, and then extract the `gmct` binary or `gmct.exe`
Just put it in the `$GOPATH/bin` directory, then open a command line and execute `gmct -- Help`. If you have help information for `GMCT`, the installation is successful.

## GENERATE WEB PROJECT

GMCT initializes projects by default using `go mod` to manage dependencies. The project path starts with: `$GOPATH/src`.
Only one parameter `--pkg` is the project path is needed to initialize the project.

Operate the steps and execute the following commands sequentially:

```shell
gmct new web --pkg foo.com/foo/myweb
cd $GOPATH/src/foo.com/foo/myweb
gmct run
```

Open the browser and visit http://127.0.0.1:7080 to see the new Web project in action.

## GENERATE API PROJECT

GMCT initializes projects by default using `go mod` to manage dependencies. The project path starts with: `$GOPATH/src`.
Only one parameter `--pkg` is the project path is needed to initialize the project.

Operate the steps and execute the following commands sequentially:

```shell
gmct new api --pkg foo.com/foo/myapi
cd $GOPATH/src/foo.com/foo/myapi
gmct run
```

Open the browser to `http://127.0.0.1:7081`, and you can see the new API project in action.

## GENERATE SIMPLE API PROJECT

GMCT initializes projects by default using `go mod` to manage dependencies. The project path starts with: `$GOPATH/src`.
Only one parameter `--pkg` is the project path is needed to initialize the project.

Operate the steps and execute the following commands sequentially:

```shell
gmct new api-simple --pkg foo.com/foo/myapi0
cd $GOPATH/src/foo.com/foo/myapi0
gmct run
```

Open the browser at `http://127.0.0.1:7082`, and you can see the new API lightweight project in action.

## GENERATE ADMIN PROJECT

GMCT initializes projects by default using `go mod` to manage dependencies. The project path starts with: `$GOPATH/src`.
Only one parameter `--pkg` is the project path is needed to initialize the project.

Operate the steps and execute the following commands sequentially:

```shell
gmct new admin --pkg foo.com/foo/myadmin
cd $GOPATH/src/foo.com/foo/myadmin
```

1. Admin Console using MySQL database, you need to create a database, and then import data file, located in: `docs/db.sql`.
1. Modify the database configuration in file `conf/app.tml`.
1. Execute `gmct run` in the project directory.

Open the browser visit `http://127.0.0.1:7082`, and you can see the new ADMIN project in action.

username：`root`

password：`123456`

## GENERATE CONTROLLER

GMCT can generate a controller with several empty methods, eliminating the need to manually set up the controller.

The following commands need to be executed under the controller directory.

To operate, execute the following command:

```shell
gmct controller -n Member
```

After executing, we found one more controller file in the current directory: 'member.go'.

Parameters description:

`-n` is the name of the controller 'struct', and the file name USES all lowercase controller names.

`-f` if the file to be generated already exists and is not overwritten by default, the '-f' parameter can be used to enforce overwriting.

## GENERATE TABLE MODEL

GMCT can generate a table model file with common methods to CURD table data without the hassle of writing a table model.

The following commands need to be executed under the model directory.

To operate, execute the following command:

```shell
gmct model -n table
```

After executing, we found a model file in the current directory: `user.go`.

Parameter description:

-n is the name of the database table. If the table prefix is set in the configuration file, there is no need to write the table prefix, such as the table name: `user`, `system_config`.

-t generates a model of MySQL by default, and also supports SQLite3, with values of: `mysql` or `sqlite3`, default is `mysql`.

-f if the file to be generated already exists and is not overwritten by default, the `-f` parameter can be used to enforce overwriting.

## PACK TEMPLATE FILES INTO BINARY

The GMC View module supports packaging view files into compiled binaries. Since the packaging functionality is related to the project directory structure, it is assumed that the directory structure is the Web project directory structure generated by GMCT.

The directory structure is as follows:

```text
new_web/
├── conf
├── controller
├── initialize
├── router
├── static
└── views
```

The following steps are required to complete this function:

```shell
cd initialize
gmct tpl --dir ../views
```

When you execute the command, you will notice that there is an extra go file in the initialize directory that is prefixed by `gmc_templates_bindata_`,
For example: `gmc_templates_bindata_2630881503983182670.go`. There's an init method in this file, this is done automatically when the `initialize` package is referenced and view file binary data is injected into the GMC view module.

When you go Build your project, avoid later development, run the code and always use the view data in the go file, the initialize execute on the directory is: `GMCT TPL -- Clean` to safely clean the go files generated above.

## PACK STATIC FILES INTO BINARY

GMC's HTTP static file module supports packaging static files into compiled binaries. Since the packaging function is related to the project directory structure. So it is assumed that the directory structure is the web project directory structure generated by GMCT.

The directory structure is as follows:

```text
new_web/
├── conf
├── controller
├── initialize
├── router
├── static
└── views
```

The following steps are required to complete this function:

```shell
cd initialize
gmct static --dir ../static
```

After executing the command, you will notice that there is an extra go file in the initialize directory that is prefixed by `gmc_static_bindata_`,
For example: `gmc_static_bindata_1780615241186372497. Go`. There's an init method in this file,
This is done automatically when the `initialize` package is referenced, injecting view file binary data into GMC's HTTP static file service module.

When you go Build your project, avoid later development, run the code using the static file data in the go file,
The initialize execute on the directory is: `gmct static --clean` to safely clean the go files generated above.

In addition, after static files are packed into binary, the order to find static files is:
1. Find out if the file is in binary data.
2. Find the static directory static if there is this file.  

## PACK I18N FILES INTO BINARY

GMC's i18n module supports packaging i18n files into compiled binaries. Since the packaging function is related to the project directory structure. So it is assumed that the directory structure is the web project directory structure generated by GMCT.

The directory structure is as follows:

```text
new_web/
├── conf
├── controller
├── initialize
├── i18n
├── router
├── static
└── views
```

The following steps are required to complete this function:

```shell
cd initialize
gmct i18n --dir ../i18n
```

After executing the command, you will notice that there is an extra go file in the initialize directory that is prefixed by `gmc_i18n_bindata_`,
For example: `gmc_i18n_bindata_1780615241186372497. Go`. There's an init method in this file,
This is done automatically when the `initialize` package is referenced, injecting view file binary data into GMC's i18n module.

When you go Build your project, avoid later development, run the code using the i18n files data in the go file,
The initialize execute on the directory is: `gmct i18n --clean` to safely clean the go files generated above.

## HOT BUILD & AUTO-RE-RUN

During the development of the project, we will constantly modify the go file or view file, and then need to manually recompile and run it.
In order to see the modified effect, which took a lot of time, in order to solve this problem, just need to execute `gmct run` in the project compile directory, 
the GMCT tool detects changes you make to the project, automatically recompiles and runs the project, you change the project file, and just refresh the browser. You can see the latest changes.

Executing `gmct run` will generate a configuration file named `gmcrun.toml` in the current directory. You can modify this file to customize `gmct run` `compile behavior.

The default configuration is as follows:

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

Configuration description:

1. `monitor_dirs` monitor directory, the default is the current directory, you can specify multiple, array form.
1. `args`         extra parameter passed to` go build `, array, multiple parameters written separately.
1. `env`          sets the environment variable at the time of` go build `execution. Multiple, array forms can be specified.
1. `include_exts` monitor file suffix, only monitor changes to this suffix file. You can specify multiple arrays.
1. `include_files` sets additional monitored files, supports relative path and absolute path, can specify multiple, array form.
1. `exclude_files` sets additional unmonitored files, supports relative path and absolute path, can specify multiple, array form.
1. `exclude_dirs`  sets additional unmonitored directories, supports relative path and absolute path, can specify multiple, array form.

`${DIR}` is available in `monitor_dirs`,` include_files`, `exclude_files`,` exclude_dirs`.

The variable `${DIR}` represents the absolute path of the current directory, with no `/` at the end.
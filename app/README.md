## GMC APP

## Features

A GMC APP boot your net/any service with hot reload , and auto manange many services. 

## HOT RELOAD APP CODE

If your service use none gmc.Config, just call `SetNoneMainConfigFile(true)`.

```golang
// api is a object implements gmc.Service interface.
api:=fooService

// create a app, and set no real config file to parse.
app := gmc.NewAPP().SetNoneMainConfigFile(true)

// parse init something
err := app.ParseConfig()
if err != nil {
    panic(err)
}
// add your service to app
app.AddService(gmc.ServiceItem{
    Service: api,
})

app.Logger().Panic(app.Run())
    
```

## Service

```golang
type Service interface {
	Init(cfg *gmcconfig.GMCConfig) error
	Start() error
	Stop()
	GracefulStop()
	SetLog(*log.Logger)
	InjectListeners(net.Listener)
	Listener() net.Listener
}
```
### INTRO

1.When hot reload.,call stack: `Init()->InjectListeners()->Start()`, so you should using InjectListeners's net.Listener in Start().  

2.When hot reload requested, `Listener()` will be called, to obtain the net.Listener FD pass to sub process.  

`GracefulStop()` will be called to stop your service.  

## HOT RELOAD COMMAND

Only worked on linux.  

command:  

`pkill -USR2 yourappname`

the `-USR2` signal will trigger the gmc app to hot relaod.
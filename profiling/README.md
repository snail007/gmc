## profiling

profiling package help you do profiling your program in easy way.

## demo

```golang
package main

func main() {
	//StartArg will search command line arguments "profiling" , using it's value as store data folder.
	StartArg("profiling")
    defer Stop()
    //do something, keep long time, such 15 minutes
}
```
start your programe with `./foobar -profiling=debug`


```golang
package main

func main() {
    //Start will using the argument as store data folder
	Start("debug")
    defer Stop()
    //do something, keep long time, such 15 minutes
}
```

start your programe with `./foobar`

cpu.prof, memory.prof, block.prof, goroutine.prof, threadcreate.prof files will be create in debug folder after your program exited.

then you can using `go tool pprof foobar memory.prof` to profiling your program.


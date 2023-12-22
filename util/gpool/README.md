## GPool
gpool is a goroutine pool, it is concurrency safed, it can using a few goroutine to run a huge tasks.  
a worker is a goroutine, a task is a function, gpool using a few workers to run a huage tasks.  

- dynamic change workers count support
- worker idle support
- lazy start worker support
- pre alloc worker support

## Demo

```golang
package main

import (
    "fmt"
    "github.com/snail007/gmc/sync/gpool"
)

func main() {
	//we create a poll named "p" with 3 workers
	p := gpool.New(3)
	//after New, you can submit a function as a task, you can repeat Submit() many times anywhere as you need.
	a := make(chan bool)
	p.Submit(func() {
		a <- true
	})
	fmt.Println(<-a)
}

```

## Testing And Code coverage

```text
ok      github.com/snail007/gmc/util/gpool      9.341s  coverage: 95.2%
total:                                                  (statements)            95.2%
```
## Benchmark

```text
go test -bench=. -run=none
goos: darwin
goarch: amd64
pkg: github.com/snail007/gmc/util/gpool
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkSubmit-16       5806988               193.8 ns/op
BenchmarkWorker-16       5232226               232.5 ns/op
PASS
ok      github.com/snail007/gmc/util/gpool      3.071s
```

```text
go test -bench=. -benchtime=3s -run=none
goos: darwin
goarch: amd64
pkg: github.com/snail007/gmc/util/gpool
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkSubmit-16      18357210               190.3 ns/op
BenchmarkWorker-16      15653331               232.4 ns/op
PASS
ok      github.com/snail007/gmc/util/gpool      7.921s
```
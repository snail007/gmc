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
BenchmarkSubmit/pool_size:20-16                   717519              3822 ns/op
BenchmarkSubmit/pool_size:40-16                   932514              3944 ns/op
BenchmarkSubmit/pool_size:60-16                   789867              4295 ns/op
BenchmarkSubmit/pool_size:80-16                  1000000              5250 ns/op
BenchmarkSubmit/pool_size:100-16                  972837              5719 ns/op
BenchmarkSubmit/pool_size:200-16                  798679              6224 ns/op
BenchmarkSubmit/pool_size:400-16                  683112              6566 ns/op
BenchmarkSubmit/pool_size:600-16                  571062              5244 ns/op
BenchmarkSubmit/pool_size:800-16                  664258              9264 ns/op
BenchmarkSubmit/pool_size:1000-16                 495985              5359 ns/op
BenchmarkSubmit/pool_size:10000-16                564003              6340 ns/op
BenchmarkSubmit/pool_size:20000-16                563130              6611 ns/op
BenchmarkSubmit/pool_size:30000-16                572671              6293 ns/op
BenchmarkSubmit/pool_size:40000-16                529896              5777 ns/op
BenchmarkSubmit/pool_size:50000-16                495811              5074 ns/op
BenchmarkJob/pool_size:20-16                      546973              4891 ns/op
BenchmarkJob/pool_size:40-16                      525769              4606 ns/op
BenchmarkJob/pool_size:60-16                      514962              5270 ns/op
BenchmarkJob/pool_size:80-16                      522291              5347 ns/op
BenchmarkJob/pool_size:100-16                     537969              4681 ns/op
BenchmarkJob/pool_size:200-16                     609165              5018 ns/op
BenchmarkJob/pool_size:400-16                     513234              5614 ns/op
BenchmarkJob/pool_size:600-16                     591480              5476 ns/op
BenchmarkJob/pool_size:800-16                     537184              5458 ns/op
BenchmarkJob/pool_size:1000-16                    475809              5273 ns/op
BenchmarkJob/pool_size:10000-16                   723447              6300 ns/op
BenchmarkJob/pool_size:20000-16                   591313              4874 ns/op
BenchmarkJob/pool_size:30000-16                   508342              4536 ns/op
BenchmarkJob/pool_size:40000-16                   484904              5399 ns/op
BenchmarkJob/pool_size:50000-16                   458240              5261 ns/op
PASS
ok      github.com/snail007/gmc/util/gpool      101.870s
```

```text
go test -bench=. -benchtime=3s -run=none
goos: darwin
goarch: amd64
pkg: github.com/snail007/gmc/util/gpool
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkSubmit/pool_size:20-16                  1000000              3702 ns/op
BenchmarkSubmit/pool_size:40-16                  1000000              6413 ns/op
BenchmarkSubmit/pool_size:60-16                  1000000              4236 ns/op
BenchmarkSubmit/pool_size:80-16                  1000000              4683 ns/op
BenchmarkSubmit/pool_size:100-16                 1000000              7908 ns/op
BenchmarkSubmit/pool_size:200-16                 1000000              6421 ns/op
BenchmarkSubmit/pool_size:400-16                 1000000              7677 ns/op
BenchmarkSubmit/pool_size:600-16                 1000000             10708 ns/op
BenchmarkSubmit/pool_size:800-16                 1000000              9914 ns/op
BenchmarkSubmit/pool_size:1000-16                1000000              7588 ns/op
BenchmarkSubmit/pool_size:10000-16               1000000              7316 ns/op
BenchmarkSubmit/pool_size:20000-16               1000000              8698 ns/op
BenchmarkSubmit/pool_size:30000-16               1000000              7268 ns/op
BenchmarkSubmit/pool_size:40000-16               1000000              7404 ns/op
BenchmarkSubmit/pool_size:50000-16               1000000              9545 ns/op
BenchmarkJob/pool_size:20-16                     1000000              6091 ns/op
BenchmarkJob/pool_size:40-16                     1000000              6476 ns/op
BenchmarkJob/pool_size:60-16                     1000000              4791 ns/op
BenchmarkJob/pool_size:80-16                     1000000              5697 ns/op
BenchmarkJob/pool_size:100-16                    1000000              5325 ns/op
BenchmarkJob/pool_size:200-16                    1000000              6210 ns/op
BenchmarkJob/pool_size:400-16                    1000000              5936 ns/op
BenchmarkJob/pool_size:600-16                    1000000              6310 ns/op
BenchmarkJob/pool_size:800-16                    1000000              8020 ns/op
BenchmarkJob/pool_size:1000-16                   1000000              7428 ns/op
BenchmarkJob/pool_size:10000-16                  1000000              6842 ns/op
BenchmarkJob/pool_size:20000-16                  1000000              7807 ns/op
BenchmarkJob/pool_size:30000-16                  1000000              5834 ns/op
BenchmarkJob/pool_size:40000-16                  1000000              5572 ns/op
BenchmarkJob/pool_size:50000-16                  1000000              6033 ns/op
PASS
ok      github.com/snail007/gmc/util/gpool      204.891s
```
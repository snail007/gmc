## Testing And Code coverage

```text
=== RUN   TestNew
--- PASS: TestNew (0.00s)
    gpool_test.go:14: New is okay
=== RUN   TestSubmit
--- PASS: TestSubmit (0.00s)
    gpool_test.go:38: Submit is okay
2020/09/10 11:08:04 GPool: a task stopped unexceptedly, err: runtime error: integer divide by zero
=== RUN   TestStop
--- PASS: TestStop (0.00s)
    gpool_test.go:55: Stop is okay
=== RUN   TestSetLogger
--- PASS: TestSetLogger (0.00s)
=== RUN   TestRunning
--- PASS: TestRunning (0.04s)
    gpool_test.go:80: Running is okay
=== RUN   TestAwaiting
--- PASS: TestAwaiting (0.04s)
    gpool_test.go:103: Awaiting is okay
PASS
coverage: 94.0% of statements
ok      github.com/snail007/gmc/sync/gpool      0.090s
```
## Benchmark

```text
go test -bench=. -run=none
goos: darwin
goarch: amd64
pkg: github.com/snail007/gmc/sync/gpool
BenchmarkSubmit-12      10000000               140 ns/op
BenchmarkWorker-12      30000000               104 ns/op
PASS
ok      github.com/snail007/gmc/sync/gpool      13.985s
```

```text
go test -bench=. -benchtime=3s -run=none
goos: darwin
goarch: amd64
pkg: github.com/snail007/gmc/sync/gpool
BenchmarkSubmit-12      20000000               231 ns/op
BenchmarkWorker-12      100000000               33.1 ns/op
PASS
ok      github.com/snail007/gmc/sync/gpool      34.705s
```
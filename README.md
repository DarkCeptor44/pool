# pool

[![made](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)
[![powered](https://forthebadge.com/images/badges/powered-by-black-magic.svg)](https://forthebadge.com)

This project started out as _"a Pure Go equivalent-ish of Python's `multiprocessing.Pool` with zero dependencies"_ but became a typical semaphore implementation using channels due to how complex multiprocessing is.

## Installation

```bash
go get -u github.com/DarkCeptor44/pool/v2
```

## Example

### Run

```go
import (
    "fmt"

    "github.com/DarkCeptor44/pool/v2"
)

func job(index, value int) {
    fmt.Printf("worker%d got %d\n", index, value)
}

func main(){
    arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // makes pool with 5 workers and the slice of ints
    _ := pool.Run(10, arr, job)
}
```

### RunAndReturn

```go
import (
    "fmt"

    "github.com/DarkCeptor44/pool/v2"
)

func multiply(_, value int) int {
    return value * 2
}

func main(){
    arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // makes pool with 5 workers and the slice of ints, results is the slice of ints that were changed by the job, keep in mind the results are in a different order
    results, _ := pool.RunAndReturn(5, arr, multiply)
}
```

## Tests

```bash
$ go test -v
=== RUN   TestNormal
=== RUN   TestNormal/Normal
    pool_test.go:100: 'Normal' Took 222.7155ms
--- PASS: TestNormal (1.11s)
    --- PASS: TestNormal/Normal (1.11s)
=== RUN   TestPool
=== RUN   TestPool/NoValues
=== RUN   TestPool/Pool
    pool_test.go:100: 'Pool' Took 63.05026ms
--- PASS: TestPool (0.32s)
    --- PASS: TestPool/NoValues (0.00s)
    --- PASS: TestPool/Pool (0.32s)
=== RUN   TestPoolWithReturn
=== RUN   TestPoolWithReturn/NoValues
--- PASS: TestPoolWithReturn (0.00s)
    --- PASS: TestPoolWithReturn/NoValues (0.00s)
PASS
ok      github.com/DarkCeptor44/pool/v2 1.448s
```

## Benchmarks

```bash
$ go test -bench .
goos: windows
goarch: amd64
pkg: github.com/DarkCeptor44/pool/v2
cpu: AMD Ryzen 7 3800X 8-Core Processor
BenchmarkNormal-16      13556533                83.97 ns/op           48 B/op          1 allocs/op
BenchmarkPool-16          168236              7033 ns/op             738 B/op         14 allocs/op
PASS
ok      github.com/DarkCeptor44/pool/v2 4.209s
```

## Vulnerabilities

Checked with [govulncheck](https://github.com/golang/vuln):

```bash
$ govulncheck .
No vulnerabilities found.
```

## License

This project is licensed under the MIT License, see [LICENSE](LICENSE) for details.

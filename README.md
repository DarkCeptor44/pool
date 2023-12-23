# pool

Pure Go equivalent-ish of Python's `multiprocessing.Pool`, with zero dependencies.

## Installation

```bash
go get -u github.com/DarkCeptor44/pool
```

## Example

### Action

```go
import (
    "fmt"

    "github.com/DarkCeptor44/pool"
)

func multiply(worker pool.Worker, value int) error {
    n := value * 2
    fmt.Printf("%s got %d\n", worker.Name(), n)
    return nil
}

func main(){
    arr := []int{1, 2, 3, 4, 5}

    // makes pool with 10 workers and the slice of ints
    p := pool.NewPool(10, arr)
    _ := p.Run(multiply)
}
```

### Return

```go
import (
    "fmt"

    "github.com/DarkCeptor44/pool"
)

func multiply(worker pool.Worker, value int) (int, error) {
    return value * 2, nil
}

func main(){
    arr := []int{1, 2, 3, 4, 5}

    // makes pool with 10 workers, the slice of ints as given values and slice of ints as return
    // first generic type would be inferred from arr but since the second needs to be explicit both are
    p := pool.NewReturnPool[int, int](10, arr)
    results, _ := p.Run(multiply)
}
```

## Tests

```bash
$ go test -v
=== RUN   TestAction
    pool_test.go:21: Manipulating file file4.txt by worker5
    pool_test.go:21: Manipulating file file1.txt by worker1
    pool_test.go:21: Manipulating file file2.txt by worker3
    pool_test.go:21: Manipulating file file3.txt by worker4
    pool_test.go:21: Manipulating file file5.txt by worker2
--- PASS: TestAction (0.00s)
=== RUN   TestReturn
    pool_test.go:42: [false false false false false]
--- PASS: TestReturn (0.00s)
PASS
ok      github.com/DarkCeptor44/pool 0.194s
```

## Benchmarks

```bash
$ go test -bench . -benchmem
goos: windows
goarch: amd64
pkg: github.com/DarkCeptor44/pool
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkStd/Read-8          787           1448731 ns/op            6975 B/op
     40 allocs/op
BenchmarkStd/Write-8                 583           2423881 ns/op         6979 B/op
     40 allocs/op
BenchmarkAction/Read-8              1868            650108 ns/op         8264 B/op
     60 allocs/op
BenchmarkAction/Write-8              830           1312837 ns/op         8260 B/op
     60 allocs/op
BenchmarkReturn/Read-8              1836            654401 ns/op         8604 B/op
     66 allocs/op
PASS
ok      github.com/DarkCeptor44/pool 10.906s
```

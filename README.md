# skiplist
implement [skiplist](https://en.wikipedia.org/wiki/Skip_list)

# Example

```bash
go get github.com/ISSuh/skiplist
```

```go

import github.com/ISSuh/skiplist


list := skipList.New(5)
list.Set("key", "value")

item := list.Get("key")
fmt.Printf("key : %s / value : %s", item.key, item.value)
```

# Test Case

```bash
$ go test -bench . -benchtime 5000000x
goos: darwin
goarch: arm64
pkg: github.com/ISSuh/skiplist
BenchmarkSet-8           5000000              1587 ns/op        3150955.24 MB/s      144 B/op          4 allocs/op
BenchmarkGet-8           5000000              1834 ns/op        2725794.73 MB/s        7 B/op          0 allocs/op
PASS
ok      github.com/ISSuh/skiplist       18.181s
```

```bash
$ go test -cover                
PASS
        github.com/ISSuh/skiplist       coverage: 100.0% of statements
ok      github.com/ISSuh/skiplist       0.337s
```
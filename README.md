# Largest N values
This program returns URL's from input file corresponding to N largest integer values

File's entry format is the following: `<url><space><long value>`

To run the program use `make run` along with *abosulte file path* provided in the stdin

By default N is set to `10` and batch size is 10k

## Algorithm
Basically we break the input in batches of fixed size, sort the batch and push (or skip if less) the values with their respective indexies in a binary heap. To print the result we iterate over the file again to get the urls.

This approach allows us to have `n * m / k` heap pushes in worst-case scenario, where m = total amount of rows, k = size of the batch. Also we can run multiple batch processing concurrently with few additions to the code.

Overall the complexity is `O(N)`


### Assumtions
- The `long value` part will be non-negative 64-bit integer, unique accross the file and represented as ASCII digits
- The `url` part can be as large as 2048 bytes
- Input files can be extremely large
- The worst-case scenario is a file containing integer values in non-descending order

### Steps
0. Buffers for batch and count sorting are initialized as well as binary heap
1. File get's consumed starting from zero row
    - Every entry gets it's unique index: `batch number + row number inside the batch`
2. Every K rows (where K = batch size) batch buffer gets filled and then
    - we find the maximum sort key in the batch
    - if the heap len = N and the maximum is smaller than smallest heap element we skip the batch
    - else the batch gets sorted using `count sort` (since we know that sort key will always be non-negative integer) and then N largest amounts with their respective indexies get pushed into the heap using the same logic as stated in previous step
3. After all batches are processed (i.e. the whole file is processed) we get the heap filled with indexies of 10 largest entries
4. We do one more read through the file and get the respective URLs

## Benchmarks
To run the benchmark use `make bench_X`

### bench_1000
```bash
go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/Barugoo/largest_n
BenchmarkLargest10-8       43413             25800 ns/op          495973 B/op          9 allocs/op
PASS
ok      github.com/Barugoo/largest_n    2.143s
```

### bench_100000
```bash
go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/Barugoo/largest_n
BenchmarkLargest10-8       45669             29512 ns/op          495972 B/op          9 allocs/op
PASS
ok      github.com/Barugoo/largest_n    2.576s
```

### bench_1000000
```bash
go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/Barugoo/largest_n
BenchmarkLargest10-8       35420             29878 ns/op          495973 B/op          9 allocs/op
PASS
ok      github.com/Barugoo/largest_n    2.488s
```

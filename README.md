# package anchor

package `anchor` provides a minimal-memory [AnchorHash](https://arxiv.org/abs/1812.09674) consistent-hash implementation for Go.

```go
import "github.com/wdamron/go-anchorhash"
```

## More Info

* [AnchorHash: A Scalable Consistent Hash (Arxiv)](https://arxiv.org/abs/1812.09674)
* [Docs (godoc)](https://godoc.org/github.com/wdamron/go-anchorhash)

## Benchmarks

* Go 1.12.1
* 2017 Macbook Pro; noisy, with a number of applications running
* 2.9 GHz Intel Core i7
* 16 GB 2133 MHz LPDDR3

* **Capacity = 10**
  * `NewCompactAnchor(10, 10)`: `BenchmarkGetBucket_10_10  200000000  5.81 ns/op`
  * `NewCompactAnchor(10, 9)`: `BenchmarkGetBucket_9_10  200000000  6.71 ns/op`
  * `NewCompactAnchor(10, 5)`: `BenchmarkGetBucket_5_10  100000000  10.8 ns/op`
* **Capacity = 1,000,000**
  * `NewAnchor(1000000, 1000000)`: `BenchmarkGetBucket_1m_1m   200000000  7.47 ns/op`
  * `NewAnchor(1000000, 900000)`: `BenchmarkGetBucket_900k_1m  200000000  9.26 ns/op`
  * `NewAnchor(1000000, 500000)`: `BenchmarkGetBucket_500k_1m  100000000  17.6 ns/op`

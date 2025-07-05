# Phase 1

Phase 1 of **kache** included a basic use of `map[string]string` to store key-value pairs, with basic operations like `Set`, `Get`, `Delete`, and `Exists`. The initial implementation was straightforward, focusing on simplicity and ease of use.
I used locks to ensure thread safety using `sync.RWMutex`, allowing multiple readers or a single writer at a time. This was sufficient for basic use cases but did not handle more complex scenarios like TTL (Time To Live) for keys or advanced concurrency control.

Benchmarks
```bash
goos: linux
goarch: amd64
pkg: github.com/vr-varad/kache
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkKacheSet-16                    	 1867969	       627.8 ns/op
BenchmarkKacheGet-16                    	 1405707	       915.9 ns/op
BenchmarkKacheDelete-16                 	 1404774	       881.7 ns/op
BenchmarkKacheExists-16                 	 1264402	       919.3 ns/op
BenchmarkKacheFlush-16                  	 2025438	       734.1 ns/op
BenchmarkKacheConcurrentSet-16          	 1103444	      1038 ns/op
BenchmarkKacheConcurrentGet-16          	 1286689	       870.0 ns/op
BenchmarkKacheConcurrentDelete-16       	  928138	      1363 ns/op
BenchmarkKacheConcurrentExists-16       	 1379953	       772.7 ns/op
BenchmarkKacheConcurrentFlush-16        	 1305098	       875.9 ns/op
PASS
ok  	github.com/vr-varad/kache	27.325s
```

# Phase 2

Phase 2 introduced a more sophisticated data structure using shards to improve concurrency and performance. Each shard is a separate instance of the basic key-value store (map[string]string) and is protected by its own mutex. This allows multiple shards to be accessed concurrently, significantly improving performance in multi-threaded scenarios.
The `ShardsCount` variable was introduced to define the number of shards, allowing for better scalability. The `Set`, `Get`, `Delete`, and `Exists` methods were updated to work with shards, distributing keys across them based on a hash function. This approach reduced contention and improved throughput for concurrent operations.

Benchmarks
```bash
goos: linux
goarch: amd64
pkg: github.com/vr-varad/kache
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkKacheSet-16                    	 1395861	       782.5 ns/op
BenchmarkKacheGet-16                    	  998008	      1356 ns/op
BenchmarkKacheDelete-16                 	  938962	      1336 ns/op
BenchmarkKacheExists-16                 	  936752	      1363 ns/op
BenchmarkKacheFlush-16                  	  587154	      2037 ns/op
BenchmarkKacheConcurrentSet-16          	 4442014	       253.9 ns/op
BenchmarkKacheConcurrentGet-16          	 1389817	       915.1 ns/op
BenchmarkKacheConcurrentDelete-16       	 1225384	      1026 ns/op
BenchmarkKacheConcurrentExists-16       	 1427024	       924.5 ns/op
BenchmarkKacheConcurrentFlush-16        	  414411	      2842 ns/op
ok  	github.com/vr-varad/kache	27.373s
```

The one issue I faced was with `flush` operations. In a multi-threaded environment, flushing all shards at once could lead to contention and performance degradation. Currently what I did is using a for loop to clear each shard individually, which is not optimal but works for now. I went for one more approach of using wait groups to ensure all shards are flushed before returning, but it added complexity without significant performance gains.

As per the benchmarks, the performance of concurrent operations improved significantly, especially for `Set` and `Get` operations. The use of shards allowed for better distribution of keys and reduced contention, leading to faster response times. But the performance of operation in single-threaded enviromnet was more of a hit and miss, depending on the operation type.

The slower performance of operations in single threaded environment was due to the overhead of managing multiple shards, which mainly included getting the index of the shard based on the key's hash. This overhead was negligible in a multi-threaded environment where contention was reduced, but in a single-threaded scenario, it added unnecessary complexity and latency.

Another feature of this phase is **TTL - Time to live** for keys.

At first, I went with the passive approach of checking the TTL during `Get` and `Exists` operations. If the key had expired, it would be deleted from the store. This approach worked but was not efficient, as it required checking the TTL every time a key was accessed. And, if the key was not accessed frequently, it would remain in the store until it was explicitly deleted or the store was flushed.


So for solving the problem of unaccessed keys, I implemented an active cleanup mechanism using Janitor function (which runs in a separate goroutine). The Janitor periodically checks all shards for expired keys and removes them. This approach ensures that expired keys are removed from the store without requiring access to them, keeping the store clean and efficient.
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
BenchmarkKacheSet-16                 	 2672880	       485.2 ns/op
BenchmarkKacheGet-16                 	 2050446	       629.0 ns/op
BenchmarkKacheDelete-16              	 1958770	       597.5 ns/op
BenchmarkKacheExists-16              	 2083227	       631.4 ns/op
BenchmarkKacheFlush-16               	  714211	      1630 ns/op
BenchmarkKacheConcurrentSet-16       	12819930	        89.70 ns/op
BenchmarkKacheConcurrentGet-16       	 2325151	       483.3 ns/op
BenchmarkKacheConcurrentDelete-16    	 2629929	       463.6 ns/op
BenchmarkKacheConcurrentExists-16    	 2600062	       479.4 ns/op
BenchmarkKacheConcurrentFlush-16     	 1000000	      1077 ns/op
PASS
ok  	github.com/vr-varad/kache	33.133s
```

The one issue I faced was with `flush` operations. In a multi-threaded environment, flushing all shards at once could lead to contention and performance degradation. Currently what I did is using a for loop to clear each shard individually, which is not optimal but works for now. I went for one more approach of using wait groups to ensure all shards are flushed before returning, but it added complexity without significant performance gains.

As per the benchmarks, the performance of concurrent operations improved significantly, especially for `Set` and `Get` operations. The use of shards allowed for better distribution of keys and reduced contention, leading to faster response times. But the performance of operation in single-threaded enviromnet was more of a hit and miss, depending on the operation type.

The slower performance of operations in single threaded environment was due to the overhead of managing multiple shards, which mainly included getting the index of the shard based on the key's hash. This overhead was negligible in a multi-threaded environment where contention was reduced, but in a single-threaded scenario, it added unnecessary complexity and latency.

Another feature of this phase is **TTL - Time to live** for keys.

At first, I went with the passive approach of checking the TTL during `Get` and `Exists` operations. If the key had expired, it would be deleted from the store. This approach worked but was not efficient, as it required checking the TTL every time a key was accessed. And, if the key was not accessed frequently, it would remain in the store until it was explicitly deleted or the store was flushed.


So for solving the problem of unaccessed keys, I implemented an active cleanup mechanism using Janitor function (which runs in a separate goroutine). The Janitor periodically checks all shards for expired keys and removes them. This approach ensures that expired keys are removed from the store without requiring access to them, keeping the store clean and efficient.

One thing I focused was to keep a limit on the number of keys in the store. If the number of keys exceeds a certain threshold, the Janitor will also remove the least recently used (LRU) keys to maintain a manageable size. This helps prevent memory bloat and ensures that the store remains performant even with a large number of keys.

The implementation of LRU is simple, a linked list and a index map to keep track of the order of keys. When a key is accessed, it is moved to the front of the list, and when a key is evicted, it is removed from both the list and the index map. This allows for efficient LRU eviction while maintaining the order of keys.

This helps to keep a limit on the number of keys in the store, ensuring that the store remains performant even with a large number of keys. The Janitor runs periodically to check for expired keys and remove them, keeping the store clean and efficient.
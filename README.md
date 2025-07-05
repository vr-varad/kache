# Kache - In Memory DataStore

## Components

### Data Structure


### Key Operations

#### `Set("key", "value", options)`

Sets a value for a given key with optional parameters.

#### `Get("key")`

Retrieves the value for a given key.

#### `Delete("key")`

Deletes a key-value pair from the datastore.

#### `Exixts("key")`

Checks if a key exists in the datastore.

#### `Flush()`

Clears all key-value pairs in the datastore.

### TTL

### Concurrency Control

### Benchmarking

```
goos: linux
goarch: amd64
pkg: github.com/vr-varad/kache
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkKacheSet-16       	 1323517	       890.2 ns/op
BenchmarkKacheGet-16       	  752610	      1401 ns/op
BenchmarkKacheDelete-16    	  966981	      1347 ns/op
BenchmarkKacheExists-16    	  729972	      1396 ns/op
BenchmarkKacheFlush-16     	  550447	      2345 ns/op
PASS
ok  	github.com/vr-varad/kache	7.537s
```

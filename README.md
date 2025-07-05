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

#### `Keys()`

Returns all keys in the datastore.

#### `Size()`

Returns the number of key-value pairs in the datastore.

#### `Flush()`

Clears all key-value pairs in the datastore.

### TTL

### Concurrency Control

### Benchmarking

```
go test -bench .
goos: linux
goarch: amd64
pkg: github.com/vr-varad/kache
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkKacheSet-16       	 2049855	       625.9 ns/op
BenchmarkKacheGet-16       	 1323562	       866.8 ns/op
BenchmarkKacheDelete-16    	 1290350	       896.8 ns/op
BenchmarkKacheExists-16    	 1289115	       903.7 ns/op
BenchmarkKacheSize-16      	 1790197	       623.7 ns/op
BenchmarkKacheFlush-16     	 2026003	       747.3 ns/op
BenchmarkKacheKeys-16      	   10000	    245246 ns/op
PASS
ok  	github.com/vr-varad/kache	15.781s
```

# Kache - In Memory DataStore

## Components

1. Sharded Map Data Structure -> When a key is set, it is hashed to determine which shard it belongs to.
2. Shard -> Each shard is a map that holds key-value pairs.
3. TTL (Time To Live) -> Each key can have an optional TTL, after which it will be automatically deleted.
4. Janitor -> A background process that periodically checks for expired keys and removes them.
5. Max Size -> Each Datastore has a maximum size limit, which can be configured. If the limit is reached, the oldest keys are removed to make space for new ones.
6. LRU -> The datastore uses a Least Recently Used (LRU) eviction policy to manage memory efficiently.

### Data Structure

The datastore is implemented as a sharded map, where each shard is a map that holds key-value pairs. The keys are hashed to determine which shard they belong to. Each key can have an optional TTL, and the datastore supports operations like setting, getting, deleting, and checking for keys.

### TTL

The Time To Live (TTL) feature allows each key to have an expiration time. If a key is set with a TTL, it will be automatically deleted after the specified duration. The janitor process runs periodically to check for expired keys and remove them from the datastore.

### Max Size

The datastore has a maximum size limit (1000 Keys) that can be configured. When the limit is reached, the oldest keys are removed to make space for new ones. This ensures that the datastore does not grow indefinitely and manages memory efficiently.

### LRU Eviction Policy

The Least Recently Used (LRU) eviction policy is used to manage memory efficiently. When the datastore reaches its maximum size, the oldest keys (the least recently used) are removed to make space for new keys. This helps in keeping the most frequently accessed data available while removing less frequently used data.

### Key Operations

#### `Set("key", "value", kache.options)`

Sets a value for a given key with optional parameters.

#### `Get("key")`

Retrieves the value for a given key.

#### `Delete("key")`

Deletes a key-value pair from the datastore.

#### `Exixts("key")`

Checks if a key exists in the datastore.

#### `Flush()`

Clears all key-value pairs in the datastore.

## Installation

```bash
go get github.com/vr-varad/kache
```

---

## Usage

### Basic Example

```go
import (
    "github.com/vr-varad/kache"
    "time"
)

func main() {
    cache := kache.NewKache()

    cache.Set("foo", "bar", kache.Options{TTL: 10}) // TTL in seconds

    value, exists := cache.Get("foo")
    if exists {
        fmt.Println("Value:", value) // Output: bar
    }
}
```

---

## ✅ Test Cases & Use Cases

### 1. **Set and Get**

Set a key and get it back before TTL expires.

```go
kache := NewKache()
kache.Set("key1", "value1", kache.Options{TTL: 10})

value, exists := kache.Get("key1")
// ✅ value should be "value1"
```

---

### 2. **Get Non-Existent Key**

Try to get a key that was never set.

```go
_, exists := kache.Get("nonexistent")
// ❌ should return false
```

---

### 3. **Delete Key**

Delete an existing key and ensure it no longer exists.

```go
kache.Set("key1", "value1", kache.Options{TTL: 10})
kache.Delete("key1")
_, exists := kache.Get("key1")
// ❌ should return false
```

---

### 4. **Exists Check**

Verify whether a key exists in the cache.

```go
kache.Set("key1", "value1", kache.Options{TTL: 10})

kache.Exists("key1")        // ✅ should return true
kache.Exists("nonexistent") // ❌ should return false
```

---

### 5. **Flush the Cache**

Flush the entire cache and verify that no keys exist afterward.

```go
kache.Set("key1", "value1", kache.Options{TTL: 10})
kache.Set("key2", "value2", kache.Options{TTL: 10})
kache.Flush()
// ❌ kache.Exists("key1") or kache.Exists("key2") should be false
```

---

### 6. **Time-to-Live (TTL) Expiry**

Keys should expire after their TTL has passed.

```go
kache.Set("key1", "value1", kache.Options{TTL: 1}) // 1 second TTL

time.Sleep(2 * time.Second)
_, exists := kache.Get("key1")
// ❌ should return false as TTL has expired
```

---

### 7. **Background Janitor**

Expired keys are automatically removed from memory by the janitor goroutine (runs every 10s by default).

```go
kache.Set("key1", "value1", kache.Options{TTL: 1})
time.Sleep(15 * time.Second)
_, exists := kache.Get("key1")
// ❌ key should be cleaned up by janitor
```

---

## ⌛️ TTL vs No TTL

* ✅ **With TTL**: Keys will expire after `n` seconds.

  ```go
  kache.Set("key", "value", kache.Options{TTL: 30}) // expires in 30s
  ```

* 🟢 **Without TTL** (default: never expires): You can pass `kache.Options{TTL: 0}` or skip TTL logic in your implementation.

  ```go
  kache.Set("key", "value", kache.Options{}) // no expiration
  ```

> You can optionally enhance `Set()` to interpret `0` TTL as "never expire."

---

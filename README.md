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

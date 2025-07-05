package kache

import (
	"sync"
)

var ShardsCount = 16

type Shard struct {
	mu    sync.RWMutex
	store map[string]string
}

type ShardedMap []*Shard

func NewKache() *ShardedMap {
	shards := make(ShardedMap, ShardsCount)
	for i := range shards {
		shardsMap := make(map[string]string)
		shards[i] = &Shard{store: shardsMap}
	}
	return &shards
}

func (shardmap *ShardedMap) Set(key, value string) {
	kache := (*shardmap).getShard(key)
	kache.mu.Lock()
	defer kache.mu.Unlock()
	kache.store[key] = value
}

func (shardmap *ShardedMap) Get(key string) (string, bool) {
	kache := (*shardmap).getShard(key)
	kache.mu.RLock()
	defer kache.mu.RUnlock()

	if _, ok := kache.store[key]; !ok {
		return "", false
	}

	return kache.store[key], true
}

func (shardmap *ShardedMap) Delete(key string) {
	kache := (*shardmap).getShard(key)
	kache.mu.Lock()
	defer kache.mu.Unlock()
	delete(kache.store, key)
}

func (shardmap *ShardedMap) Exists(key string) bool {
	kache := (*shardmap).getShard(key)
	kache.mu.RLock()
	defer kache.mu.RUnlock()
	_, ok := kache.store[key]
	return ok
}

func (shardmap *ShardedMap) Flush() {
	for _, shard := range *shardmap {
		shard.mu.Lock()
		shard.store = make(map[string]string)
		shard.mu.Unlock()
	}
}

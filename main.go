package kache

import (
	"sync"
	"time"
)

var ShardsCount = 16

type Item struct {
	value string
	ttl   time.Time
}

type Shard struct {
	mu    sync.RWMutex
	store map[string]Item
}

type ShardedMap []*Shard

func NewKache() *ShardedMap {
	shards := make(ShardedMap, ShardsCount)
	for i := range shards {
		shardsMap := make(map[string]Item)
		shards[i] = &Shard{store: shardsMap}
	}
	go shards.StartJanitor(10 * time.Second) // Start janitor to clean up expired items every 5 seconds
	return &shards
}

type Options struct {
	TTL int64
}

func (shardmap *ShardedMap) Set(key, value string, options Options) {
	kache := (*shardmap).getShard(key)
	kache.mu.Lock()
	defer kache.mu.Unlock()

	ttl := time.Now().Add(time.Duration(options.TTL) * time.Second)
	if options.TTL <= 0 {
		ttl = time.Time{}
	}
	kache.store[key] = Item{
		value: value,
		ttl:   ttl,
	}
}

func (shardmap *ShardedMap) Get(key string) (string, bool) {
	kache := (*shardmap).getShard(key)
	kache.mu.RLock()
	item, ok := kache.store[key]
	kache.mu.RUnlock()

	if !ok {
		return "", false
	}

	if time.Now().After(kache.store[key].ttl) {
		kache.mu.Lock()
		delete(kache.store, key)
		kache.mu.Unlock()
		return "", false
	}

	return item.value, true
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
	_, ok := kache.store[key]
	kache.mu.RUnlock()

	if time.Now().After(kache.store[key].ttl) {
		kache.mu.Lock()
		delete(kache.store, key)
		kache.mu.Unlock()
		return false
	}

	return ok
}

func (shardmap *ShardedMap) Flush() {
	for _, shard := range *shardmap {
		shard.mu.Lock()
		shard.store = make(map[string]Item)
		shard.mu.Unlock()
	}
}

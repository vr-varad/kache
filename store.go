package kache

import (
	"container/list"
	"time"
)

func (shardmap *shardedMap) Set(key, value string, options options) {
	kache := getShard(key, *shardmap)
	kache.mu.Lock()
	defer kache.mu.Unlock()

	ttl := time.Now().Add(time.Duration(options.TTL) * time.Second)
	if options.TTL <= 0 {
		ttl = time.Time{}
	}

	if item, exists := kache.store[key]; exists {

		kache.lru.MoveToFront(kache.index[key])

		item.value = value
		item.ttl = ttl
		kache.store[key] = item
		return
	}

	if len(kache.store) >= MaxEntries {
		lruElement := kache.lru.Back()
		if lruElement != nil {
			lruKey := lruElement.Value.(string)
			delete(kache.store, lruKey)
			kache.lru.Remove(lruElement)
			delete(kache.index, lruKey)
		}
	}

	kache.index[key] = kache.lru.PushFront(key)

	kache.store[key] = item{
		value: value,
		ttl:   ttl,
	}
}

func (shardmap *shardedMap) Get(key string) (string, bool) {
	kache := getShard(key, *shardmap)
	kache.mu.Lock()
	item, ok := kache.store[key]
	defer kache.mu.Unlock()

	if !ok {
		return "", false
	}

	if time.Now().After(kache.store[key].ttl) {
		delete(kache.store, key)
		kache.lru.Remove(kache.index[key])
		delete(kache.index, key)
		return "", false
	}

	kache.lru.MoveToFront(kache.index[key])

	return item.value, true
}

func (shardmap *shardedMap) Delete(key string) {
	kache := getShard(key, *shardmap)
	kache.mu.Lock()
	defer kache.mu.Unlock()
	if _, ok := kache.index[key]; !ok {
		return
	}
	kache.lru.Remove(kache.index[key])
	delete(kache.store, key)
	delete(kache.index, key)
}

func (shardmap *shardedMap) Exists(key string) bool {
	kache := getShard(key, *shardmap)
	kache.mu.Lock()
	_, ok := kache.store[key]
	defer kache.mu.Unlock()

	if !ok {
		return false
	}
	if time.Now().After(kache.store[key].ttl) {
		kache.lru.Remove(kache.index[key])
		delete(kache.store, key)
		delete(kache.index, key)
		return false
	}

	kache.lru.MoveToFront(kache.index[key])

	return true
}

func (shardmap *shardedMap) Flush() {
	for _, shard := range *shardmap {
		shard.mu.Lock()
		shard.store = make(map[string]item)
		shard.lru = list.New()
		shard.index = make(map[string]*list.Element)
		shard.mu.Unlock()
	}
}

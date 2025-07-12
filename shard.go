package kache

import (
	"container/list"
	"sync"
	"time"
)

var ShardsCount = 16  // Number of shards in the sharded map
var MaxEntries = 1000 // Maximum number of entries per shard

type item struct {
	value string
	ttl   time.Time
}

type shard struct {
	mu    sync.RWMutex
	store map[string]item
	lru   *list.List
	index map[string]*list.Element
}

type shardedMap []*shard

func NewKache() *shardedMap {
	shards := make(shardedMap, ShardsCount)
	for i := range shards {
		shardsMap := make(map[string]item)
		shards[i] = &shard{
			store: shardsMap,
			lru:   list.New(),
			index: make(map[string]*list.Element),
		}
	}
	go startJanitor(10*time.Second, &shards) // Start janitor to clean up expired items every 5 seconds
	return &shards
}

type Options struct {
	TTL int64
}

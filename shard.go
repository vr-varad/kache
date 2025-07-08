package kache

import (
	"container/list"
	"sync"
	"time"
)

var ShardsCount = 16
var MaxEntries = 1000 // Maximum number of entries per shard

type Item struct {
	value string
	ttl   time.Time
}

type Shard struct {
	mu    sync.RWMutex
	store map[string]Item
	lru   *list.List
	index map[string]*list.Element
}

type ShardedMap []*Shard

func NewKache() *ShardedMap {
	shards := make(ShardedMap, ShardsCount)
	for i := range shards {
		shardsMap := make(map[string]Item)
		shards[i] = &Shard{
			store: shardsMap,
			lru:   list.New(),
			index: make(map[string]*list.Element),
		}
	}
	go shards.StartJanitor(10 * time.Second) // Start janitor to clean up expired items every 5 seconds
	return &shards
}

type Options struct {
	TTL int64
}

package kache

import (
	"crypto/sha1"
	"time"
)

func (sm ShardedMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	n := int(checksum[15])
	return n % len(sm)
}

func (sm ShardedMap) getShard(key string) *Shard {
	idx := sm.getShardIndex(key)
	return sm[idx]
}

func (sm *ShardedMap) StartJanitor(interval time.Duration) {
	for _, shard := range *sm {
		go shard.runJanitor(interval)
	}
}

func (s *Shard) runJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		s.mu.Lock()
		for key, item := range s.store {
			if time.Now().After(item.ttl) {
				delete(s.store, key)
			}
		}
		s.mu.Unlock()
	}
}

package kache

import (
	"time"
)

func startJanitor(interval time.Duration, sm *shardedMap) {
	for _, shard := range *sm {
		go runJanitor(interval, shard)
	}
}

func runJanitor(interval time.Duration, s *shard) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for key, item := range s.store {
			if item.ttl.IsZero() {
				continue // No TTL set, skip this item
			}
			if time.Now().After(item.ttl) {
				delete(s.store, key)
			}
		}
		s.mu.Unlock()
	}
}

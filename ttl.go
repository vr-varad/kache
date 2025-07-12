package kache

import "time"

func StartJanitor(interval time.Duration, sm *shardedMap) {
	for _, shard := range *sm {
		go runJanitor(interval, shard)
	}
}

func runJanitor(interval time.Duration, s *shard) {
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

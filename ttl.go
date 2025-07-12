package kache

import "time"

func (sm *shardedMap) StartJanitor(interval time.Duration) {
	for _, shard := range *sm {
		go shard.runJanitor(interval)
	}
}

func (s *shard) runJanitor(interval time.Duration) {
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

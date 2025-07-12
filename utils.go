package kache

func (s *shard) Evict(key string) {
	delete(s.store, key)
	s.lru.Remove(s.index[key])
	delete(s.index, key)
}

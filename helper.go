package kache

import "crypto/sha1"

func (sm ShardedMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	n := int(checksum[15])
	return n % len(sm)
}

func (sm ShardedMap) getShard(key string) *Shard {
	idx := sm.getShardIndex(key)
	return sm[idx]
}

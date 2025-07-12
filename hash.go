package kache

import "crypto/sha1"

func (sm shardedMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	n := int(checksum[15])
	return n % len(sm)
}

func (sm shardedMap) getShard(key string) *shard {
	idx := sm.getShardIndex(key)
	return sm[idx]
}

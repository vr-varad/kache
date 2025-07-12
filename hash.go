package kache

import "crypto/sha1"

func getShardIndex(key string, sm shardedMap) int {
	checksum := sha1.Sum([]byte(key))
	n := int(checksum[15])
	return n % len(sm)
}

func getShard(key string, sm shardedMap) *shard {
	idx := getShardIndex(key, sm)
	return sm[idx]
}

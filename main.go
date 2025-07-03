package kache

import (
	"sync"
)

type Kache struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewKache() *Kache {
	return &Kache{
		store: map[string]string{},
	}
}

func (kache *Kache) Set(key, value string) {
	kache.mu.Lock()
	defer kache.mu.Unlock()
	kache.store[key] = value
}

func (kache *Kache) Get(key string) (string, bool) {
	kache.mu.RLock()
	defer kache.mu.RUnlock()

	if _, ok := kache.store[key]; !ok {
		return "", false
	}

	return kache.store[key], true
}

func (kache *Kache) Delete(key string) {
	kache.mu.Lock()
	defer kache.mu.Unlock()
	delete(kache.store, key)
}

func (kache *Kache) Exists(key string) bool {
	kache.mu.RLock()
	defer kache.mu.RUnlock()
	_, ok := kache.store[key]
	return ok
}

func (kache *Kache) Size() int {
	kache.mu.RLock()
	defer kache.mu.RUnlock()
	return len(kache.store)
}

func (kache *Kache) Flush() {
	kache.mu.Lock()
	defer kache.mu.Unlock()
	clear(kache.store)
}

func (kache *Kache) Keys() []string {
	kache.mu.RLock()
	defer kache.mu.RUnlock()
	keys := make([]string, 0, len(kache.store))
	for key := range kache.store {
		keys = append(keys, key)
	}

	return keys
}

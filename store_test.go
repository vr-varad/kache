package kache

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestKache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1", Options{TTL: 10})

		value, exists := kache.Get("key1")
		if !exists || value != "value1" {
			t.Errorf("Expected value1, got %s", value)
		}
	})
	t.Run("Get Non-Existent Key", func(t *testing.T) {
		kache := NewKache()
		_, exists := kache.Get("nonexistent")
		if exists {
			t.Error("Expected key to not exist")
		}
	})
	t.Run("Delete Key", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1", Options{TTL: 10})
		kache.Delete("key1")

		_, exists := kache.Get("key1")
		if exists {
			t.Error("Expected key to be deleted")
		}
	})
	t.Run("Exists Key", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1", Options{TTL: 10})

		if !kache.Exists("key1") {
			t.Error("Expected key to exist")
		}
		if kache.Exists("nonexistent") {
			t.Error("Expected key to not exist")
		}
	})
	t.Run("Flush Kache", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1", Options{TTL: 10})
		kache.Set("key2", "value2", Options{TTL: 10})
		kache.Flush()
		if kache.Exists("key1") || kache.Exists("key2") {
			t.Error("Expected kache to be empty after flush")
		}
	})
	t.Run("Time-to-Live (TTL)", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1", Options{TTL: 1}) // 1 second TTL

		value, exists := kache.Get("key1")
		if !exists || value != "value1" {
			t.Errorf("Expected value1, got %s", value)
		}
		// Wait for TTL to expire
		time.Sleep(2 * time.Second)
		_, exists = kache.Get("key1")
		if exists {
			t.Error("Expected key to be expired after TTL")
		}
	})

	t.Run("Janitor", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1", Options{TTL: 1}) // 1 second TTL
		kache.Set("key21", "value21", Options{})     // 10 seconds TTL

		// Wait for janitor to clean up expired items
		fmt.Println("Waiting for janitor to clean up expired items...")
		time.Sleep(15 * time.Second)

		_, exists := kache.Get("key1")
		if exists {
			t.Error("Expected key is not cleaned up by janitor")
		}

		value, exists := kache.Get("key21")
		if !exists || value != "value21" {
			t.Errorf("Expected value21, got %s", value)
		}
	})
}

func BenchmarkKacheSet(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
}
func BenchmarkKacheGet(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	for i := 0; i < b.N; i++ {
		kache.Get("key" + strconv.Itoa(i))
	}
}
func BenchmarkKacheDelete(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	for i := 0; i < b.N; i++ {
		kache.Delete("key" + strconv.Itoa(i))
	}
}
func BenchmarkKacheExists(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	for i := 0; i < b.N; i++ {
		kache.Exists("key" + strconv.Itoa(i))
	}
}
func BenchmarkKacheFlush(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	for i := 0; i < b.N; i++ {
		kache.Flush()
	}
}

func BenchmarkKacheConcurrentSet(b *testing.B) {
	kache := NewKache()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key" + strconv.Itoa(rand.Intn(10000))
			value := "value" + strconv.Itoa(rand.Intn(10000))
			kache.Set(key, value, Options{TTL: 10})
		}
	})
}

func BenchmarkKacheConcurrentGet(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key" + strconv.Itoa(rand.Intn(10000))
			kache.Get(key)
		}
	})
}

func BenchmarkKacheConcurrentDelete(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key" + strconv.Itoa(rand.Intn(10000))
			kache.Delete(key)
		}
	})
}

func BenchmarkKacheConcurrentExists(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := "key" + strconv.Itoa(rand.Intn(10000))
			kache.Exists(key)
		}
	})
}
func BenchmarkKacheConcurrentFlush(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i), Options{TTL: 10})
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			kache.Flush()
		}
	})
}

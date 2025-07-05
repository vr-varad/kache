package kache

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestKache(t *testing.T) {
	t.Run("Set and Get", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1")

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
		kache.Set("key1", "value1")
		kache.Delete("key1")

		_, exists := kache.Get("key1")
		if exists {
			t.Error("Expected key to be deleted")
		}
	})
	t.Run("Exists Key", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1")

		if !kache.Exists("key1") {
			t.Error("Expected key to exist")
		}
		if kache.Exists("nonexistent") {
			t.Error("Expected key to not exist")
		}
	})
	t.Run("Flush Kache", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1")
		kache.Set("key2", "value2")
		kache.Flush()
		if kache.Exists("key1") || kache.Exists("key2") {
			t.Error("Expected kache to be empty after flush")
		}
	})
}

func BenchmarkKacheSet(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
}
func BenchmarkKacheGet(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		kache.Get("key" + strconv.Itoa(i))
	}
}
func BenchmarkKacheDelete(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		kache.Delete("key" + strconv.Itoa(i))
	}
}
func BenchmarkKacheExists(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		kache.Exists("key" + strconv.Itoa(i))
	}
}
func BenchmarkKacheFlush(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
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
			kache.Set(key, value)
		}
	})
}

func BenchmarkKacheConcurrentGet(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
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
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
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
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
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
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			kache.Flush()
		}
	})
}

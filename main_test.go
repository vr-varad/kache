package kache

import (
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
	t.Run("Size of Kache", func(t *testing.T) {
		kache := NewKache()
		if kache.Size() != 0 {
			t.Error("Expected size to be 0")
		}

		kache.Set("key1", "value1")
		if kache.Size() != 1 {
			t.Error("Expected size to be 1")
		}

		kache.Set("key2", "value2")
		if kache.Size() != 2 {
			t.Error("Expected size to be 2")
		}

		kache.Delete("key1")
		if kache.Size() != 1 {
			t.Error("Expected size to be 1 after deletion")
		}
	})
	t.Run("Flush Kache", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1")
		kache.Set("key2", "value2")
		kache.Flush()
		if kache.Size() != 0 {
			t.Error("Expected size to be 0 after flush")
		}
	})
	t.Run("Keys in Kache", func(t *testing.T) {
		kache := NewKache()
		kache.Set("key1", "value1")
		kache.Set("key2", "value2")
		keys := kache.Keys()
		if len(keys) != 2 {
			t.Error("Expected 2 keys in kache")
		}
		if keys[0] != "key1" && keys[0] != "key2" {
			t.Error("Expected keys to contain key1 or key2")
		}
		if keys[1] != "key1" && keys[1] != "key2" {
			t.Error("Expected keys to contain key1 or key2")
		}
		if kache.Exists("key1") && kache.Exists("key2") {
			t.Log("Keys test passed")
		} else {
			t.Error("Expected both keys to exist in kache")
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
func BenchmarkKacheSize(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		kache.Size()
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
func BenchmarkKacheKeys(b *testing.B) {
	kache := NewKache()
	for i := 0; i < b.N; i++ {
		kache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		kache.Keys()
	}
}

package cache_test

import (
	"testing"

	"concurrency/cache"
)

func TestNewCache(t *testing.T) {
	cache_ := cache.NewCache[string, string](5)

	if cache_ == nil {
		t.Fatal("error: got a nil cache")
	}

	if cache_.Limit != 5 {
		t.Errorf("Expected cache limit to 15, got %d", cache_.Limit)
	}
}

func TestPutAndGet(t *testing.T) {
	cache_ := cache.NewCache[string, int](3)

	cache_.Put("a", 1)
	cache_.Put("b", 2)
	cache_.Put("c", 3)

	val, ok := cache_.Get("a")

	if !ok || val != 1 {
		t.Errorf("Expected to get 1, got %d, ok=%v", val, ok)
	}

	val, ok = cache_.Get("b")
	if !ok || val != 2 {
		t.Errorf("Expected to get 2, got %d, ok=%v", val, ok)
	}

	val, ok = cache_.Get("c")
	if !ok || val != 3 {
		t.Errorf("Expected to get 3, got %d, ok=%v", val, ok)
	}
}

func TestGetNonExistent(t *testing.T) {
	cache_ := cache.NewCache[string, string](2)
	val, ok := cache_.Get("nonexistent")
	if val != "" {
		t.Errorf("Expected to get '', got %s, ok=%v", val, ok)
	}
}

func TestLRUEviction(t *testing.T) {
	cache_ := cache.NewCache[string, int](3)

	cache_.Put("a", 1)
	cache_.Put("b", 2)
	cache_.Put("c", 3)

	cache_.Put("d", 4)

	_, ok := cache_.Get("a")
	if ok {
		t.Error("Expected 'a' to be evicted, but it's still in cache")
	}
}

func TestGetRefreshesItem(t *testing.T) {
	cache_ := cache.NewCache[string, int](3)

	cache_.Put("a", 1)
	cache_.Put("b", 2)
	cache_.Put("c", 3)

	cache_.Get("a")

	cache_.Put("d", 4)

	_, ok := cache_.Get("b")
	if ok {
		t.Error("Expected 'b' to be evicted")
	}

	val, ok := cache_.Get("a")
	if !ok || val != 1 {
		t.Error("Expected 'a' to still be in cache")
	}
}

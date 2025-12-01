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

	cacheLimit := cache_.GetLimit()

	if cacheLimit != 5 {
		t.Errorf("Expected cache limit to 15, got %d", cacheLimit)
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

// this test case also acounts for evicted key, hence the second
// cache.Get("a")
func TestReadCount(t *testing.T) {
	cache_ := cache.NewCache[string, int](2)
	cache_.Put("a", 1)
	cache_.Put("b", 2)
	cache_.Get("a")

	cache_.Put("c", 3)
	cache_.Get("b")
	cache_.Get("c")
	cache_.Get("a")

	stats := cache_.GetStats()
	readCounts := stats.CacheLevelStats.Reads
	if readCounts != 4 {
		t.Errorf("expected 4 total reads, got %d", readCounts)
	}
}

func TestPutCount(t *testing.T) {
	cache_ := cache.NewCache[string, int](4)
	cache_.Put("a", 1)
	cache_.Put("b", 2)
	cache_.Put("c", 3)

	stats := cache_.GetStats()
	writeCounts := stats.CacheLevelStats.Writes
	if writeCounts != 3 {
		t.Errorf("expected 3 total writes, got %d", writeCounts)
	}
}

func TestHitRate(t *testing.T) {
	cache_ := cache.NewCache[int, string](3)
	cache_.Put(1, "z")
	cache_.Put(2, "y")
	cache_.Put(3, "x")

	cache_.Get(1)
	cache_.Get(2)
	cache_.Get(3)

	cache_.Get(4)
	cache_.Get(5)
	cache_.Get(6)

	stats := cache_.GetStats()
	hitRate := stats.CacheLevelStats.HitRate

	if hitRate != 0.5 {
		t.Errorf("expected hit rate value of 0.5, got %f", hitRate)
	}
}

func TestAverageReads(t *testing.T) {
	cache_ := cache.NewCache[int, string](2)
	cache_.Put(1, "z")
	cache_.Put(2, "y")
	cache_.Put(3, "x")

	cache_.Get(1)
	cache_.Get(2)
	cache_.Get(3)
	cache_.Get(4)
	cache_.Get(5)
	cache_.Get(2)

	stats := cache_.GetStats()
	avgReads := stats.CacheLevelStats.AverageReads

	if avgReads != 1.5 {
		t.Errorf("expected avg reads of 1.5, got %f", avgReads)
	}
}

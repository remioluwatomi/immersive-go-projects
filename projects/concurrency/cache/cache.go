package cache

import (
	"sync"
)

type Node[K comparable, V any] struct {
	prev *Node[K, V]
	next *Node[K, V]
	key  K
	val  V
}

type EntryLevelStats[K comparable] struct {
	HitCount int
	key      K
}

func createLinkedDummyNodes[K comparable, V any]() (*Node[K, V], *Node[K, V]) {
	head := &Node[K, V]{}
	tail := &Node[K, V]{}

	head.next = tail
	tail.prev = head

	return head, tail
}

type CacheLevelStats struct {
	Reads        int
	Writes       int
	HitRate      float64
	AverageReads float64
	cacheHit     int
}

type Cache[K comparable, V any] struct {
	cacheMap        map[K]*Node[K, V]
	limit           int
	head            *Node[K, V]
	tail            *Node[K, V]
	stats           CacheLevelStats
	entryLevelStats map[K]EntryLevelStats[K]
	mu              sync.RWMutex
	statsMu         sync.Mutex
}

func (c *Cache[K, V]) unlinkNode(n *Node[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (c *Cache[K, V]) addNode(n *Node[K, V]) {
	prevNode := c.tail.prev
	prevNode.next = n

	n.prev = prevNode
	n.next = c.tail

	c.tail.prev = n
}

func (c *Cache[K, V]) GetLimit() int {
	return c.limit
}

func (c *Cache[K, V]) moveExistingNodeToTail(n *Node[K, V]) {
	c.unlinkNode(n)
	c.addNode(n)
}

func (c *Cache[K, V]) incrementCacheStatReadOrWrite(statType string) {
	switch statType {
	case "READ":
		c.stats.Reads++
		return
	case "WRITE":
		c.stats.Writes++
		return
	}
}

func (c *Cache[K, V]) Put(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	node, ok := c.cacheMap[key]

	c.statsMu.Lock()
	c.incrementCacheStatReadOrWrite("WRITE")
	c.statsMu.Unlock()

	if ok {
		node.val = val
		c.moveExistingNodeToTail(node)
		return
	}

	newNode := &Node[K, V]{
		key: key,
		val: val,
	}

	c.cacheMap[key] = newNode
	c.entryLevelStats[key] = EntryLevelStats[K]{
		key: key,
	}

	c.addNode(newNode)

	if len(c.cacheMap) > c.limit {
		firstOutNode := c.head.next
		c.unlinkNode(firstOutNode)
		delete(c.cacheMap, firstOutNode.key)
		delete(c.entryLevelStats, firstOutNode.key)
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	node, ok := c.cacheMap[key]
	c.mu.RUnlock()

	// this is basically the total lookup/write
	c.statsMu.Lock()
	c.incrementCacheStatReadOrWrite("READ")
	c.statsMu.Unlock()

	var val V
	if !ok {
		return val, false
	}

	c.mu.Lock()

	if nodeRecheck, ok := c.cacheMap[key]; !ok || nodeRecheck != node {
		c.mu.Unlock()
		return val, false
	}

	c.moveExistingNodeToTail(node)
	c.mu.Unlock()

	c.statsMu.Lock()
	c.stats.cacheHit++

	entryStat := c.entryLevelStats[key]
	entryStat.HitCount++
	c.entryLevelStats[key] = entryStat
	c.statsMu.Unlock()

	return node.val, true
}

type CacheStats struct {
	CacheLevelStats CacheLevelStats
}

func (c *Cache[K, V]) averageReads() float64 {
	if len(c.cacheMap) == 0 {
		return 0
	}

	totalEntryHitCount := 0
	for _, stats := range c.entryLevelStats {
		totalEntryHitCount += stats.HitCount
	}

	return float64(totalEntryHitCount) / float64(len(c.cacheMap))
}

func (c *Cache[K, V]) GetStats() CacheStats {
	// compute hit HitRate i.e cacheHit divided by total Reads/total lookup
	c.statsMu.Lock()
	defer c.statsMu.Unlock()

	var hitRate float64
	if c.stats.Reads > 0 {
		hitRate = float64(c.stats.cacheHit) / float64(c.stats.Reads)
	}

	return CacheStats{
		CacheLevelStats: CacheLevelStats{
			Reads:        c.stats.Reads,
			Writes:       c.stats.Writes,
			HitRate:      hitRate,
			AverageReads: c.averageReads(),
		},
	}
}

func NewCache[K comparable, V any](limit int) *Cache[K, V] {
	if limit < 1 {
		panic("cache limit must be at least 1")
	}

	head, tail := createLinkedDummyNodes[K, V]()

	return &Cache[K, V]{
		cacheMap:        make(map[K]*Node[K, V]),
		limit:           limit,
		head:            head,
		tail:            tail,
		stats:           CacheLevelStats{},
		entryLevelStats: make(map[K]EntryLevelStats[K]),
		mu:              sync.RWMutex{},
		statsMu:         sync.Mutex{},
	}
}

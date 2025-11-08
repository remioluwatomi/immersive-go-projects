package cache

type Node[K comparable, V any] struct {
	prev *Node[K, V]
	next *Node[K, V]
	key  K
	val  V
}

func createLinkedDummyNodes[K comparable, V any]() (*Node[K, V], *Node[K, V]) {
	head := &Node[K, V]{}
	tail := &Node[K, V]{}

	head.next = tail
	tail.prev = head

	return head, tail
}

type CacheLevelStats struct {
	Reads  int
	Writes int
}

type Cache[K comparable, V any] struct {
	cacheMap map[K]*Node[K, V]
	Limit    int
	head     *Node[K, V]
	tail     *Node[K, V]
	stats    CacheLevelStats
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
	node, ok := c.cacheMap[key]
	defer c.incrementCacheStatReadOrWrite("WRITE")

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
	c.addNode(newNode)

	if len(c.cacheMap) > c.Limit {
		firstOutNode := c.head.next
		c.unlinkNode(firstOutNode)
		delete(c.cacheMap, firstOutNode.key)
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	node, ok := c.cacheMap[key]
	defer c.incrementCacheStatReadOrWrite("READ")

	if !ok {
		var val V
		return val, false
	}

	c.moveExistingNodeToTail(node)
	return node.val, true
}

type CacheStats struct {
	CacheLevelStats CacheLevelStats
}

func (c *Cache[K, V]) GetStats() CacheStats {
	return CacheStats{
		CacheLevelStats: CacheLevelStats{
			Reads:  c.stats.Reads,
			Writes: c.stats.Writes,
		},
	}
}

func NewCache[K comparable, V any](limit int) *Cache[K, V] {
	if limit < 1 {
		panic("cache limit must be at least 1")
	}

	head, tail := createLinkedDummyNodes[K, V]()

	return &Cache[K, V]{
		cacheMap: make(map[K]*Node[K, V]),
		Limit:    limit,
		head:     head,
		tail:     tail,
		stats:    CacheLevelStats{},
	}
}

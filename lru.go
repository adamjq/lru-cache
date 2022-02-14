package lru

import (
	"errors"
	"sync"
)

type LRUCache[K comparable, V any] struct {
	maxCapacity int
	cache       map[K]*node[K, V]
	head        *node[K, V]
	tail        *node[K, V]
	mu          sync.RWMutex
}

type node[K comparable, V any] struct {
	key   *K
	value *V
	next  *node[K, V]
	prev  *node[K, V]
}

func New[K comparable, V any](capacity int) (*LRUCache[K, V], error) {
	if capacity < 1 {
		return nil, errors.New("capacity must be greater than 0")
	}

	head, tail := node[K, V]{}, node[K, V]{}
	head.next = &tail
	tail.prev = &head

	return &LRUCache[K, V]{
		maxCapacity: capacity,
		cache:       make(map[K]*node[K, V]),
		head:        &head,
		tail:        &tail,
	}, nil
}

// Get returns a value from the cache
func (lc *LRUCache[K, V]) Get(key K) *V {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	n, exists := lc.cache[key]
	if exists {
		lc.remove(n)
		lc.insert(n)
		return n.value
	}

	return nil
}

// Put stores a value in the cache and evicts the oldest value if the cache is at maximum capacity
func (lc *LRUCache[K, V]) Put(key K, value V) {
	lc.mu.Lock()

	n, exists := lc.cache[key]
	if exists {
		lc.remove(n)
	}

	newNode := node[K, V]{
		key:   &key,
		value: &value,
	}
	lc.cache[key] = &newNode
	lc.insert(&newNode)

	if len(lc.cache) > lc.maxCapacity {
		lc.evict()
	}

	lc.mu.Unlock()
}

// insert adds a node to the head of the doubly linked list
func (lc *LRUCache[K, V]) insert(n *node[K, V]) {
	headNode, firstNode := lc.head, lc.head.next

	headNode.next = n
	firstNode.prev = n

	n.prev = headNode
	n.next = firstNode
}

// remove unlinks a node from a doubly linked list
func (lc *LRUCache[K, V]) remove(n *node[K, V]) {
	prevNode, nextNode := n.prev, n.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
}

// evict removes the least frequently used key from the cache
func (lc *LRUCache[K, V]) evict() {
	lruNode := lc.tail.prev
	delete(lc.cache, *lruNode.key)
	lc.remove(lruNode)
}

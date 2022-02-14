package lru

import (
	"errors"
	"sync"
)

type LRUCache[K comparable] struct {
	maxCapacity int
	cache       map[K]*node[K]
	head        *node[K]
	tail        *node[K]
	mu          sync.RWMutex
}

type node[K comparable] struct {
	key   *K
	value *string
	next  *node[K]
	prev  *node[K]
}

func New[K comparable](capacity int) (*LRUCache[K], error) {
	if capacity < 1 {
		return nil, errors.New("capacity must be greater than 0")
	}

	head, tail := node[K]{}, node[K]{}
	head.next = &tail
	tail.prev = &head

	return &LRUCache[K]{
		maxCapacity: capacity,
		cache:       make(map[K]*node[K]),
		head:        &head,
		tail:        &tail,
	}, nil
}

// Get returns a value from the cache
func (lc *LRUCache[K]) Get(key K) *string {
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
func (lc *LRUCache[K]) Put(key K, value string) {
	lc.mu.Lock()

	n, exists := lc.cache[key]
	if exists {
		lc.remove(n)
	}

	newNode := node[K]{
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
func (lc *LRUCache[K]) insert(n *node[K]) {
	headNode, firstNode := lc.head, lc.head.next

	headNode.next = n
	firstNode.prev = n

	n.prev = headNode
	n.next = firstNode
}

// remove unlinks a node from a doubly linked list
func (lc *LRUCache[K]) remove(n *node[K]) {
	prevNode, nextNode := n.prev, n.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
}

// evict removes the least frequently used key from the cache
func (lc *LRUCache[K]) evict() {
	lruNode := lc.tail.prev
	delete(lc.cache, *lruNode.key)
	lc.remove(lruNode)
}

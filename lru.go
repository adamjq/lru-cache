package lru

import (
	"errors"
	"sync"
)

type LRUCache struct {
	maxCapacity int
	cache       map[string]*node
	head        *node
	tail        *node
	mu          sync.RWMutex
}

type node struct {
	key   *string
	value *string
	next  *node
	prev  *node
}

func New(capacity int) (*LRUCache, error) {
	if capacity < 1 {
		return nil, errors.New("capacity must be greater than 0")
	}

	head, tail := node{}, node{}
	head.next = &tail
	tail.prev = &head

	return &LRUCache{
		maxCapacity: capacity,
		cache:       make(map[string]*node),
		head:        &head,
		tail:        &tail,
	}, nil
}

// Get returns a value from the cache
func (lc *LRUCache) Get(key string) *string {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	n, exists := lc.cache[key]
	if exists {

		// move to front of doubly linked list
		lc.remove(n)
		lc.insert(n)
		return n.value
	}

	return nil
}

// Put stores a value in the cache and evicts the oldest value if the cache is at maximum capacity
func (lc *LRUCache) Put(key string, value string) {
	lc.mu.Lock()

	n, exists := lc.cache[key]
	if exists {
		lc.remove(n)
	}

	newNode := node{
		key:   &key,
		value: &value,
	}
	lc.cache[key] = &newNode
	lc.insert(&newNode)

	// evict oldest key if cache is over capacity
	if len(lc.cache) > lc.maxCapacity {
		lruNode := lc.tail.prev
		delete(lc.cache, *lruNode.key)
		lc.remove(lruNode)
	}

	lc.mu.Unlock()
}

// insert adds a node to the head of the doubly linked list
func (lc *LRUCache) insert(n *node) {
	headNode, firstNode := lc.head, lc.head.next

	headNode.next = n
	firstNode.prev = n

	n.prev = headNode
	n.next = firstNode
}

// remove deletes a node from a doubly linked list
func (lc *LRUCache) remove(n *node) {
	prevNode, nextNode := n.prev, n.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
}

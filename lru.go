package lru

import (
	"fmt"

	"github.com/ndisidore/lru-generic/internal/list"
)

type (
	config struct {
		cap uint16
	}

	Cache[K comparable, V any] struct {
		list    *list.DLL[entry[K, V]]
		entries map[K]*list.Node[entry[K, V]]
		config  config
	}

	entry[K comparable, V any] struct {
		key   K
		value V
	}
)

// New initializes a new lru cache with the provided capacity.
func NewCache[K comparable, V any](capacity uint16) *Cache[K, V] {
	return &Cache[K, V]{
		list:    list.NewDoublyLinkedList[entry[K, V]](),
		entries: make(map[K]*list.Node[entry[K, V]], 0),
		config: config{
			cap: capacity,
		},
	}
}

// Len is the number of key value pairs in the cache.
func (c *Cache[K, V]) Len() int {
	return c.list.Len()
}

// Set sets an entry in the cache using the provided key and value. Sets are treated
// as accesses and thus the new entry is prepended to the front of the LRU access list.
func (c *Cache[K, V]) Set(key K, val V) error {
	// If the key exists, it will be treated as an updated operation
	if existing, ok := c.entries[key]; ok {
		if err := c.list.MoveToFront(existing); err != nil {
			return fmt.Errorf("unable to move existing entry: %w", err)
		}
		return nil
	}

	node, err := c.list.InsertAtFront(entry[K, V]{
		key:   key,
		value: val,
	})
	if err != nil {
		return fmt.Errorf("could not insert node at front of list: %w", err)
	}

	c.entries[key] = node

	if c.list.Len() > int(c.config.cap) {
		if _, err := c.list.RemoveFromBack(); err != nil {
			return fmt.Errorf("unable to evict oldest entry: %w", err)
		}
	}

	return nil
}

// Get retrieves an item from the cache. Gets are treated as accesses and the value associated
// with the provided key is bumped to the front of the LRU access list.
func (c *Cache[K, V]) Get(key K) (value V, err error) {
	entry, ok := c.entries[key]
	if !ok {
		return value, fmt.Errorf("could not find an entry for key '%v'", value)
	}

	if err := c.list.MoveToFront(entry); err != nil {
		return value, fmt.Errorf("unable to promote entry to head of list: %w", err)
	}
	return entry.Value.value, nil
}

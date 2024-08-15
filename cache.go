package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache is a simple in-memory key-value store
type Cache struct {
	Lock sync.RWMutex
	Data map[string][]byte
}

// NewCache creates a new Cache
func NewCache() *Cache {
	return &Cache{
		Data: make(map[string][]byte),
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) ([]byte, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	keyStr := string(key)
	val, ok := c.Data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)

	}
	return val, nil
}

// Set sets a value in the cache
func (c *Cache) Set(key, value []byte, ttl time.Duration) []byte {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Data[string(key)] = value

	return nil
}

// Has checks if a key exists in the cache
func (c *Cache) Has(key []byte) bool {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	_, ok := c.Data[string(key)]

	return ok
}

// Delete removes a key from the cache
func (c *Cache) Delete(key []byte) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	delete(c.Data, string(key))
	return nil
}

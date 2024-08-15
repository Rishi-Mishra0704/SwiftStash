package cache

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

// Get implements the Cacher interface
// It retrieves a value from the cache
func (c *Cache) Get(key []byte) ([]byte, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	keyStr := string(key)
	val, ok := c.Data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)

	}
	return val, nil
}

// Set implements the Cacher interface
// It sets a key-value pair in the cache
func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Data[string(key)] = value

	return nil
}

// It implements the Cacher interface
// Has checks if a key exists in the cache
func (c *Cache) Has(key []byte) bool {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	_, ok := c.Data[string(key)]

	return ok
}

// It implements the Cacher interface
// Delete removes a key from the cache
func (c *Cache) Delete(key []byte) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	delete(c.Data, string(key))
	return nil
}

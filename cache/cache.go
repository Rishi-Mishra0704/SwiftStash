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
	// Acquire a read lock to ensure concurrent safety during retrieval.
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	// Convert the byte slice key to a string for map lookup.
	keyStr := string(key)

	// Retrieve the value associated with the key from the internal data map.
	val, ok := c.Data[keyStr]
	if !ok {
		// Return an error if the key is not found.
		return nil, fmt.Errorf("key (%s) not found", keyStr)
	}

	// Return the retrieved value and a nil error if the key is present in the cache.
	return val, nil
}

// Set implements the Cacher interface
// It sets a key-value pair in the cache
func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	keyStr := string(key)

	c.Data[string(key)] = value

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			c.Lock.Lock()
			defer c.Lock.Unlock()
			delete(c.Data, keyStr)
		}()
	}
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

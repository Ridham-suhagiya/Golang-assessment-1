package cache

import "sync"

// Cache interface defines the supported cache operations
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type memoryCache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewMemoryCache initializes and returns a new Cache implementation.
func NewMemoryCache() Cache {
	return &memoryCache{
		data: make(map[string]interface{}),
	}
}

// Get retrieves a value from the cache.
func (c *memoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.data[key]
	return val, exists
}

// Set stores a value in the cache.
func (c *memoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

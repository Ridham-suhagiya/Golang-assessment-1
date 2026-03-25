package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryCache(t *testing.T) {
	c := NewMemoryCache()
	assert.NotNil(t, c)
}

func TestCacheSetAndGet(t *testing.T) {
	c := NewMemoryCache()
	
	// Test setting and retrieving a value
	c.Set("india", "data")
	val, exists := c.Get("india")
	assert.True(t, exists)
	assert.Equal(t, "data", val)

	// Test getting a non-existent key
	val, exists = c.Get("usa")
	assert.False(t, exists)
	assert.Nil(t, val)
}

func TestCacheRaceCondition(t *testing.T) {
	c := NewMemoryCache()
	var wg sync.WaitGroup

	// Run multiple goroutines setting and getting values to ensure thread safety
	iterations := 1000

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			c.Set("key", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			c.Get("key")
		}
	}()

	wg.Wait()
	// If the test doesn't panic or fail with the -race flag, thread-safety is verified.
}

package pokecache

import (
	"sync"
	"time"
)

// cacheEntry holds a cached value and the time it was added.
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache stores raw bytes keyed by a string identifier.
//
// Entries are removed by a background goroutine once they are older than the
// configured interval. It is safe for concurrent use.
type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

// NewCache creates an empty Cache and starts a goroutine that removes
// entries older than interval on every tick.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
	}

	go c.reapLoop(interval)

	return c
}

// Add stores val under key, replacing any existing value and resetting its age.
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get returns the bytes stored under key and true if the key exists.
// If the key does not exist, it returns nil and false.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

// reapLoop removes expired entries on every tick until the process exits.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap(interval)
	}
}

// reap removes every entry older than interval.
func (c *Cache) reap(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.entries {
		if time.Since(entry.createdAt) > interval {
			delete(c.entries, key)
		}
	}
}

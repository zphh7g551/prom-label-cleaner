package cache

import (
	"sync"
	"time"
)

// Entry holds a cached metrics payload with an expiry timestamp.
type Entry struct {
	Data      []byte
	FetchedAt time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the entry is past its TTL.
func (e *Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a simple in-memory store for scraped metrics payloads.
type Cache struct {
	mu  sync.RWMutex
	ttl time.Duration
	entry *Entry
}

// New creates a Cache with the given TTL.
func New(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &Cache{ttl: ttl}
}

// Set stores a new payload, replacing any existing entry.
func (c *Cache) Set(data []byte) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry = &Entry{
		Data:      data,
		FetchedAt: now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get returns the current entry and whether it is valid (non-nil and not expired).
func (c *Cache) Get() (*Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.entry == nil || c.entry.IsExpired() {
		return nil, false
	}
	return c.entry, true
}

// Invalidate clears the cached entry.
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry = nil
}

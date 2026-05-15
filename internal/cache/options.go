package cache

import "time"

// Option is a functional option for Cache.
type Option func(*Cache)

// WithTTL sets the cache TTL.
func WithTTL(ttl time.Duration) Option {
	return func(c *Cache) {
		if ttl > 0 {
			c.ttl = ttl
		}
	}
}

// NewWithOptions creates a Cache and applies the provided options.
func NewWithOptions(opts ...Option) *Cache {
	c := New(30 * time.Second)
	for _, o := range opts {
		o(c)
	}
	return c
}

package cache

import (
	"testing"
	"time"
)

func TestWithTTLApplied(t *testing.T) {
	c := NewWithOptions(WithTTL(10 * time.Second))
	if c.ttl != 10*time.Second {
		t.Errorf("expected TTL 10s, got %v", c.ttl)
	}
}

func TestWithTTLZeroIgnored(t *testing.T) {
	c := NewWithOptions(WithTTL(0))
	if c.ttl != 30*time.Second {
		t.Errorf("expected default TTL 30s, got %v", c.ttl)
	}
}

func TestNewWithOptionsNoOptions(t *testing.T) {
	c := NewWithOptions()
	if c.ttl != 30*time.Second {
		t.Errorf("expected default TTL 30s, got %v", c.ttl)
	}
}

func TestNewWithOptionsMultiple(t *testing.T) {
	c := NewWithOptions(
		WithTTL(1*time.Minute),
		WithTTL(2*time.Minute),
	)
	if c.ttl != 2*time.Minute {
		t.Errorf("expected TTL 2m, got %v", c.ttl)
	}
}

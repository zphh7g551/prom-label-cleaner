package cache

import (
	"testing"
	"time"
)

func TestSetAndGetValid(t *testing.T) {
	c := New(5 * time.Second)
	payload := []byte("hello metrics")
	c.Set(payload)

	entry, ok := c.Get()
	if !ok {
		t.Fatal("expected valid cache entry")
	}
	if string(entry.Data) != string(payload) {
		t.Errorf("got %q, want %q", entry.Data, payload)
	}
}

func TestGetEmptyCache(t *testing.T) {
	c := New(5 * time.Second)
	_, ok := c.Get()
	if ok {
		t.Fatal("expected no entry in empty cache")
	}
}

func TestExpiredEntry(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set([]byte("data"))
	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get()
	if ok {
		t.Fatal("expected expired entry to be invalid")
	}
}

func TestInvalidate(t *testing.T) {
	c := New(5 * time.Second)
	c.Set([]byte("data"))
	c.Invalidate()

	_, ok := c.Get()
	if ok {
		t.Fatal("expected cache to be empty after invalidation")
	}
}

func TestNewDefaultTTL(t *testing.T) {
	c := New(0)
	if c.ttl != 30*time.Second {
		t.Errorf("expected default TTL 30s, got %v", c.ttl)
	}
}

func TestSetOverwritesPreviousEntry(t *testing.T) {
	c := New(5 * time.Second)
	c.Set([]byte("first"))
	c.Set([]byte("second"))

	entry, ok := c.Get()
	if !ok {
		t.Fatal("expected valid cache entry")
	}
	if string(entry.Data) != "second" {
		t.Errorf("got %q, want %q", entry.Data, "second")
	}
}

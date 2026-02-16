package cache

import (
	"sort"
	"sync"
	"time"
)

// CacheItem represents a single cached entry with optional expiration.
type CacheItem struct {
	Value     string
	ExpiresAt time.Time
	HasTTL    bool
}

// KeyInfo holds metadata about a cached key.
type KeyInfo struct {
	Key       string
	ExpiresAt *time.Time
}

// Cache is a thread-safe in-memory key-value store with optional TTL support.
type Cache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

// NewCache returns an initialized Cache.
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set stores a key-value pair. If ttl is non-nil, the item expires after that duration.
// If ttl is nil, the item never expires.
func (c *Cache) Set(key, value string, ttl *time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := CacheItem{
		Value: value,
	}
	if ttl != nil {
		item.HasTTL = true
		item.ExpiresAt = time.Now().Add(*ttl)
	}
	c.items[key] = item
}

// Get returns the value for key and true if the key exists and has not expired.
// Expired keys return ("", false) but are not removed from the map.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return "", false
	}
	if item.HasTTL && time.Now().After(item.ExpiresAt) {
		return "", false
	}
	return item.Value, true
}

// Delete removes the key from the cache. Returns true if the key existed.
func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.items[key]
	if exists {
		delete(c.items, key)
	}
	return exists
}

// Keys returns all non-expired keys with their expiry information.
// Keys without a TTL have a nil ExpiresAt.
func (c *Cache) Keys() []KeyInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	keys := make([]KeyInfo, 0, len(c.items))

	for k, item := range c.items {
		if item.HasTTL && now.After(item.ExpiresAt) {
			continue
		}
		ki := KeyInfo{Key: k}
		if item.HasTTL {
			t := item.ExpiresAt
			ki.ExpiresAt = &t
		}
		keys = append(keys, ki)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Key < keys[j].Key
	})

	return keys
}

// CleanExpired removes all entries whose TTL has elapsed.
func (c *Cache) CleanExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, item := range c.items {
		if item.HasTTL && now.After(item.ExpiresAt) {
			delete(c.items, k)
		}
	}
}

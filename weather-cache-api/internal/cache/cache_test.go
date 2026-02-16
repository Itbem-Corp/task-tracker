package cache

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	c := NewCache()
	c.Set("key1", "value1", nil)

	value, ok := c.Get("key1")
	if !ok || value != "value1" {
		t.Errorf("Expected 'value1', got '%v'", value)
	}
}

func TestGetNonExistent(t *testing.T) {
	c := NewCache()
	value, ok := c.Get("nonexistent")
	if ok {
		t.Errorf("Expected false for non-existent key, got true")
	}
	if value != "" {
		t.Errorf("Expected empty string, got '%v'", value)
	}
}

func TestSetWithTTL(t *testing.T) {
	c := NewCache()
	ttl := 100 * time.Millisecond
	c.Set("ttl-key", "ttl-value", &ttl)

	value, ok := c.Get("ttl-key")
	if !ok || value != "ttl-value" {
		t.Errorf("Expected 'ttl-value' immediately after set, got '%v', ok=%v", value, ok)
	}

	time.Sleep(150 * time.Millisecond)
	value, ok = c.Get("ttl-key")
	if ok {
		t.Errorf("Expected expired key to return false, got true")
	}
}

func TestDelete(t *testing.T) {
	c := NewCache()
	c.Set("delete-key", "delete-value", nil)

	if !c.Delete("delete-key") {
		t.Error("Expected Delete to return true for existing key")
	}

	value, ok := c.Get("delete-key")
	if ok {
		t.Errorf("Expected key to be deleted, but still got '%v'", value)
	}

	if c.Delete("delete-key") {
		t.Error("Expected Delete to return false for non-existent key")
	}
}

func TestKeys(t *testing.T) {
	c := NewCache()
	c.Set("a", "value-a", nil)
	c.Set("b", "value-b", nil)
	ttl := 100 * time.Millisecond
	c.Set("c", "value-c", &ttl)

	keys := c.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	if keys[0].Key != "a" || keys[1].Key != "b" || keys[2].Key != "c" {
		t.Errorf("Keys not in expected order: %+v", keys)
	}

	if keys[0].ExpiresAt != nil || keys[1].ExpiresAt != nil {
		t.Error("Keys without TTL should have nil ExpiresAt")
	}

	if keys[2].ExpiresAt == nil {
		t.Error("Key with TTL should have non-nil ExpiresAt")
	}

	time.Sleep(150 * time.Millisecond)
	keys = c.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys after TTL expiry, got %d", len(keys))
	}
}

func TestCleanExpired(t *testing.T) {
	c := NewCache()
	c.Set("permanent", "value", nil)
	ttl := 50 * time.Millisecond
	c.Set("temporary", "value", &ttl)

	time.Sleep(100 * time.Millisecond)
	c.CleanExpired()

	_, ok := c.Get("permanent")
	if !ok {
		t.Error("Permanent key should still exist")
	}

	_, ok = c.Get("temporary")
	if ok {
		t.Error("Expired key should be removed by CleanExpired")
	}
}

func TestOverwriteKey(t *testing.T) {
	c := NewCache()
	c.Set("key", "value1", nil)
	c.Set("key", "value2", nil)

	value, ok := c.Get("key")
	if !ok || value != "value2" {
		t.Errorf("Expected 'value2' after overwrite, got '%v'", value)
	}
}

func TestThreadSafety(t *testing.T) {
	c := NewCache()
	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(idx int) {
			c.Set("key", "value", nil)
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			c.Get("key")
			done <- true
		}()
	}

	for i := 0; i < 20; i++ {
		<-done
	}

	value, ok := c.Get("key")
	if !ok {
		t.Error("Expected key to exist after concurrent operations")
	}
	if value == "" {
		t.Error("Expected non-empty value")
	}
}

func TestEmptyKey(t *testing.T) {
	c := NewCache()
	c.Set("", "value", nil)

	value, ok := c.Get("")
	if !ok {
		t.Error("Empty key should be valid")
	}
	if value != "value" {
		t.Errorf("Expected 'value', got '%v'", value)
	}
}

func TestEmptyValue(t *testing.T) {
	c := NewCache()
	c.Set("key", "", nil)

	value, ok := c.Get("key")
	if !ok {
		t.Error("Empty value should be storable")
	}
	if value != "" {
		t.Errorf("Expected empty string, got '%v'", value)
	}
}

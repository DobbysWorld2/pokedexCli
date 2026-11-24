package pokecache

import (
	"bytes"
	"testing"
	"time"
)

func TestAddAndGet(t *testing.T) {
	cache := NewCache(time.Second)
	cache.Add("test", []byte("value"))

	val, ok := cache.Get("test")
	if !ok {
		t.Fatalf("expected key %q to exist", "test")
	}

	if !bytes.Equal(val, []byte("value")) {
		t.Errorf("expected %q, got %q", "value", string(val))
	}
}

func TestGetMissingKey(t *testing.T) {
	cache := NewCache(time.Second)

	val, ok := cache.Get("missing")
	if ok {
		t.Errorf("expected key %q not to exist", "missing")
	}

	if val != nil {
		t.Errorf("expected nil value for missing key, got %q", string(val))
	}
}

func TestAddReplacesExistingValue(t *testing.T) {
	cache := NewCache(time.Second)
	cache.Add("test", []byte("first"))
	cache.Add("test", []byte("second"))

	val, ok := cache.Get("test")
	if !ok {
		t.Fatalf("expected key %q to exist", "test")
	}

	if !bytes.Equal(val, []byte("second")) {
		t.Errorf("expected %q, got %q", "second", string(val))
	}
}

func TestDifferentKeysDoNotOverwrite(t *testing.T) {
	cache := NewCache(time.Second)
	cache.Add("first", []byte("one"))
	cache.Add("second", []byte("two"))

	first, ok := cache.Get("first")
	if !ok {
		t.Fatalf("expected key %q to exist", "first")
	}
	if !bytes.Equal(first, []byte("one")) {
		t.Errorf("expected %q, got %q", "one", string(first))
	}

	second, ok := cache.Get("second")
	if !ok {
		t.Fatalf("expected key %q to exist", "second")
	}
	if !bytes.Equal(second, []byte("two")) {
		t.Errorf("expected %q, got %q", "two", string(second))
	}
}

func TestNewCacheIsEmpty(t *testing.T) {
	cache := NewCache(time.Second)

	if len(cache.entries) != 0 {
		t.Errorf("expected new cache to be empty, got %d entries", len(cache.entries))
	}
}


func TestOldEntryIsEventuallyRemoved(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)
	cache.Add("old", []byte("value"))

	// The entry expires at 50ms, but the ticker only removes it on a tick,
	// so allow another tick to pass before asserting.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if _, ok := cache.Get("old"); !ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("expected key %q to be reaped", "old")
}

func TestFreshEntryStaysAvailable(t *testing.T) {
	cache := NewCache(500 * time.Millisecond)
	cache.Add("fresh", []byte("value"))

	// The entry is younger than the interval, so it must remain available.
	time.Sleep(100 * time.Millisecond)

	val, ok := cache.Get("fresh")
	if !ok {
		t.Fatalf("expected key %q to still exist", "fresh")
	}

	if !bytes.Equal(val, []byte("value")) {
		t.Errorf("expected %q, got %q", "value", string(val))
	}
}

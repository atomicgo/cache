package cache_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"atomicgo.dev/cache"
)

func TestNew(t *testing.T) {
	t.Run("Default options", func(t *testing.T) {
		t.Parallel()

		c := cache.New[string]()
		if c == nil {
			t.Fatal("expected cache to be initialized")
		}

		if c.Options.DefaultExpiration != 0 {
			t.Errorf("expected DefaultExpiration to be 0, got %v", c.Options.DefaultExpiration)
		}

		if c.Options.AutoPurgeInterval != 0 {
			t.Errorf("expected AutoPurgeInterval to be 0, got %v", c.Options.AutoPurgeInterval)
		}
	})

	t.Run("With options", func(t *testing.T) {
		t.Parallel()

		opts := cache.Options{
			DefaultExpiration: time.Second * 10,
			AutoPurgeInterval: time.Minute,
		}
		c := cache.New[string](opts)

		if c.Options != opts {
			t.Errorf("expected options %v, got %v", opts, c.Options)
		}
	})
}

func TestCache_SetGet(t *testing.T) {
	c := cache.New[string]()

	t.Run("Set and Get", func(t *testing.T) {
		c.Set("key1", "value1")
		c.Set("key2", "value2")

		t.Run("Get existing keys", func(t *testing.T) {
			t.Parallel()

			if val := c.Get("key1"); val != "value1" {
				t.Errorf("expected 'value1', got '%v'", val)
			}

			if val := c.Get("key2"); val != "value2" {
				t.Errorf("expected 'value2', got '%v'", val)
			}
		})

		t.Run("Get non-existing key", func(t *testing.T) {
			t.Parallel()

			if val := c.Get("key3"); val != "" {
				t.Errorf("expected '', got '%v'", val)
			}
		})
	})
}

func TestCache_Contains(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1")

	t.Run("Contains existing key", func(t *testing.T) {
		t.Parallel()

		if !c.Contains("key1") {
			t.Error("expected cache to contain 'key1'")
		}
	})

	t.Run("Does not contain non-existing key", func(t *testing.T) {
		t.Parallel()

		if c.Contains("key2") {
			t.Error("expected cache not to contain 'key2'")
		}
	})
}

func TestCache_Delete(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1")
	c.Delete("key1")

	t.Run("Deleted key should not be present", func(t *testing.T) {
		t.Parallel()

		if c.Contains("key1") {
			t.Error("expected 'key1' to be deleted")
		}
	})
}

func TestCache_Purge(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1")
	c.Set("key2", "value2")
	c.Purge()

	t.Run("Cache should be empty after purge", func(t *testing.T) {
		t.Parallel()

		if size := c.Size(); size != 0 {
			t.Errorf("expected size 0, got %d", size)
		}
	})
}

func TestCache_Size(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1")
	c.Set("key2", "value2")

	t.Run("Size should reflect number of items", func(t *testing.T) {
		t.Parallel()

		if size := c.Size(); size != 2 {
			t.Errorf("expected size 2, got %d", size)
		}
	})
}

func TestCache_Expiration(t *testing.T) {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Millisecond * 100,
	})

	t.Run("Set values with varying expiration", func(t *testing.T) {
		c.Set("key1", "value1")                       // Expires in 100ms (default)
		c.Set("key2", "value2", time.Millisecond*200) // Expires in 200ms
		c.Set("key3", "value3", time.Millisecond*300) // Expires in 300ms
		c.Set("key4", "value4", 0)                    // No expiration
	})

	time.Sleep(time.Millisecond * 150) // Wait 150ms

	t.Run("Check expiration after 150ms", func(t *testing.T) {
		if c.Contains("key1") {
			t.Error("expected 'key1' to be expired")
		}

		if !c.Contains("key2") {
			t.Error("expected 'key2' to be valid")
		}

		if !c.Contains("key3") {
			t.Error("expected 'key3' to be valid")
		}

		if !c.Contains("key4") {
			t.Error("expected 'key4' to be valid")
		}
	})

	time.Sleep(time.Millisecond * 100) // Wait additional 100ms (250ms total)

	t.Run("Check expiration after 250ms", func(t *testing.T) {
		if c.Contains("key2") {
			t.Error("expected 'key2' to be expired")
		}

		if !c.Contains("key3") {
			t.Error("expected 'key3' to be valid")
		}

		if !c.Contains("key4") {
			t.Error("expected 'key4' to be valid")
		}
	})

	time.Sleep(time.Millisecond * 100) // Wait additional 100ms (350ms total)

	t.Run("Check expiration after 350ms", func(t *testing.T) {
		if c.Contains("key3") {
			t.Error("expected 'key3' to be expired")
		}

		if !c.Contains("key4") {
			t.Error("expected 'key4' to be valid")
		}
	})
}

func TestCache_GetOrSet(t *testing.T) {
	c := cache.New[string]()

	t.Run("Key does not exist", func(t *testing.T) {
		t.Parallel()

		value := c.GetOrSet("key1", func() string { return "value1" })
		if value != "value1" {
			t.Errorf("expected 'value1', got '%v'", value)
		}
	})

	t.Run("Key exists", func(t *testing.T) {
		t.Parallel()

		value := c.GetOrSet("key1", func() string { return "new_value" })
		if value != "value1" {
			t.Errorf("expected 'value1', got '%v'", value)
		}
	})
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := cache.New[*int64]()
	var wg sync.WaitGroup
	key := "counter"
	var initialValue int64
	c.Set(key, &initialValue)

	// Concurrent increments
	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			c.GetOrSet(key, func() *int64 { return new(int64) }) // Ensure key exists
			ptr := c.Get(key)
			atomic.AddInt64(ptr, 1)
		}()
	}

	wg.Wait()

	t.Run("Concurrent increments result", func(t *testing.T) {
		val := atomic.LoadInt64(c.Get(key))
		if val != 1000 {
			t.Errorf("expected '1000', got '%v'", val)
		}
	})
}

func TestCache_AutoPurge(t *testing.T) {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Millisecond * 10,
	})
	c.Set("key1", "value1")
	c.EnableAutoPurge(time.Millisecond * 5)

	time.Sleep(time.Millisecond * 20) // Wait for auto purge to happen

	t.Run("Auto purge should remove expired keys", func(t *testing.T) {
		if c.Contains("key1") {
			t.Error("expected 'key1' to be auto purged")
		}

		if size := c.Size(); size != 0 {
			t.Errorf("expected size 0, got %d", size)
		}
	})

	c.StopAutoPurge()
}

func TestCache_KeysValues(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1")
	c.Set("key2", "value2")

	t.Run("Keys", func(t *testing.T) {
		t.Parallel()

		keys := c.Keys()
		if len(keys) != 2 {
			t.Errorf("expected 2 keys, got %d", len(keys))
		}

		expectedKeys := map[string]bool{"key1": true, "key2": true}
		for _, key := range keys {
			if !expectedKeys[key] {
				t.Errorf("unexpected key '%s'", key)
			}
		}
	})

	t.Run("Values", func(t *testing.T) {
		t.Parallel()

		values := c.Values()
		if len(values) != 2 {
			t.Errorf("expected 2 values, got %d", len(values))
		}

		expectedValues := map[string]bool{"value1": true, "value2": true}
		for _, value := range values {
			if !expectedValues[value] {
				t.Errorf("unexpected value '%s'", value)
			}
		}
	})
}

func TestCache_SetExpiration(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1")
	c.Set("key2", "value2")

	c.SetExpiration(time.Millisecond * 10)
	time.Sleep(time.Millisecond * 15)

	t.Run("All keys should be expired after SetExpiration", func(t *testing.T) {
		if c.Contains("key1") || c.Contains("key2") {
			t.Error("expected all keys to be expired")
		}
	})
}

func TestCache_Close(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1")
	c.EnableAutoPurge()
	c.Close()

	t.Run("Cache should be empty after Close", func(t *testing.T) {
		if size := c.Size(); size != 0 {
			t.Errorf("expected size 0, got %d", size)
		}
	})
}

func TestCache_ExpirationEdgeCases(t *testing.T) {
	c := cache.New[string]()

	t.Run("Zero expiration (no expiration)", func(t *testing.T) {
		c.Set("key1", "value1", 0)
		time.Sleep(time.Millisecond * 5)

		if !c.Contains("key1") {
			t.Error("expected 'key1' to not expire")
		}
	})

	t.Run("Negative expiration (immediate expiration)", func(t *testing.T) {
		c.Set("key2", "value2", -time.Second)

		if c.Contains("key2") {
			t.Error("expected 'key2' to be expired immediately")
		}
	})

	t.Run("Immediate expiration", func(t *testing.T) {
		c.Set("key3", "value3", time.Nanosecond)
		time.Sleep(time.Millisecond * 1)

		if c.Contains("key3") {
			t.Error("expected 'key3' to be expired")
		}
	})
}

func TestCache_GetEntry(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1", time.Minute)

	t.Run("GetEntry returns correct entry", func(t *testing.T) {
		t.Parallel()

		entry := c.GetEntry("key1")
		if entry == nil {
			t.Fatal("expected entry to not be nil")
		}

		if entry.Value != "value1" {
			t.Errorf("expected value 'value1', got '%v'", entry.Value)
		}

		if entry.Expiration != time.Minute {
			t.Errorf("expected expiration %v, got %v", time.Minute, entry.Expiration)
		}
	})

	t.Run("GetEntry on non-existing key", func(t *testing.T) {
		t.Parallel()

		if entry := c.GetEntry("key2"); entry != nil {
			t.Error("expected entry to be nil for non-existing key")
		}
	})
}

func TestCache_ValidKeys(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1", time.Millisecond*10)
	c.Set("key2", "value2", time.Millisecond*20)
	c.Set("key3", "value3") // No expiration

	time.Sleep(time.Millisecond * 15)

	t.Run("ValidKeys returns unexpired keys", func(t *testing.T) {
		t.Parallel()

		keys := c.ValidKeys()
		expectedKeys := map[string]bool{"key2": true, "key3": true}

		if len(keys) != 2 {
			t.Errorf("expected 2 valid keys, got %d", len(keys))
		}

		for _, key := range keys {
			if !expectedKeys[key] {
				t.Errorf("unexpected valid key '%s'", key)
			}
		}
	})
}

func TestCache_ValidSize(t *testing.T) {
	c := cache.New[string]()
	c.Set("key1", "value1", time.Millisecond*10)
	c.Set("key2", "value2", time.Millisecond*20)
	c.Set("key3", "value3") // No expiration

	time.Sleep(time.Millisecond * 15)

	t.Run("ValidSize returns correct count of unexpired items", func(t *testing.T) {
		t.Parallel()

		if size := c.ValidSize(); size != 2 {
			t.Errorf("expected valid size 2, got %d", size)
		}
	})
}

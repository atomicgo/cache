package cache_test

import (
	"testing"
	"testing/synctest"
	"time"

	"atomicgo.dev/cache"
)

func TestNew(_ *testing.T) {
	cache.New[string]()
}

func TestCache(t *testing.T) {
	c := cache.New[string]()

	t.Run("Get", func(_ *testing.T) {
		c.Set("1", "one")
		c.Set("2", "two")
		c.Set("3", "three")
	})

	t.Run("Contains", func(t *testing.T) {
		if !c.Contains("1") {
			t.Error("expected to contain key 1")
		}

		if !c.Contains("2") {
			t.Error("expected to contain key 2")
		}

		if !c.Contains("3") {
			t.Error("expected to contain key 3")
		}
	})

	t.Run("Get", func(t *testing.T) {
		if c.Get("1") != "one" {
			t.Error("expected key '1' to be \"one\"")
		}

		if c.Get("2") != "two" {
			t.Error("expected key '2' to be \"two\"")
		}

		if c.Get("3") != "three" {
			t.Error("expected key '3' to be \"three\"")
		}
	})

	t.Run("Delete", func(_ *testing.T) {
		c.Delete("1")
	})

	t.Run("Does not contain deleted key", func(t *testing.T) {
		if c.Contains("1") {
			t.Error("expected to not contain key 1")
		}

		if !c.Contains("2") {
			t.Error("expected to contain key 2")
		}

		if !c.Contains("3") {
			t.Error("expected to contain key 3")
		}
	})

	t.Run("Purge", func(_ *testing.T) {
		c.Purge()
	})

	t.Run("Does not contain any keys", func(t *testing.T) {
		if c.Contains("1") {
			t.Error("expected to not contain key 1")
		}

		if c.Contains("2") {
			t.Error("expected to not contain key 2")
		}

		if c.Contains("3") {
			t.Error("expected to not contain key 3")
		}
	})

	t.Run("Size is 0", func(t *testing.T) {
		if c.Size() != 0 {
			t.Error("expected length to be 0")
		}
	})
}

func TestCache_Expiration(t *testing.T) {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second,
	})
	c.Set("expired", "old")
	expireEntry(t, c, "expired")

	if !c.Expired("expired") {
		t.Error("expected key to be expired")
	}

	if c.Contains("expired") {
		t.Error("expected expired key to not be contained")
	}

	if got := c.Get("expired"); got != "" {
		t.Errorf("expected expired key to return zero value, got %q", got)
	}

	if c.Size() != 0 {
		t.Error("expected Get to delete expired key")
	}
}

func TestCache_MissingKeyAccessors(t *testing.T) {
	c := cache.New[string]()

	if c.Expired("missing") {
		t.Error("expected missing key to not be expired")
	}

	if got := c.GetExpiration("missing"); got != 0 {
		t.Errorf("expected missing key expiration to be 0, got %s", got)
	}
}

func TestCache_GetOrSet(t *testing.T) {
	t.Run("Not existing", func(t *testing.T) {
		c := cache.New[string]()
		res := c.GetOrSet("1", func() string {
			return "one"
		})

		if res != "one" {
			t.Error("expected to return \"one\"")
		}
	})

	t.Run("Existing", func(t *testing.T) {
		c := cache.New[string]()
		c.Set("1", "one")
		res := c.GetOrSet("1", func() string {
			return "asd"
		})

		if res != "one" {
			t.Error("expected to return \"one\"")
		}
	})

	t.Run("Expired", func(t *testing.T) {
		c := cache.New[string]()
		c.Set("1", "one", time.Second)
		expireEntry(t, c, "1")

		calls := 0
		res := c.GetOrSet("1", func() string {
			calls++

			return "two"
		}, time.Minute)

		if res != "two" {
			t.Error("expected to return \"two\"")
		}

		if calls != 1 {
			t.Errorf("expected callback to be called once, got %d", calls)
		}

		if c.Get("1") != "two" {
			t.Error("expected key '1' to be \"two\"")
		}

		if c.GetExpiration("1") != time.Minute {
			t.Error("expected replacement value to use new expiration")
		}
	})
}

func TestCacheExpirationWithSynctest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := cache.New[string](cache.Options{
			DefaultExpiration: 10 * time.Millisecond,
		})
		c.Set("1", "one")
		c.Set("2", "two", 20*time.Millisecond)
		c.Set("3", "three", 30*time.Millisecond)

		requireContains(t, c, "1", true)
		requireContains(t, c, "2", true)
		requireContains(t, c, "3", true)

		time.Sleep(10*time.Millisecond - time.Nanosecond)
		requireContains(t, c, "1", true)
		requireContains(t, c, "2", true)
		requireContains(t, c, "3", true)

		time.Sleep(time.Nanosecond)
		requireContains(t, c, "1", false)
		requireContains(t, c, "2", true)
		requireContains(t, c, "3", true)

		time.Sleep(10 * time.Millisecond)
		requireContains(t, c, "1", false)
		requireContains(t, c, "2", false)
		requireContains(t, c, "3", true)

		time.Sleep(10 * time.Millisecond)
		requireContains(t, c, "1", false)
		requireContains(t, c, "2", false)
		requireContains(t, c, "3", false)
	})
}

func TestCacheAutoPurgeWithSynctest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c := cache.New[string](cache.Options{
			DefaultExpiration: 10 * time.Millisecond,
			AutoPurgeInterval: 5 * time.Millisecond,
		})
		c.EnableAutoPurge()
		c.Set("short", "short")
		c.Set("long", "long", 20*time.Millisecond)

		time.Sleep(10 * time.Millisecond)
		synctest.Wait()

		if got := c.Size(); got != 1 {
			t.Fatalf("expected auto purge to remove one key, got size %d", got)
		}

		requireContains(t, c, "long", true)

		c.StopAutoPurge()
		c.StopAutoPurge()
		c.Set("stopped", "stopped", 5*time.Millisecond)

		time.Sleep(10 * time.Millisecond)
		synctest.Wait()

		if got := c.Size(); got != 2 {
			t.Fatalf("expected stopped auto purge to keep expired keys, got size %d", got)
		}

		requireContains(t, c, "stopped", false)

		c.EnableAutoPurge(5 * time.Millisecond)
		time.Sleep(5 * time.Millisecond)
		synctest.Wait()

		if got := c.Size(); got != 0 {
			t.Fatalf("expected restarted auto purge to remove expired keys, got size %d", got)
		}

		c.Close()
	})
}

func requireContains(t *testing.T, c *cache.Cache[string], key string, want bool) {
	t.Helper()

	if got := c.Contains(key); got != want {
		t.Fatalf("Contains(%q) = %t, want %t", key, got, want)
	}
}

func expireEntry(t *testing.T, c *cache.Cache[string], key string) {
	t.Helper()

	entry := c.GetEntry(key)
	if entry == nil {
		t.Fatalf("expected key %q to exist", key)
	}

	entry.CachedAt = entry.CachedAt.Add(-entry.Expiration)
}

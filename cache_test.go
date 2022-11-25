package cache_test

import (
	"atomicgo.dev/cache"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cache.New[string]()
}

func TestCache(t *testing.T) {
	c := cache.New[string]()

	t.Run("Get", func(t *testing.T) {
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

	t.Run("Delete", func(t *testing.T) {
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

	t.Run("Purge", func(t *testing.T) {
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
		DefaultExpiration: time.Millisecond * 10,
	})

	t.Run("Set values", func(t *testing.T) {
		c.Set("1", "one")                        // should use default of 10ms
		c.Set("2", "two", time.Millisecond*20)   // should use 20ms
		c.Set("3", "three", time.Millisecond*30) // should use 30ms
	})

	t.Run("No entries should be expired", func(t *testing.T) {
		time.Sleep(time.Millisecond * 5) // 5ms total

		// Expect all to be present
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

	t.Run("First entry should be expired", func(t *testing.T) {
		time.Sleep(time.Millisecond * 5) // 10ms total

		// Expect 1 to be gone, 2 and 3 to be present
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

	t.Run("First and second entry should be expired", func(t *testing.T) {
		time.Sleep(time.Millisecond * 10) // 20ms total

		// Expect 1 and 2 to be gone, 3 to be present
		if c.Contains("1") {
			t.Error("expected to not contain key 1")
		}
		if c.Contains("2") {
			t.Error("expected to not contain key 2")
		}
		if !c.Contains("3") {
			t.Error("expected to contain key 3")
		}
	})

	t.Run("First, second and third entry should be expired", func(t *testing.T) {
		time.Sleep(time.Millisecond * 10) // 30ms total

		// Expect 1, 2, and 3 to be gone
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
}

package cache_test

import (
	"atomicgo.dev/cache"
	"fmt"
	"sort"
	"time"
)

func ExampleNew() {
	// Create a cache for string values, without any options.
	cache.New[string]()

	// Create a cache for string values, with options.
	cache.New[string](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Create a cache for int values, without any options.
	cache.New[int]()

	// Create a cache for int values, with options.
	cache.New[int](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Create a cache for any values, without any options.
	cache.New[any]()

	// Create a cache for any values, with options.
	cache.New[any](cache.Options{
		DefaultExpiration: time.Second * 10,
	})
}

func ExampleCache_Set() {
	c := cache.New[string]()

	// Set a value for a key.
	c.Set("1", "one")

	// Set a value for a key with a custom expiration.
	c.Set("2", "two", time.Second*10)
}

func ExampleCache_Get() {
	c := cache.New[string]()

	// Set a value for a key.
	c.Set("1", "one")

	// Get the value for a key.
	fmt.Println(c.Get("1"))

	// Output: one
}

func ExampleCache_ValidKeys() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")
	c.Set("expiresSoon", "data", time.Millisecond*10)

	validKeys := c.ValidKeys()
	sort.Strings(validKeys)

	fmt.Println(validKeys)

	// Sleep for 10ms to let the "expiresSoon" key expire.
	time.Sleep(time.Millisecond * 10)

	// Get the valid keys.
	validKeys = c.ValidKeys()
	sort.Strings(validKeys)
	fmt.Println(validKeys)

	// Output: [1 2 3 expiresSoon]
	// [1 2 3]
}

func ExampleCache_ValidSize() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")
	c.Set("expiresSoon", "data", time.Millisecond*10)

	fmt.Println(c.ValidSize())

	// Sleep for 10ms to let the "expiresSoon" key expire.
	time.Sleep(time.Millisecond * 10)

	// Get the valid size.
	fmt.Println(c.ValidSize())

	// Output: 4
	// 3
}

func ExampleCache_Contains() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Check if the cache contains a key.
	fmt.Println(c.Contains("1"))
	fmt.Println(c.Contains("4"))

	// Output: true
	// false
}

func ExampleCache_Purge() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Purge the cache.
	c.Purge()

	// Check if the cache is empty.
	fmt.Println(c.Size())

	// Output: 0
}

func ExampleCache_Keys() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Get the keys.
	keys := c.Keys()
	sort.Strings(keys)

	fmt.Println(keys)

	// Output: [1 2 3]
}

func ExampleCache_Size() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Get the size.
	fmt.Println(c.Size())

	// Output: 3
}

func ExampleCache_GetExpiration() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")

	// Get the expiration for a key.
	fmt.Println(c.GetExpiration("1"))

	// Output: 10s
}

func ExampleCache_GetEntry() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")

	// Get the entry for a key.
	entry := c.GetEntry("1")
	fmt.Println(entry.Value)
	fmt.Println(entry.Expiration)

	// Output: one
	// 10s
}

func ExampleCache_SetExpiration() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three", time.Second*1337)

	// Set the expiration for a key.
	c.SetExpiration(time.Second * 20)

	// Get the expiration for a key.
	fmt.Println(c.GetExpiration("1"))
	fmt.Println(c.GetExpiration("2"))
	fmt.Println(c.GetExpiration("3"))

	// Output: 20s
	// 20s
	// 20s
}

func ExampleCache_Expired() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Millisecond * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")

	// Check if the key is expired.
	fmt.Println(c.Expired("1"))

	// Sleep for 10ms to let the key expire.
	time.Sleep(time.Millisecond * 10)

	// Check if the key is expired.
	fmt.Println(c.Expired("1"))

	// Output: false
	// true
}

func ExampleCache_PurgeExpired() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second,
	})

	// Fill the cache with some values.
	c.Set("1", "one", time.Millisecond*10)
	c.Set("2", "two", time.Millisecond*10)
	c.Set("3", "three")

	// Purge the expired keys.
	c.PurgeExpired()
	fmt.Println(c.Size())

	// Sleep for 10ms to let the first two keys expire.
	time.Sleep(time.Millisecond * 10)

	// Purge the expired keys.
	c.PurgeExpired()
	fmt.Println(c.Size())

	// Output: 3
	// 1
}

func ExampleCache_GetOrSet() {
	c := cache.New[string]()

	// Try to get or set a value for a key.
	c.GetOrSet("1", func() string {
		return "one"
	})

	fmt.Println(c.Get("1"))

	// try to get or set a value for an existing key.
	c.GetOrSet("1", func() string {
		return "something else"
	})
	fmt.Println(c.Get("1"))

	// delete the key
	c.Delete("1")

	// try to get or set a value for a non-existing key.
	c.GetOrSet("1", func() string {
		return "something else"
	})
	fmt.Println(c.Get("1"))

	// Output: one
	// one
	// something else
}

func ExampleCache_Values() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Get the values.
	values := c.Values()
	sort.Strings(values)

	fmt.Println(values)

	// Output: [one three two]
}

func ExampleCache_StopAutoPurge() {
	c := cache.New[string](cache.Options{
		AutoPurgeInterval: time.Millisecond * 10,
	})
	c.EnableAutoPurge()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Stop the auto purge.
	c.StopAutoPurge()

	// Sleep for 10 milliseconds to let the auto purge run.
	time.Sleep(time.Millisecond * 10)

	// Get the size.
	fmt.Println(c.Size())

	// Output: 3
}

func ExampleCache_Delete() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Delete a key.
	c.Delete("3")

	// Get the size.
	fmt.Println(c.Size())

	// Output: 2
}

func ExampleCache_Close() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Close the cache.
	c.Close()

	// Get the size.
	fmt.Println(c.Size())

	// Output: 0
}

func ExampleEntry_Expired() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Millisecond * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")

	// Get the entry for a key.
	entry := c.GetEntry("1")

	// Check if the entry is expired.
	fmt.Println(entry.Expired())

	// Sleep for 10ms to let the entry expire.
	time.Sleep(time.Millisecond * 10)

	// Check if the entry is expired.
	fmt.Println(entry.Expired())

	// Output: false
	// true
}

func ExampleEntry_SetExpiration() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")

	// Get the entry for a key.
	entry := c.GetEntry("1")

	// Set the expiration for the entry.
	entry.SetExpiration(time.Second * 20)

	// Get the expiration for a key.
	fmt.Println(c.GetExpiration("1"))

	// Output: 20s
}

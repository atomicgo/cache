package cache

import (
	"sync"
	"time"
)

// Cache is a fast and thread-safe cache implementation with expiration.
type Cache[T any] struct {
	mutex           sync.Mutex
	cache           map[string]*Entry[T]
	autoPurgeActive bool
	cancelAutoPurge chan struct{}
	Options         Options
}

// Options can be passed to New to configure the cache.
type Options struct {
	DefaultExpiration time.Duration
	AutoPurgeInterval time.Duration
}

// New returns a new Cache.
// The default Expiration time is 0, which means no Expiration.
func New[T any](options ...Options) *Cache[T] {
	opts := Options{}
	if len(options) > 0 {
		opts = options[0]
	}
	return &Cache[T]{
		cache:           make(map[string]*Entry[T]),
		cancelAutoPurge: make(chan struct{}),
		Options:         opts,
	}
}

// EnableAutoPurge starts a goroutine that purges expired keys from the cache.
// The interval is the time between purges.
// If the interval is 0, the default interval of the cache options is used.
// If the cache options do not specify a default interval, the default interval is 1 minute.
func (c *Cache[T]) EnableAutoPurge(purgeInterval ...time.Duration) *Cache[T] {
	if c.autoPurgeActive {
		return c
	}

	interval := c.Options.AutoPurgeInterval
	if len(purgeInterval) > 0 {
		interval = purgeInterval[0]
	}

	if interval == 0 {
		interval = time.Minute
	}

	c.autoPurgeActive = true

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				c.PurgeExpired()
			case <-c.cancelAutoPurge:
				ticker.Stop()
				return
			}
		}

	}()
	return c
}

// Set sets the value for a key.
// If the key already exists, the value is overwritten.
// If the Expiration time is 0, the default Expiration time is used.
func (c *Cache[T]) Set(key string, value T, expiration ...time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	exp := c.Options.DefaultExpiration
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	c.cache[key] = &Entry[T]{
		Value:      value,
		CachedAt:   time.Now(),
		Expiration: exp,
	}
}

// Get returns the value for a key.
// If the key does not exist, nil is returned.
// If the key is expired, the zero value is returned and the key is deleted.
func (c *Cache[T]) Get(key string) T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v, ok := c.cache[key]
	if !ok || v.Expired() {
		return *new(T)
	}
	return v.Value
}

// GetEntry returns the Entry for a key.
// If the key does not exist, nil is returned.
func (c *Cache[T]) GetEntry(key string) *Entry[T] {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.cache[key]
}

// GetOrSet returns the value for a key.
// If the key does not exist, the value is set to the result of the callback function and returned.
// If the key is expired, the value is set to the result of the callback function and returned.
func (c *Cache[T]) GetOrSet(key string, callback func() T, expiration ...time.Duration) T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v, ok := c.cache[key]
	if !ok || v.Expired() {
		exp := c.Options.DefaultExpiration
		if len(expiration) > 0 {
			exp = expiration[0]
		}
		c.cache[key] = &Entry[T]{
			Value:      callback(),
			CachedAt:   time.Now(),
			Expiration: exp,
		}
		return c.cache[key].Value
	}
	return v.Value
}

// Expired returns true if the key is expired.
func (c *Cache[T]) Expired(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.cache[key].Expired()
}

// Contains returns true if the key is in the cache, and the key is not expired.
func (c *Cache[T]) Contains(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.cache[key]
	return ok && !c.cache[key].Expired()
}

// SetExpiration sets the Expiration time all keys in the cache.
func (c *Cache[T]) SetExpiration(expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, v := range c.cache {
		v.SetExpiration(expiration)
	}
}

// GetExpiration returns the Expiration time for a key.
func (c *Cache[T]) GetExpiration(key string) time.Duration {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.cache[key].Expiration
}

// Delete removes the key from the cache.
func (c *Cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)
}

// Purge removes all keys from the cache.
func (c *Cache[T]) Purge() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*Entry[T])
}

// PurgeExpired removes all expired keys from the cache.
func (c *Cache[T]) PurgeExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for k, v := range c.cache {
		if v.Expired() {
			delete(c.cache, k)
		}
	}
}

// Size returns the number of keys in the cache, expired or not.
func (c *Cache[T]) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.cache)
}

// ValidSize returns the number of keys in the cache that are not expired.
func (c *Cache[T]) ValidSize() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	size := 0
	for _, v := range c.cache {
		if !v.Expired() {
			size++
		}
	}
	return size
}

// Keys returns a slice of all keys in the cache, expired or not.
func (c *Cache[T]) Keys() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	keys := make([]string, len(c.cache))
	i := 0
	for k := range c.cache {
		keys[i] = k
		i++
	}
	return keys
}

// ValidKeys returns a slice of all valid keys in the cache.
func (c *Cache[T]) ValidKeys() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	keys := make([]string, 0)
	for k, v := range c.cache {
		if !v.Expired() {
			keys = append(keys, k)
		}
	}
	return keys
}

// Values returns a slice of all values in the cache.
func (c *Cache[T]) Values() []T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	values := make([]T, len(c.cache))
	i := 0
	for _, v := range c.cache {
		values[i] = v.Value
		i++
	}
	return values
}

// StopAutoPurge stops the auto purge goroutine.
func (c *Cache[T]) StopAutoPurge() {
	if c.autoPurgeActive {
		c.cancelAutoPurge <- struct{}{}
	}
}

// Close purges the cache and stops the auto purge goroutine, if active.
func (c *Cache[T]) Close() {
	if c.autoPurgeActive {
		c.StopAutoPurge()
	}
	c.Purge()
}

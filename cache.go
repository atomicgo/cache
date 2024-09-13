package cache

import (
	"sync"
	"time"
)

// Cache is a fast and thread-safe cache implementation with expiration for any type T.
type Cache[T any] struct {
	cache           map[string]*Entry[T]
	cancelAutoPurge chan struct{}
	Options         Options
	mutex           sync.RWMutex
	autoPurgeActive bool
}

// Options can be passed to New to configure the cache.
type Options struct {
	// DefaultExpiration is the default duration before an entry expires.
	DefaultExpiration time.Duration
	// AutoPurgeInterval is the interval between automatic purging of expired entries.
	AutoPurgeInterval time.Duration
}

// New creates a new Cache with optional Options.
// The default expiration time is 0, which means no expiration.
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

// EnableAutoPurge starts a goroutine that periodically purges expired entries from the cache.
// The interval between purges can be specified; if not provided, the cache's AutoPurgeInterval option is used,
// or defaults to 1 minute if not set.
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
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.PurgeExpired()
			case <-c.cancelAutoPurge:
				return
			}
		}
	}()

	return c
}

// Set sets the value for a key with an optional expiration.
// If the key already exists, the value is overwritten.
// If the expiration time is 0, the cache's DefaultExpiration is used.
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

// Get returns the value associated with the key.
// If the key does not exist or is expired, the zero value of T is returned.
func (c *Cache[T]) Get(key string) T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	val, ok := c.cache[key]
	if !ok || val.Expired() {
		var zero T
		return zero
	}

	return val.Value
}

// GetEntry returns the cache Entry associated with the key.
// If the key does not exist, nil is returned.
func (c *Cache[T]) GetEntry(key string) *Entry[T] {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	val, ok := c.cache[key]
	if !ok {
		return nil
	}

	return val
}

// GetOrSet returns the value associated with the key.
// If the key does not exist or is expired, the value is set to the result of the callback function and returned.
// The callback function is called without holding a lock.
func (c *Cache[T]) GetOrSet(key string, callback func() T, expiration ...time.Duration) T {
	// First, try to get the value with a read lock
	c.mutex.RLock()

	val, ok := c.cache[key]
	if ok && !val.Expired() {
		c.mutex.RUnlock()
		return val.Value
	}
	c.mutex.RUnlock()

	// Not found or expired, compute value without holding the lock
	value := callback()

	// Now, acquire write lock to set the value
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check again if the key was set by another goroutine
	val, ok = c.cache[key]
	if ok && !val.Expired() {
		return val.Value
	}

	// Set the new value
	exp := c.Options.DefaultExpiration
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	c.cache[key] = &Entry[T]{
		Value:      value,
		CachedAt:   time.Now(),
		Expiration: exp,
	}

	return value
}

// Expired checks if the key is expired.
func (c *Cache[T]) Expired(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	v, ok := c.cache[key]
	if !ok {
		return false
	}

	return v.Expired()
}

// Contains returns true if the key exists in the cache and is not expired.
func (c *Cache[T]) Contains(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	v, ok := c.cache[key]

	return ok && !v.Expired()
}

// SetExpiration sets the expiration time for all entries in the cache.
func (c *Cache[T]) SetExpiration(expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, v := range c.cache {
		v.SetExpiration(expiration)
	}
}

// GetExpiration returns the expiration time for a specific key.
// If the key does not exist, zero duration is returned.
func (c *Cache[T]) GetExpiration(key string) time.Duration {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	v, ok := c.cache[key]
	if !ok {
		return 0
	}

	return v.Expiration
}

// Delete removes the key from the cache.
func (c *Cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.cache, key)
}

// Purge removes all entries from the cache.
func (c *Cache[T]) Purge() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*Entry[T])
}

// PurgeExpired removes all expired entries from the cache.
func (c *Cache[T]) PurgeExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for k, v := range c.cache {
		if v.ExpiredAt(now) {
			delete(c.cache, k)
		}
	}
}

// Size returns the total number of entries in the cache, including expired entries.
func (c *Cache[T]) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.cache)
}

// ValidSize returns the number of valid (not expired) entries in the cache.
func (c *Cache[T]) ValidSize() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	size := 0

	now := time.Now()
	for _, v := range c.cache {
		if !v.ExpiredAt(now) {
			size++
		}
	}

	return size
}

// Keys returns a slice of all keys in the cache, including expired keys.
func (c *Cache[T]) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.cache))
	for k := range c.cache {
		keys = append(keys, k)
	}

	return keys
}

// ValidKeys returns a slice of all valid (not expired) keys in the cache.
func (c *Cache[T]) ValidKeys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]string, 0, len(c.cache))

	now := time.Now()
	for k, v := range c.cache {
		if !v.ExpiredAt(now) {
			keys = append(keys, k)
		}
	}

	return keys
}

// Values returns a slice of all values in the cache, including expired values.
func (c *Cache[T]) Values() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	values := make([]T, 0, len(c.cache))
	for _, v := range c.cache {
		values = append(values, v.Value)
	}

	return values
}

// StopAutoPurge stops the auto purge goroutine.
func (c *Cache[T]) StopAutoPurge() {
	if c.autoPurgeActive {
		c.cancelAutoPurge <- struct{}{}
		c.autoPurgeActive = false
	}
}

// Close purges the cache and stops the auto purge goroutine, if active.
func (c *Cache[T]) Close() {
	c.StopAutoPurge()
	c.Purge()
}

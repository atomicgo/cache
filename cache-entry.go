package cache

import "time"

// Entry is a cache entry.
type Entry[T any] struct {
	Value      T
	CachedAt   time.Time
	Expiration time.Duration
}

// Expired returns if the Entry is expired.
func (e Entry[T]) Expired() bool {
	if e.Expiration == 0 {
		return false
	}
	return time.Now().After(e.CachedAt.Add(e.Expiration))
}

// SetExpiration sets the Expiration time for the Entry.
func (e *Entry[T]) SetExpiration(expiration time.Duration) {
	e.Expiration = expiration
}

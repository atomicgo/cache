package cache

import "time"

// Entry is a cache entry storing a value of type T.
type Entry[T any] struct {
	Value      T
	CachedAt   time.Time
	Expiration time.Duration
}

// Expired checks if the entry is expired.
func (e *Entry[T]) Expired() bool {
	if e.Expiration == 0 {
		return false
	}

	return time.Now().After(e.CachedAt.Add(e.Expiration))
}

// ExpiredAt checks if the entry is expired at the given time.
func (e *Entry[T]) ExpiredAt(t time.Time) bool {
	if e.Expiration == 0 {
		return false
	}

	return t.After(e.CachedAt.Add(e.Expiration))
}

// SetExpiration sets the expiration duration for the entry.
func (e *Entry[T]) SetExpiration(expiration time.Duration) {
	e.Expiration = expiration
}

package cache_test

import (
	"atomicgo.dev/cache"
	"testing"
)

func BenchmarkCache_Set_EmptyStruct(b *testing.B) {
	c := cache.New[struct{}]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("1", struct{}{})
	}
}

func BenchmarkCache_Set_Int(b *testing.B) {
	c := cache.New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("1", 1)
	}
}

func BenchmarkCache_Set_String(b *testing.B) {
	c := cache.New[string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("1", "one")
	}
}

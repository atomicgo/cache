<h1 align="center">AtomicGo | cache</h1>

<p align="center">
<img src="https://img.shields.io/endpoint?url=https://atomicgo.dev/api/shields/cache&style=flat-square" alt="Downloads">

<a href="https://github.com/atomicgo/cache/releases">
<img src="https://img.shields.io/github/v/release/atomicgo/cache?style=flat-square" alt="Latest Release">
</a>

<a href="https://codecov.io/gh/atomicgo/cache" target="_blank">
<img src="https://img.shields.io/github/workflow/status/atomicgo/cache/Go?label=tests&style=flat-square" alt="Tests">
</a>

<a href="https://codecov.io/gh/atomicgo/cache" target="_blank">
<img src="https://img.shields.io/codecov/c/gh/atomicgo/cache?color=magenta&logo=codecov&style=flat-square" alt="Coverage">
</a>

<a href="https://codecov.io/gh/atomicgo/cache">
<!-- unittestcount:start --><img src="https://img.shields.io/badge/Unit_Tests-38-magenta?style=flat-square" alt="Unit test count"><!-- unittestcount:end -->
</a>

<a href="https://opensource.org/licenses/MIT" target="_blank">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>
  
<a href="https://goreportcard.com/report/github.com/atomicgo/cache" target="_blank">
<img src="https://goreportcard.com/badge/github.com/atomicgo/cache?style=flat-square" alt="Go report">
</a>   

</p>

---

<p align="center">
<strong><a href="https://pkg.go.dev/atomicgo.dev/cache#section-documentation" target="_blank">Documentation</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CONTRIBUTING.md" target="_blank">Contributing</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CODE_OF_CONDUCT.md" target="_blank">Code of Conduct</a></strong>
</p>

---

<p align="center">
  <img src="https://raw.githubusercontent.com/atomicgo/atomicgo/main/assets/header.png" alt="AtomicGo">
</p>

<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get atomicgo.dev/cache</pre></h3>
<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>

<!-- gomarkdoc:embed:start -->

<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# cache

```go
import "atomicgo.dev/cache"
```

Package cache is a generic, fast and thread\-safe cache implementation to improve performance of your Go applications.

## Index

- [type Cache](<#type-cache>)
  - [func New[T any](options ...Options) *Cache[T]](<#func-new>)
  - [func (c *Cache[T]) Close()](<#func-cachet-close>)
  - [func (c *Cache[T]) Contains(key string) bool](<#func-cachet-contains>)
  - [func (c *Cache[T]) Delete(key string)](<#func-cachet-delete>)
  - [func (c *Cache[T]) EnableAutoPurge(purgeInterval ...time.Duration) *Cache[T]](<#func-cachet-enableautopurge>)
  - [func (c *Cache[T]) Expired(key string) bool](<#func-cachet-expired>)
  - [func (c *Cache[T]) Get(key string) T](<#func-cachet-get>)
  - [func (c *Cache[T]) GetEntry(key string) *Entry[T]](<#func-cachet-getentry>)
  - [func (c *Cache[T]) GetExpiration(key string) time.Duration](<#func-cachet-getexpiration>)
  - [func (c *Cache[T]) GetOrSet(key string, callback func() T, expiration ...time.Duration) T](<#func-cachet-getorset>)
  - [func (c *Cache[T]) Keys() []string](<#func-cachet-keys>)
  - [func (c *Cache[T]) Purge()](<#func-cachet-purge>)
  - [func (c *Cache[T]) PurgeExpired()](<#func-cachet-purgeexpired>)
  - [func (c *Cache[T]) Set(key string, value T, expiration ...time.Duration)](<#func-cachet-set>)
  - [func (c *Cache[T]) SetExpiration(expiration time.Duration)](<#func-cachet-setexpiration>)
  - [func (c *Cache[T]) Size() int](<#func-cachet-size>)
  - [func (c *Cache[T]) StopAutoPurge()](<#func-cachet-stopautopurge>)
  - [func (c *Cache[T]) ValidKeys() []string](<#func-cachet-validkeys>)
  - [func (c *Cache[T]) ValidSize() int](<#func-cachet-validsize>)
  - [func (c *Cache[T]) Values() []T](<#func-cachet-values>)
- [type Entry](<#type-entry>)
  - [func (e Entry[T]) Expired() bool](<#func-entryt-expired>)
  - [func (e *Entry[T]) SetExpiration(expiration time.Duration)](<#func-entryt-setexpiration>)
- [type Options](<#type-options>)


## type [Cache](<https://github.com/atomicgo/cache/blob/main/cache.go#L9-L15>)

Cache is a fast and thread\-safe cache implementation with expiration.

```go
type Cache[T any] struct {
    Options Options
    // contains filtered or unexported fields
}
```

### func [New](<https://github.com/atomicgo/cache/blob/main/cache.go#L25>)

```go
func New[T any](options ...Options) *Cache[T]
```

New returns a new Cache. The default Expiration time is 0, which means no Expiration.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"time"
)

func main() {
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
```

</p>
</details>

### func \(\*Cache\[T\]\) [Close](<https://github.com/atomicgo/cache/blob/main/cache.go#L272>)

```go
func (c *Cache[T]) Close()
```

Close purges the cache and stops the auto purge goroutine, if active.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
)

func main() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Close the cache.
	c.Close()

	// Get the size.
	fmt.Println(c.Size())

}
```

#### Output

```
0
```

</p>
</details>

### func \(\*Cache\[T\]\) [Contains](<https://github.com/atomicgo/cache/blob/main/cache.go#L146>)

```go
func (c *Cache[T]) Contains(key string) bool
```

Contains returns true if the key is in the cache, and the key is not expired.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
)

func main() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Check if the cache contains a key.
	fmt.Println(c.Contains("1"))
	fmt.Println(c.Contains("4"))

}
```

#### Output

```
true
false
```

</p>
</details>

### func \(\*Cache\[T\]\) [Delete](<https://github.com/atomicgo/cache/blob/main/cache.go#L173>)

```go
func (c *Cache[T]) Delete(key string)
```

Delete removes the key from the cache.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
)

func main() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Delete a key.
	c.Delete("3")

	// Get the size.
	fmt.Println(c.Size())

}
```

#### Output

```
2
```

</p>
</details>

### func \(\*Cache\[T\]\) [EnableAutoPurge](<https://github.com/atomicgo/cache/blob/main/cache.go#L41>)

```go
func (c *Cache[T]) EnableAutoPurge(purgeInterval ...time.Duration) *Cache[T]
```

EnableAutoPurge starts a goroutine that purges expired keys from the cache. The interval is the time between purges. If the interval is 0, the default interval of the cache options is used. If the cache options do not specify a default interval, the default interval is 1 minute.

### func \(\*Cache\[T\]\) [Expired](<https://github.com/atomicgo/cache/blob/main/cache.go#L138>)

```go
func (c *Cache[T]) Expired(key string) bool
```

Expired returns true if the key is expired.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
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

}
```

#### Output

```
false
true
```

</p>
</details>

### func \(\*Cache\[T\]\) [Get](<https://github.com/atomicgo/cache/blob/main/cache.go#L94>)

```go
func (c *Cache[T]) Get(key string) T
```

Get returns the value for a key. If the key does not exist, nil is returned. If the key is expired, the zero value is returned and the key is deleted.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
)

func main() {
	c := cache.New[string]()

	// Set a value for a key.
	c.Set("1", "one")

	// Get the value for a key.
	fmt.Println(c.Get("1"))

}
```

#### Output

```
one
```

</p>
</details>

### func \(\*Cache\[T\]\) [GetEntry](<https://github.com/atomicgo/cache/blob/main/cache.go#L107>)

```go
func (c *Cache[T]) GetEntry(key string) *Entry[T]
```

GetEntry returns the Entry for a key. If the key does not exist, nil is returned.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")

	// Get the entry for a key.
	entry := c.GetEntry("1")
	fmt.Println(entry.Value)
	fmt.Println(entry.Expiration)

}
```

#### Output

```
one
10s
```

</p>
</details>

### func \(\*Cache\[T\]\) [GetExpiration](<https://github.com/atomicgo/cache/blob/main/cache.go#L165>)

```go
func (c *Cache[T]) GetExpiration(key string) time.Duration
```

GetExpiration returns the Expiration time for a key.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
	c := cache.New[string](cache.Options{
		DefaultExpiration: time.Second * 10,
	})

	// Set a value for a key.
	c.Set("1", "one")

	// Get the expiration for a key.
	fmt.Println(c.GetExpiration("1"))

}
```

#### Output

```
10s
```

</p>
</details>

### func \(\*Cache\[T\]\) [GetOrSet](<https://github.com/atomicgo/cache/blob/main/cache.go#L117>)

```go
func (c *Cache[T]) GetOrSet(key string, callback func() T, expiration ...time.Duration) T
```

GetOrSet returns the value for a key. If the key does not exist, the value is set to the result of the callback function and returned. If the key is expired, the value is set to the result of the callback function and returned.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
)

func main() {
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

}
```

#### Output

```
one
one
something else
```

</p>
</details>

### func \(\*Cache\[T\]\) [Keys](<https://github.com/atomicgo/cache/blob/main/cache.go#L223>)

```go
func (c *Cache[T]) Keys() []string
```

Keys returns a slice of all keys in the cache, expired or not.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"sort"
)

func main() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Get the keys.
	keys := c.Keys()
	sort.Strings(keys)

	fmt.Println(keys)

}
```

#### Output

```
[1 2 3]
```

</p>
</details>

### func \(\*Cache\[T\]\) [Purge](<https://github.com/atomicgo/cache/blob/main/cache.go#L181>)

```go
func (c *Cache[T]) Purge()
```

Purge removes all keys from the cache.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
)

func main() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Purge the cache.
	c.Purge()

	// Check if the cache is empty.
	fmt.Println(c.Size())

}
```

#### Output

```
0
```

</p>
</details>

### func \(\*Cache\[T\]\) [PurgeExpired](<https://github.com/atomicgo/cache/blob/main/cache.go#L189>)

```go
func (c *Cache[T]) PurgeExpired()
```

PurgeExpired removes all expired keys from the cache.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
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

}
```

#### Output

```
3
1
```

</p>
</details>

### func \(\*Cache\[T\]\) [Set](<https://github.com/atomicgo/cache/blob/main/cache.go#L76>)

```go
func (c *Cache[T]) Set(key string, value T, expiration ...time.Duration)
```

Set sets the value for a key. If the key already exists, the value is overwritten. If the Expiration time is 0, the default Expiration time is used.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"time"
)

func main() {
	c := cache.New[string]()

	// Set a value for a key.
	c.Set("1", "one")

	// Set a value for a key with a custom expiration.
	c.Set("2", "two", time.Second*10)
}
```

</p>
</details>

### func \(\*Cache\[T\]\) [SetExpiration](<https://github.com/atomicgo/cache/blob/main/cache.go#L155>)

```go
func (c *Cache[T]) SetExpiration(expiration time.Duration)
```

SetExpiration sets the Expiration time all keys in the cache.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
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

}
```

#### Output

```
20s
20s
20s
```

</p>
</details>

### func \(\*Cache\[T\]\) [Size](<https://github.com/atomicgo/cache/blob/main/cache.go#L201>)

```go
func (c *Cache[T]) Size() int
```

Size returns the number of keys in the cache, expired or not.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
)

func main() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Get the size.
	fmt.Println(c.Size())

}
```

#### Output

```
3
```

</p>
</details>

### func \(\*Cache\[T\]\) [StopAutoPurge](<https://github.com/atomicgo/cache/blob/main/cache.go#L265>)

```go
func (c *Cache[T]) StopAutoPurge()
```

StopAutoPurge stops the auto purge goroutine.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
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

}
```

#### Output

```
3
```

</p>
</details>

### func \(\*Cache\[T\]\) [ValidKeys](<https://github.com/atomicgo/cache/blob/main/cache.go#L237>)

```go
func (c *Cache[T]) ValidKeys() []string
```

ValidKeys returns a slice of all valid keys in the cache.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"sort"
	"time"
)

func main() {
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

}
```

#### Output

```
[1 2 3 expiresSoon]
[1 2 3]
```

</p>
</details>

### func \(\*Cache\[T\]\) [ValidSize](<https://github.com/atomicgo/cache/blob/main/cache.go#L209>)

```go
func (c *Cache[T]) ValidSize() int
```

ValidSize returns the number of keys in the cache that are not expired.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
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

}
```

#### Output

```
4
3
```

</p>
</details>

### func \(\*Cache\[T\]\) [Values](<https://github.com/atomicgo/cache/blob/main/cache.go#L251>)

```go
func (c *Cache[T]) Values() []T
```

Values returns a slice of all values in the cache.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"sort"
)

func main() {
	c := cache.New[string]()

	// Fill the cache with some values.
	c.Set("1", "one")
	c.Set("2", "two")
	c.Set("3", "three")

	// Get the values.
	values := c.Values()
	sort.Strings(values)

	fmt.Println(values)

}
```

#### Output

```
[one three two]
```

</p>
</details>

## type [Entry](<https://github.com/atomicgo/cache/blob/main/cache-entry.go#L6-L10>)

Entry is a cache entry.

```go
type Entry[T any] struct {
    Value      T
    CachedAt   time.Time
    Expiration time.Duration
}
```

### func \(Entry\[T\]\) [Expired](<https://github.com/atomicgo/cache/blob/main/cache-entry.go#L13>)

```go
func (e Entry[T]) Expired() bool
```

Expired returns if the Entry is expired.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
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

}
```

#### Output

```
false
true
```

</p>
</details>

### func \(\*Entry\[T\]\) [SetExpiration](<https://github.com/atomicgo/cache/blob/main/cache-entry.go#L21>)

```go
func (e *Entry[T]) SetExpiration(expiration time.Duration)
```

SetExpiration sets the Expiration time for the Entry.

<details><summary>Example</summary>
<p>

```go
package main

import (
	"atomicgo.dev/cache"
	"fmt"
	"time"
)

func main() {
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

}
```

#### Output

```
20s
```

</p>
</details>

## type [Options](<https://github.com/atomicgo/cache/blob/main/cache.go#L18-L21>)

Options can be passed to New to configure the cache.

```go
type Options struct {
    DefaultExpiration time.Duration
    AutoPurgeInterval time.Duration
}
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)


<!-- gomarkdoc:embed:end -->

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)

package lesson

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 4. Map Leak (Unbounded Cache)
type Cache struct {
	sync.RWMutex
	items map[string][]byte
}

// Bad: No cleanup mechanism
func (c *Cache) Set(key string, value []byte) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = value
}

// Good: With TTL and cleanup
type CacheItem struct {
	value     []byte
	timestamp time.Time
}

type BetterCache struct {
	sync.RWMutex
	items map[string]CacheItem
	ttl   time.Duration
}

func (c *BetterCache) Set(key string, value []byte) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = CacheItem{value: value, timestamp: time.Now()}
}

func (c *BetterCache) Cleanup() {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			c.Lock()
			now := time.Now()
			for k, v := range c.items {
				if now.Sub(v.timestamp) > c.ttl {
					delete(c.items, k)
				}
			}
			c.Unlock()
		}
	}()
}

func MapLeak() {
	var m runtime.MemStats

	cache := &Cache{items: make(map[string][]byte)}
	cache.Set("key", []byte("value"))
	runtime.ReadMemStats(&m)
	fmt.Println(m.Alloc)

	betterCache := &BetterCache{items: make(map[string]CacheItem), ttl: time.Minute}
	betterCache.Cleanup()
	runtime.ReadMemStats(&m)
	fmt.Println(m.Alloc)
}

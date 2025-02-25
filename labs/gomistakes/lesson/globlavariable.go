package lesson

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// 7. Global Variable Leak
var globalCache = make(map[string]*LargeObject) // Bad: Never cleaned up

// Good: Use sync.Map with cleanup mechanism
var betterCache sync.Map

func cleanup() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			betterCache.Range(func(key, value interface{}) bool {
				// Implement cleanup logic
				return true
			})
		}
	}()
}

func SetGlobalCache(key string, value *LargeObject) {
	betterCache.Store(key, value)
}

func GetGlobalCache(key string) *LargeObject {
	value, _ := betterCache.Load(key)
	return value.(*LargeObject)
}

func GlobalVariableLeak() {
	SetGlobalCache("key", &LargeObject{data: make([]byte, 1024*1024)})
	fmt.Println(len(GetGlobalCache("key").data))

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println(m.Alloc)
}

package lesson

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Leak struct {
	id       int
	name     string
	birthday time.Time
}

var leakSlice = []Leak{}

// 1. Goroutine Leak
func GoroutineLeak() {
	// Bad: Goroutine never exits
	go func() {
		ticker := time.NewTicker(time.Second)
		largeString := ""
		for {
			select {
			case <-ticker.C:
				newLeak := Leak{
					id:       rand.Intn(1000),
					name:     fmt.Sprintf("Leak %d", rand.Intn(1000)),
					birthday: time.Now(),
				}
				leakSlice = append(leakSlice, newLeak)
				largeString += "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				log.Printf("\nWorking... %d", len(leakSlice))
				ClosureLeak()
			}
		}
	}()
	select {}
}

// Good: Proper cancellation
func GoroutineLeakWithContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("Working...")
				
			case <-ctx.Done():
				return
			}
		}
	}()
	select {}
}

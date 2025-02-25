package lesson

import (
	"context"
	"fmt"
	"time"
)

// 5. Channel Leak
func ChannelLeak() {
	// Bad: Channel and goroutine never cleanup
	ch := make(chan int)
	go func() {
		val := <-ch // Blocked forever if nothing sends
		fmt.Println(val)
	}()

	// Good: Use context for cancellation
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch = make(chan int)
	go func() {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-ctx.Done():
			return
		}
	}()

}

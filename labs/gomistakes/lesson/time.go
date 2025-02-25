package lesson

import (
	"context"
	"fmt"
	"time"
)

// 6. Timer/Ticker Leak
func TimerLeak() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Bad: Timer never stopped
	timer := time.NewTimer(time.Hour)
	go func() {
		<-timer.C
		fmt.Println("Done!")
	}()

	// Good: Properly stop timer
	timer = time.NewTimer(time.Second)
	defer timer.Stop()
	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("Good:Done!")
			case <-ctx.Done():
				return
			}
		}
	}()

	select {}
}

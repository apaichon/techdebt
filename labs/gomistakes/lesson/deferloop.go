package lesson

import (
	"fmt"
	"os"
)

// 9. Defer in Loop Leak
func DeferInLoopLeak() {
	// Bad: Defers accumulate until function returns
	/* for i := 0; i < 100_000; i++ {
		file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			continue
		}
		defer file.Close() // Won't be called until function returns

		// Write to file
		if _, err := file.WriteString(fmt.Sprintf("Line %d\n", i)); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
	}*/

	// Good: Close in the same loop iteration
	for i := 0; i < 100_000; i++ {
		func() {
			file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				return
			}
			defer file.Close() // Called when anonymous function returns

			// Write to file
			if _, err := file.WriteString(fmt.Sprintf("Line %d\n", i)); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
			}
		}()
	}
}

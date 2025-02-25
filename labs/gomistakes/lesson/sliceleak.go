package lesson

import "fmt"

// 3. Slice Leak
func SliceLeak() {
	// Bad: Original array stays in memory
	data := make([]int, 1000000)
	small := data[len(data)-3:]

	// Good: Copy only what's needed
	small = make([]int, 3)
	copy(small, data[len(data)-3:])
	fmt.Println(len(small))
}

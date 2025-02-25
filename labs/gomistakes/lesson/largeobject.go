package lesson

import "fmt"

// 2. Closure Capturing Large Objects
type LargeObject struct {
	data []byte
}

func ClosureLeak() {
	obj := &LargeObject{
		data: make([]byte, 1024*1024), // 1MB
	}

	// Bad: Captures entire obj
	handler := func() {
		fmt.Println(len(obj.data))
	}
	handler()

	/*
	   // Good: Capture only what's needed
	   size := len(obj.data)
	   handler = func() {
	       fmt.Println(size)
	   }
	*/
}

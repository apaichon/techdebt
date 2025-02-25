package lesson

import (
	"fmt"
	"io"
	"net/http"
)

// 8. HTTP Response Body Leak
func HTTPBodyLeak() error {
	// Bad: Response body not closed
	resp, err := http.Get("https://example.com")
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(len(body))

	// Good: Always close response body
	resp, err = http.Get("https://example.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(len(body))
	return nil
}

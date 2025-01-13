package main

import (
	"fmt"
	"time"
)

func asyncOperation(callback func(), delay time.Duration) {
	time.Sleep(delay)
	callback()
}

func main() {
	done := make(chan struct{})

	fmt.Println("Starting...")
	asyncOperation(func() {
		fmt.Println("Async operation completed.")
		close(done)
	}, 2*time.Second)
	fmt.Println("Doing other work...")

	<-done
	fmt.Println("All operations complete.")
}

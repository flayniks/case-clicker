//go:build !windows

package main

import (
	"fmt"
	"time"
)

// The non-Windows harness never actually installs; it exists so the server
// and update code can be exercised in tests.
func runInstallerAndExit(path string) { fmt.Println("would install:", path) }

func main() {
	url, err := startServer()
	if err != nil {
		panic(err)
	}
	checkForUpdate() // synchronous here so tests can read the result immediately
	fmt.Println(url)
	WaitForPageGone(14*time.Second, 90*time.Second)
	fmt.Println("page gone; shutting down")
}

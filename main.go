package main

import (
	"fmt"
	"os"
)

func main() {
	var failed bool

	for path, errors := range DataErrors() {
		if len(errors) == 0 {
			fmt.Printf("✅ %s is valid\n", path)
		} else {
			failed = true
			fmt.Printf("⛔️ %s had errors:\n", path)
			for _, err := range errors {
				fmt.Printf("   - %s\n", err)
			}
		}
	}

	if failed {
		os.Exit(1)
	}
}

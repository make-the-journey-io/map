package main

import (
	"fmt"
	"os"
)

func main() {
	var failed bool

	for _, stage := range LoadMap().stages {
		if len(stage.errors) == 0 {
			fmt.Printf("✅ %s is valid\n", stage.path)
		} else {
			failed = true
			fmt.Printf("⛔️ %s had errors:\n", stage.path)
			for _, err := range stage.errors {
				fmt.Printf("   - %s\n", err)
			}
		}

		for _, link := range stage.Requires {
			fmt.Printf("   - requires: '%s' ➡ %s\n", link.stage.DisplayName, link.stage.path)
		}
	}

	if failed {
		os.Exit(1)
	}
}

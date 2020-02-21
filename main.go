package main

import (
	"fmt"
	"os"
)

func showLinks(links []Link) {
	for _, link := range links {
		fmt.Printf("   - requires: '%s' ➡ %s\n", link.stage.DisplayName, link.stage.path)
	}
}

func main() {
	var failed bool

	for _, stage := range LoadMap().stages {
		if len(stage.errors) == 0 {
			fmt.Printf("✅ %s is valid\n", stage.path)
			showLinks(stage.Requires)
		} else {
			failed = true
			fmt.Printf("⛔️ %s had errors:\n", stage.path)
			for _, err := range stage.errors {
				fmt.Printf("   - %s\n", err)
			}
		}
	}

	if failed {
		os.Exit(1)
	}
}

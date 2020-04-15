package main

import (
	"fmt"
	"os"
)

func showSuccess(s *Stage) {
	fmt.Printf("✅ %s is valid\n", s.path)
	for _, link := range s.Requires {
		fmt.Printf("   - requires: '%s' ➡ %s\n", link.stage.DisplayName, link.stage.path)
	}
}

func showFailure(s *Stage) {
	fmt.Printf("⛔️ %s had errors:\n", s.path)
	for _, err := range s.errors {
		fmt.Printf("   - %s\n", err)
	}
}

func main() {
	var failed bool

	for _, stage := range LoadMap().stages {
		if len(stage.errors) == 0 {
			showSuccess(stage)
		} else {
			failed = true
			showFailure(stage)
		}
	}

	if failed {
		os.Exit(1)
	}
}

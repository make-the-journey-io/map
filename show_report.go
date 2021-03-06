package main

import (
	"fmt"
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

// ShowValidation prints ✅ or detailed ⛔️ results to stdout about each node in the map
func ShowValidation(m *JourneyMap) {
	for _, stage := range m.stages {
		if len(stage.errors) == 0 {
			showSuccess(stage)
		} else {
			showFailure(stage)
		}
	}
}

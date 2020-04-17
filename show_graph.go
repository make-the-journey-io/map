package main

import "fmt"

func printConnections(s *Stage) {
	for _, r := range s.Requires {
		fmt.Printf(`  "%s"->"%s" [label="%s"]`+"\n", s.id, r.stage.id, "requires")
	}

	for _, r := range s.RelatesTo {
		fmt.Printf(`  "%s"->"%s" [label="%s"]`+"\n", s.id, r.stage.id, "relates to")
	}
}

// ShowGraph prints the map as a directed graph in Grapviz DOT format to stdout
func ShowGraph(m *JourneyMap) {
	fmt.Println("digraph G {")
	fmt.Println("  node [shape=box];")

	for _, s := range m.stages {
		fmt.Println()
		fmt.Printf(`  "%s" [label="%s"];`, s.id, s.DisplayName)
		fmt.Println()
		printConnections(s)
	}

	fmt.Println("}")
}

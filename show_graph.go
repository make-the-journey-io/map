package main

import "fmt"

func printConnections(s *Stage) {
	for _, r := range s.Requires {
		fmt.Printf(`  "%s"->"%s" [label="%s"]`+"\n", s.id, r.stage.id, "requires")
	}
}

// ShowGraph prints the map as a directed graph in Grapviz DOT format to stdout
func ShowGraph(m *JourneyMap) {
	fmt.Println("digraph G {")
	fmt.Println(`  node [shape="box" fontname="Roboto" fontsize="14" margin="0.15,0.10" height="0"];`)
	fmt.Println(`  edge [fontname="Roboto" fontsize="12"];`)

	for _, s := range m.stages {
		fmt.Println()
		fmt.Printf(`  "%s" [label="%s"];`, s.id, s.DisplayName)
		fmt.Println()
		printConnections(s)
	}

	fmt.Println("}")
}

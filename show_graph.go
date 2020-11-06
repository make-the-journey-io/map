package main

import (
	"fmt"
	"strings"
)

func resolveRelativeUrl(maybeRelative string, original string) string {
	url := maybeRelative
	if strings.Index(url, "#") == 0 {
		url = original + url
	}
	return url
}

func printConnections(s *Stage) {
	for _, r := range s.Requires {
		url := resolveRelativeUrl(r.CitedInURL, s.DefinitionURL)
		fmt.Printf(`  "%s"->"%s" [label="%s", URL="%s"]`+"\n", s.id, r.stage.id, "requires", url)
	}
}

// ShowGraph prints the map as a directed graph in Grapviz DOT format to stdout
func ShowGraph(m *JourneyMap) {
	fmt.Println("digraph G {")
	fmt.Println(`  node [shape="box" fontname="Roboto" fontsize="14" margin="0.15,0.10" height="0"];`)
	fmt.Println(`  edge [fontname="Roboto" fontsize="12"];`)

	for _, s := range m.stages {
		fmt.Println()
		fmt.Printf(`  "%s" [label="%s", URL="%s"];`, s.id, s.DisplayName, s.DefinitionURL)
		fmt.Println()
		printConnections(s)
	}

	fmt.Println("}")
}

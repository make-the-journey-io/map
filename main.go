package main

import (
	"flag"
	"os"
)

func main() {
	graphPtr := flag.Bool("graph", false, "Output a dependency graph.")
	flag.Parse()

	valid, jmap := LoadMap()
	if *graphPtr {
		ShowGraph(&jmap)
	} else {
		ShowValidation(&jmap)
	}

	if !valid {
		os.Exit(1)
	}
}

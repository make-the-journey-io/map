package main

import (
	"os"
	"path/filepath"
	"strings"
)

const stageDir = "data"

func pathToID(path string) string {
	return strings.Replace(strings.Replace(path, ".yaml", "", 1), stageDir+"/", "", 1)
}

func idToPath(id string) string {
	return stageDir + "/" + id + ".yaml"
}

func walk(f func(path string, info os.FileInfo, err error) error) {
	filepath.Walk("./"+stageDir, f)
}

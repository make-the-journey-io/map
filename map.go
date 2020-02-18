package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

// Stage represents the individual stages on the map
type Stage struct {
	path   string
	errors []error
}

// JourneyMap contains the complete definition of stages on the map
type JourneyMap struct {
	stages []*Stage
}

func loadStageFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func loadStage(schemaLoader gojsonschema.JSONLoader, path string) *Stage {
	stage := Stage{path: path}

	node, err := loadStageFile(path)
	if err != nil {
		stage.errors = append(stage.errors, err)
	}

	documentLoader := gojsonschema.NewBytesLoader(node)
	validationResult, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		stage.errors = append(stage.errors, err)
	}

	if !validationResult.Valid() {
		for _, desc := range validationResult.Errors() {
			stage.errors = append(stage.errors, errors.New(desc.String()))
		}
	}

	return &stage
}

// LoadMap builds the entire journey from the YAML data files
func LoadMap() JourneyMap {
	gojsonschema.FormatCheckers.Add("cross-referenced-data", CrossReferencedDataChecker{})
	schemaLoader := gojsonschema.NewReferenceLoader("file://./schema/node.json")

	m := JourneyMap{}

	filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		m.stages = append(m.stages, loadStage(schemaLoader, path))
		return nil
	})

	return m
}

package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	jsonyaml "github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

type Link struct {
	LinkTo     string `yaml:"link-to"`
	CitedInURL string `yaml:"cited-in-url"`
}

// Stage represents the individual stages on the map
type Stage struct {
	DisplayName   string `yaml:"display-name"`
	DefinitionURL string `yaml:"definition-url"`
	Requires      []Link `yaml:"requires"`

	path   string
	errors []error
}

// JourneyMap contains the complete definition of stages on the map
type JourneyMap struct {
	stages []*Stage
}

var stageSchemaLoader gojsonschema.JSONLoader

func init() {
	stageSchemaLoader = gojsonschema.NewReferenceLoader("file://./schema/node.json")
}

func loadStageFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data, err = jsonyaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func validateStage(content []byte) []error {
	var errs []error

	documentLoader := gojsonschema.NewBytesLoader(content)
	validationResult, err := gojsonschema.Validate(stageSchemaLoader, documentLoader)
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, desc := range validationResult.Errors() {
			errs = append(errs, errors.New(desc.String()))
		}
	}

	return errs
}

func loadStage(path string) *Stage {
	stage := Stage{path: path}

	content, err := loadStageFile(path)
	if err != nil {
		stage.errors = append(stage.errors, err)
		return &stage
	}

	schemaErrors := validateStage(content)
	stage.errors = append(stage.errors, schemaErrors...)

	err = yaml.Unmarshal(content, &stage)
	if err != nil {
		stage.errors = append(stage.errors, err)
	}

	return &stage
}

// LoadMap builds the entire journey from the YAML data files
func LoadMap() JourneyMap {
	m := JourneyMap{}

	filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		m.stages = append(m.stages, loadStage(path))
		return nil
	})

	return m
}

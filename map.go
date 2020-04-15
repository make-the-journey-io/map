package main

import (
	"errors"
	"io/ioutil"
	"os"

	jsonyaml "github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

// Link refers to another Stage by their identifier
type Link struct {
	LinkTo     string `yaml:"link-to"`
	CitedInURL string `yaml:"cited-in-url"`

	stage *Stage
}

// Stage represents the individual stages on the map
type Stage struct {
	DisplayName   string `yaml:"display-name"`
	DefinitionURL string `yaml:"definition-url"`
	Requires      []Link `yaml:"requires"`
	RelatesTo     []Link `yaml:"relates-to"`

	id     string
	path   string
	errors []error
}

// JourneyMap contains the complete definition of stages on the map
type JourneyMap struct {
	stages map[string]*Stage
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

func load(path string) *Stage {
	stage := Stage{path: path, id: pathToID(path)}

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

func resolveStage(m *JourneyMap, id string) *Stage {
	ref := m.stages[id]
	if ref == nil {
		ref = &Stage{id: "error", DisplayName: "Error"}
	}
	return ref
}

func resolveLinks(m *JourneyMap) {
	for _, s := range m.stages {
		for lid := range s.Requires {
			link := &s.Requires[lid]
			link.stage = resolveStage(m, link.LinkTo)
		}
		for rid := range s.RelatesTo {
			link := &s.RelatesTo[rid]
			link.stage = resolveStage(m, link.LinkTo)
		}
	}
}

// LoadMap builds the entire journey from the YAML data files
func LoadMap() (bool, JourneyMap) {
	m := JourneyMap{}
	m.stages = make(map[string]*Stage)
	valid := true

	walk(func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		stage := load(path)
		valid = valid && len(stage.errors) == 0
		m.stages[stage.id] = stage
		return nil
	})

	resolveLinks(&m)

	return valid, m
}

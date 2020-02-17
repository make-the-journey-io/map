package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

func loadNode(filename string) ([]byte, error) {
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

func checkNode(schemaLoader gojsonschema.JSONLoader, path string) []error {
	var errs []error

	node, err := loadNode(path)
	if err != nil {
		errs = append(errs, err)
	}

	documentLoader := gojsonschema.NewBytesLoader(node)
	validationResult, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		errs = append(errs, err)
	}

	if !validationResult.Valid() {
		for _, desc := range validationResult.Errors() {
			errs = append(errs, errors.New(desc.String()))
		}
	}

	return errs
}

// DataErrors checks whether the data conforms to the schema and returns all data file errors (if any)
func DataErrors() map[string][]error {
	gojsonschema.FormatCheckers.Add("cross-referenced-data", CrossReferencedDataChecker{})
	schemaLoader := gojsonschema.NewReferenceLoader("file://./schema/node.json")

	files := make(map[string][]error)

	filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		files[path] = checkNode(schemaLoader, path)
		return nil
	})

	return files
}

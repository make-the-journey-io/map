package main

import (
	"errors"
	"fmt"
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

func main() {
	schemaLoader := gojsonschema.NewReferenceLoader("file://./schema/node.json")
	var failed bool

	filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		errors := checkNode(schemaLoader, path)
		if len(errors) == 0 {
			fmt.Printf("✅ %s is valid\n", path)
		} else {
			failed = true
			fmt.Printf("⛔️ %s had errors:\n", path)
			for _, err := range errors {
				fmt.Printf("   - %s\n", err)
			}
		}

		return nil
	})

	if failed {
		os.Exit(1)
	}
}

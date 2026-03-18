package parser

import (
	"io"
	"os"
	"synctl/internal/domain/parser"

	"go.yaml.in/yaml/v3"
)

type YamlLoader struct {
}

func (l *YamlLoader) LoadAll(path string) ([]parser.ResourceYAML, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var resources []parser.ResourceYAML

	for {
		var doc parser.ResourceYAML

		err := decoder.Decode(&doc)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if doc.Kind == "" {
			continue
		}

		resources = append(resources, doc)
	}
	return resources, nil
}

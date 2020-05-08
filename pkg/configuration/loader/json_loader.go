package loader

import (
	"encoding/json"

	"github.com/stretchy/stretchy/pkg/configuration"
)

type JSONLoader struct {
	basePath string
}

func NewJSONLoader(basePath string) *JSONLoader {
	return &JSONLoader{basePath: basePath}
}

func (jl *JSONLoader) Supports() []string {
	return []string{"json"}
}

func (jl *JSONLoader) LoadAll() (configuration.IndexCollection, error) {
	indexCollection := configuration.NewIndexCollection()

	templates, err := listFilesByExtensions(jl.basePath, jl.Supports()...)
	if err != nil {
		return nil, err
	}

	for _, t := range templates {
		f, err := getFile(t)
		if err != nil {
			return nil, err
		}

		index := configuration.Index{}
		if err := json.Unmarshal([]byte(f.Content), &index); err != nil {
			return nil, err
		}

		indexCollection.Load(f.Name, index)
	}

	return indexCollection, nil
}

func (jl *JSONLoader) Load(configurationName string) (configuration.Index, error) {
	indexCollection, err := jl.LoadAll()
	if err != nil {
		return configuration.Index{}, err
	}

	return indexCollection.Get(configurationName)
}

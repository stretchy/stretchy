package loader

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/stretchy/stretchy/pkg/configuration"
	"gopkg.in/yaml.v3"

	"github.com/Masterminds/sprig/v3"
)

type YAMLLoader struct {
	basePath string
}

func NewYAMLLoader(basePath string) *YAMLLoader {
	return &YAMLLoader{basePath: basePath}
}

func (yl *YAMLLoader) Supports() []string {
	return []string{"yml", "yaml"}
}

func (yl *YAMLLoader) LoadAll() (configuration.IndexCollection, error) {
	indexCollection := configuration.NewIndexCollection()
	helpers, err := yl.loadHelpers()

	if err != nil {
		return nil, err
	}

	templater := template.New("yaml-loader")
	yl.initFunMap(templater)

	if len(helpers) > 0 {
		if _, err := templater.ParseFiles(helpers...); err != nil {
			return nil, err
		}
	}

	templates, err := yl.loadTemplates()
	if err != nil {
		return nil, err
	}

	for _, t := range templates {
		f, err := getFile(t)
		if err != nil {
			return nil, err
		}

		if _, err = templater.Parse(fmt.Sprintf("{{define `%s`}}%s{{end}}", f.Name, f.Content)); err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		if err := templater.ExecuteTemplate(buf, f.Name, nil); err != nil {
			return nil, err
		}

		index := configuration.Index{}
		if err := yaml.NewDecoder(buf).Decode(&index); err != nil {
			return nil, err
		}

		indexCollection.Load(f.Name, index)
	}

	return indexCollection, nil
}

func (yl *YAMLLoader) Load(configurationName string) (configuration.Index, error) {
	indexCollection, err := yl.LoadAll()
	if err != nil {
		return configuration.Index{}, err
	}

	return indexCollection.Get(configurationName)
}

func (yl *YAMLLoader) loadTemplates() ([]string, error) {
	templates, err := listFilesByExtensions(yl.basePath, yl.Supports()...)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func (yl *YAMLLoader) loadHelpers() ([]string, error) {
	helpers, err := listFilesByExtensions(yl.basePath, "tmpl")
	if err != nil {
		return nil, err
	}

	return helpers, nil
}

const recursionMaxNums = 1000

func (yl *YAMLLoader) initFunMap(t *template.Template) {
	funcMap := sprig.TxtFuncMap()

	includedNames := make(map[string]int)
	// Load some extra functionality

	funcMap["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder

		if v, ok := includedNames[name]; ok {
			if v > recursionMaxNums {
				return "", errors.Wrapf(
					fmt.Errorf("unable to execute template"),
					"rendering template has a nested reference name: %s",
					name,
				)
			}
			includedNames[name]++
		} else {
			includedNames[name] = 1
		}

		err := t.ExecuteTemplate(&buf, name, data)
		includedNames[name]--

		return buf.String(), err
	}

	t.Funcs(funcMap)
}

package configuration

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type Index struct {
	Mappings Mappings `json:"mappings,omitempty" yaml:"mappings"`
	Settings Settings `json:"settings,omitempty" yaml:"settings"`
}

func (i *Index) GetSettings() Settings {
	return i.Settings
}

func (i *Index) GetMappings() Mappings {
	return i.Mappings
}

func New(mappings Mappings, settings Settings) Index {
	index := Index{
		Mappings: mappings,
		Settings: settings,
	}

	index.Settings.CleanUp()

	return index
}

func (i *Index) UnmarshalYAML(value *yaml.Node) error {
	type yamlIndex Index

	yamlI := yamlIndex{}
	if err := value.Decode(&yamlI); err != nil {
		return err
	}

	i.Settings = yamlI.Settings
	i.Mappings = yamlI.Mappings

	i.Settings.CleanUp()

	return nil
}

func (i *Index) UnmarshalJSON(data []byte) error {
	type jsonIndex Index

	jsonI := jsonIndex{}

	if err := json.Unmarshal(data, &jsonI); err != nil {
		return err
	}

	i.Mappings = jsonI.Mappings
	i.Settings = jsonI.Settings

	i.Settings.CleanUp()

	return nil
}

func (i Index) Diff(mapping Index) (ChangeCollection, error) {
	changes := ChangeCollection{}

	settingsChanges, err := i.GetSettings().Diff(mapping.GetSettings())
	if err != nil {
		return nil, err
	}

	changes = append(changes, settingsChanges...)

	mappingsChanges, err := i.GetMappings().Diff(mapping.GetMappings())
	if err != nil {
		return nil, err
	}

	changes = append(changes, mappingsChanges...)

	return changes, nil
}

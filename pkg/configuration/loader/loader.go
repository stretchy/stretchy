package loader

import "github.com/stretchy/stretchy/pkg/configuration"

type Loader interface {
	LoadAll() (configuration.IndexCollection, error)
	Load(configurationName string) (configuration.Index, error)
	Supports() []string
}

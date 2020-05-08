package action

import (
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/configuration/loader"
)

type Load struct {
	basePath string
}

func NewLoad(
	basePath string,
) *Load {
	return &Load{
		basePath: basePath,
	}
}

func (load *Load) Load(configurationName string, format string) (configuration.Index, error) {
	l, err := loader.NewRegistry(load.basePath).GetByFormat(format)
	if err != nil {
		return configuration.Index{}, err
	}

	return l.Load(configurationName)
}

func (load *Load) LoadAll(format string) (configuration.IndexCollection, error) {
	l, err := loader.NewRegistry(load.basePath).GetByFormat(format)
	if err != nil {
		return nil, err
	}

	return l.LoadAll()
}

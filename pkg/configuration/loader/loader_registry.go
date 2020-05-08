package loader

import (
	"fmt"

	"github.com/stretchy/stretchy/pkg/utils"
)

type Registry struct {
	loaders []Loader
}

func NewRegistry(path string) *Registry {
	return &Registry{
		loaders: []Loader{
			NewYAMLLoader(path),
			NewJSONLoader(path),
		},
	}
}

func (r *Registry) GetByFormat(format string) (Loader, error) {
	for _, l := range r.loaders {
		if isTheLoader, _ := utils.InSlice(format, l.Supports()); isTheLoader {
			return l, nil
		}
	}

	return nil, fmt.Errorf("cannot find loader for format '%s'", format)
}

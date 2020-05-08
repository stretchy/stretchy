package configuration

import (
	"github.com/imdario/mergo"
	"github.com/r3labs/diff"
)

type Mappings map[string]interface{}

func (m Mappings) Merge(mappings Mappings) error {
	return mergo.Map(
		&m,
		mappings,
		mergo.WithOverride,
	)
}

func (m Mappings) Diff(mappings Mappings) (ChangeCollection, error) {
	mappingsChangeLogs, err := diff.Diff(m, mappings)
	if err != nil {
		return nil, err
	}

	changes := ChangeCollection{}

	for _, c := range mappingsChangeLogs {
		change := Change{
			Type: NewChangeTypeFromDiffType(c.Type),
			Path: append([]string{"mappings"}, c.Path...),
			From: c.From,
			To:   c.To,
		}

		if change.ShouldBeReported() {
			changes = append(changes, change)
		}
	}

	return changes, nil
}

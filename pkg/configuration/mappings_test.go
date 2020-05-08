package configuration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
)

func getBaseMappings() configuration.Mappings {
	return configuration.Mappings{
		"a_property":    "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
		"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
		"another_prop":  "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
	}
}

func TestMappings_Diff(t *testing.T) {
	mappings := getBaseMappings()
	anotherMappings := getBaseMappings()

	delete(anotherMappings, "another_prop")
	anotherMappings["template_prop"] = "cdda7dc9-9190-4cf0-8262-0411210ca0da"
	anotherMappings["new_property"] = "a7b8c4bb-5197-49c5-846c-514cecad5989"

	changes, err := mappings.Diff(anotherMappings)
	assert.NoError(t, err)

	assert.Len(t, changes, 3)

	assert.Contains(t, changes, configuration.Change{
		Type: configuration.ChangeTypeDelete,
		Path: []string{"mappings", "another_prop"},
		From: "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
		To:   nil,
	})
	assert.Contains(t, changes, configuration.Change{
		Type: configuration.ChangeTypeUpdate,
		Path: []string{"mappings", "template_prop"},
		From: "A9201516-B322-49B2-A284-5186FA43A306",
		To:   "cdda7dc9-9190-4cf0-8262-0411210ca0da",
	})
	assert.Contains(t, changes, configuration.Change{
		Type: configuration.ChangeTypeCreate,
		Path: []string{"mappings", "new_property"},
		From: nil,
		To:   "a7b8c4bb-5197-49c5-846c-514cecad5989",
	})
}

func TestMappings_Merge(t *testing.T) {
	mappings := getBaseMappings()

	err := mappings.Merge(
		configuration.Mappings{
			"new_property": map[string]interface{}{
				"type": "aType",
			},
		},
	)

	assert.NoError(t, err)

	assert.NotEqual(t, mappings, getBaseMappings())
	assert.Equal(
		t,
		configuration.Mappings{
			"a_property":   "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
			"another_prop": "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
			"new_property": map[string]interface{}{
				"type": "aType",
			},
			"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
		},
		mappings,
	)

	err = mappings.Merge(
		configuration.Mappings{
			"new_property": map[string]interface{}{
				"type":         "anotherType",
				"aNewProperty": "24715c47-3672-4d60-bb1f-ce0078da68d0",
			},
		},
	)

	assert.NoError(t, err)

	assert.Equal(
		t,
		configuration.Mappings{
			"a_property":   "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
			"another_prop": "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
			"new_property": map[string]interface{}{
				"type":         "anotherType",
				"aNewProperty": "24715c47-3672-4d60-bb1f-ce0078da68d0",
			},
			"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
		},
		mappings,
	)
}

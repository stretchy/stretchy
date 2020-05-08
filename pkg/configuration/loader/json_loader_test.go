package loader_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/configuration/loader"
)

func getJSONBaseTestAExpectedValue() configuration.Index {
	return configuration.New(
		configuration.Mappings{
			"a_property":    "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
			"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
			"another_prop":  "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
		},
		configuration.Settings{
			"a_setting":       "a",
			"another_setting": float64(2),
			"third_setting":   "asd",
		},
	)
}

func getJSONBaseTestBExpectedValue() configuration.Index {
	return configuration.New(
		map[string]interface{}{
			"field": "bf0a1a2a-e1a2-4ad9-8b46-d8a7f3a12370",
		},
		map[string]interface{}{
			"a_setting":       "a",
			"another_setting": float64(2),
			"third_setting":   "asd",
		},
	)
}

func TestJSONLoader_LoadAll(t *testing.T) {
	jsonLoader := loader.NewJSONLoader(getScenarioPath(t, "json-base"))

	mappingCollection, err := jsonLoader.LoadAll()

	assert.NoError(t, err)

	assert.Equal(
		t,
		configuration.IndexCollection{
			"test-a": getJSONBaseTestAExpectedValue(),
			"test-b": getJSONBaseTestBExpectedValue(),
		},
		mappingCollection,
	)
}

func TestJSONLoader_Load(t *testing.T) {
	jsonLoader := loader.NewJSONLoader(getScenarioPath(t, "json-base"))

	testCases := []struct {
		name          string
		expectedValue configuration.Index
	}{
		{
			name:          "test-a",
			expectedValue: getJSONBaseTestAExpectedValue(),
		},
		{
			name:          "test-b",
			expectedValue: getJSONBaseTestBExpectedValue(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			index, err := jsonLoader.Load(tc.name)
			assert.NoError(t, err)
			assert.Equal(
				t,
				tc.expectedValue,
				index,
			)
		})
	}
}

func TestJSONLoader_LoadAll_WithSyntaxError(t *testing.T) {
	jsonLoader := loader.NewJSONLoader(getScenarioPath(t, "json-syntax-error"))

	mappingCollection, err := jsonLoader.LoadAll()

	assert.Error(t, err)

	assert.Nil(
		t,
		mappingCollection,
	)
}

func TestJSONLoader_Load_WithSyntaxError(t *testing.T) {
	jsonLoader := loader.NewJSONLoader(getScenarioPath(t, "json-syntax-error"))

	mappingCollection, err := jsonLoader.Load("test-b")

	assert.Error(t, err)

	assert.Equal(
		t,
		configuration.Index{},
		mappingCollection,
	)
}

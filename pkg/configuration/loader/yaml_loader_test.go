package loader_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/configuration/loader"
)

func getYAMLBaseTestAExpectedValue() configuration.Index {
	return configuration.New(
		map[string]interface{}{
			"a_property":    "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
			"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
			"another_prop":  "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
		},
		map[string]interface{}{
			"a_setting":       "a",
			"another_setting": 2,
			"third_setting":   "asd",
		},
	)
}

func getYAMLBaseTestBExpectedValue() configuration.Index {
	return configuration.New(
		map[string]interface{}{
			"field": "bf0a1a2a-e1a2-4ad9-8b46-d8a7f3a12370",
		},
		map[string]interface{}{
			"a_setting":       "a",
			"another_setting": 2,
			"third_setting":   "asd",
		},
	)
}

func TestYAMLLoader_LoadAll(t *testing.T) {
	yamlLoader := loader.NewYAMLLoader(getScenarioPath(t, "yaml-base"))

	mappingCollection, err := yamlLoader.LoadAll()

	assert.NoError(t, err)

	assert.EqualValues(
		t,
		configuration.IndexCollection{
			"test-a": getYAMLBaseTestAExpectedValue(),
			"test-b": getYAMLBaseTestBExpectedValue(),
		},
		mappingCollection,
	)
}

func TestYAMLLoader_Load(t *testing.T) {
	yamlLoader := loader.NewYAMLLoader(getScenarioPath(t, "yaml-base"))

	testCases := []struct {
		name          string
		expectedValue configuration.Index
	}{
		{
			name:          "test-a",
			expectedValue: getYAMLBaseTestAExpectedValue(),
		},
		{
			name:          "test-b",
			expectedValue: getYAMLBaseTestBExpectedValue(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			index, err := yamlLoader.Load(tc.name)
			assert.NoError(t, err)
			assert.Equal(
				t,
				tc.expectedValue,
				index,
			)
		})
	}
}

func TestYAMLLoader_LoadAll_WithSyntaxError(t *testing.T) {
	yamlLoader := loader.NewYAMLLoader(getScenarioPath(t, "yaml-syntax-error"))

	mappingCollection, err := yamlLoader.LoadAll()

	assert.Error(t, err)

	assert.Nil(
		t,
		mappingCollection,
	)
}

func TestYAMLLoader_Load_WithSyntaxError(t *testing.T) {
	yamlLoader := loader.NewYAMLLoader(getScenarioPath(t, "yaml-syntax-error"))

	mappingCollection, err := yamlLoader.Load("test-b")

	assert.Error(t, err)

	assert.Equal(
		t,
		configuration.Index{},
		mappingCollection,
	)
}

func TestYAMLLoader_LoadAll_WithRecursion(t *testing.T) {
	yamlLoader := loader.NewYAMLLoader(getScenarioPath(t, "yaml-recursion"))

	mappingCollection, err := yamlLoader.LoadAll()

	assert.Error(t, err)

	assert.Nil(
		t,
		mappingCollection,
	)
}

func TestYAMLLoader_Load_WithRecursion(t *testing.T) {
	yamlLoader := loader.NewYAMLLoader(getScenarioPath(t, "yaml-recursion"))

	mappingCollection, err := yamlLoader.Load("test-b")

	assert.Error(t, err)

	assert.Equal(
		t,
		configuration.Index{},
		mappingCollection,
	)
}

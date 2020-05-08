package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/action"
	"github.com/stretchy/stretchy/pkg/configuration"
)

func TestLoad_Load(t *testing.T) {
	testCases := []struct {
		basePath          string
		format            string
		configurationName string
		expectedIndex     configuration.Index
	}{
		{
			basePath:          getScenarioPath(t, "json-base"),
			format:            "json",
			configurationName: "test-a",
			expectedIndex: configuration.Index{
				Mappings: configuration.Mappings{
					"a_property":    "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
					"another_prop":  "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
					"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
				},
				Settings: configuration.Settings{
					"a_setting":       "a",
					"another_setting": float64(2),
					"third_setting":   "asd",
				},
			},
		},
		{
			basePath:          getScenarioPath(t, "yaml-base"),
			format:            "yaml",
			configurationName: "test-a",
			expectedIndex: configuration.Index{
				Mappings: configuration.Mappings{
					"a_property":    "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
					"another_prop":  "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
					"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
				},
				Settings: configuration.Settings{
					"a_setting":       "a",
					"another_setting": 2,
					"third_setting":   "asd",
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.format, func(t *testing.T) {
			loadAction := action.NewLoad(testCase.basePath)

			index, err := loadAction.Load(testCase.configurationName, testCase.format)
			assert.NoError(t, err)
			assert.Equal(
				t,
				testCase.expectedIndex,
				index,
			)
		})
	}
}

func TestLoad_LoadAll(t *testing.T) {
	testCases := []struct {
		basePath           string
		format             string
		expectedCollection configuration.IndexCollection
	}{
		{
			basePath: getScenarioPath(t, "json-base"),
			format:   "json",
			expectedCollection: configuration.IndexCollection{
				"test-a": configuration.Index{
					Mappings: configuration.Mappings{
						"a_property":    "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
						"another_prop":  "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
						"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
					},
					Settings: configuration.Settings{
						"a_setting":       "a",
						"another_setting": float64(2),
						"third_setting":   "asd",
					},
				},
				"test-b": configuration.Index{
					Mappings: configuration.Mappings{
						"field": "bf0a1a2a-e1a2-4ad9-8b46-d8a7f3a12370",
					},
					Settings: configuration.Settings{
						"a_setting":       "a",
						"another_setting": float64(2),
						"third_setting":   "asd",
					},
				},
			},
		},
		{
			basePath: getScenarioPath(t, "yaml-base"),
			format:   "yaml",
			expectedCollection: configuration.IndexCollection{
				"test-a": configuration.Index{
					Mappings: configuration.Mappings{
						"a_property":    "e0a6877f-f62a-4cd8-9ca9-656a3a4a0fe5",
						"another_prop":  "d82c950f-ebf5-4aa9-822e-992da5ecc69e",
						"template_prop": "A9201516-B322-49B2-A284-5186FA43A306",
					},
					Settings: configuration.Settings{
						"a_setting":       "a",
						"another_setting": 2,
						"third_setting":   "asd",
					},
				},
				"test-b": configuration.Index{
					Mappings: configuration.Mappings{
						"field": "bf0a1a2a-e1a2-4ad9-8b46-d8a7f3a12370",
					},
					Settings: configuration.Settings{
						"a_setting": "a", "another_setting": 2,
						"third_setting": "asd",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.format, func(t *testing.T) {
			loadLocaleAction := action.NewLoad(testCase.basePath)
			indexCollection, err := loadLocaleAction.LoadAll(testCase.format)
			assert.NoError(t, err)
			assert.Equal(
				t,
				testCase.expectedCollection,
				indexCollection,
			)
		})
	}
}

func TestLoad_Load_NoRegistry(t *testing.T) {
	loadLocaleAction := action.NewLoad(getScenarioPath(t, "toml-base"))

	index, err := loadLocaleAction.Load("test-b", "toml")

	assert.Error(t, err)
	assert.Equal(t, configuration.Index{}, index)
}

func TestLoad_LoadAll_NoRegistry(t *testing.T) {
	loadLocaleAction := action.NewLoad(getScenarioPath(t, "toml-base"))

	indexCollection, err := loadLocaleAction.LoadAll("toml")

	assert.Error(t, err)
	assert.Nil(t, indexCollection)
}

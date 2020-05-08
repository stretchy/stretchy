package configuration_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
	"gopkg.in/yaml.v3"
)

const indexExampleName = "test-index"

func getIndexExample() configuration.Index {
	return configuration.New(
		configuration.Mappings{
			"_doc": map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": map[string]interface{}{
						"type": "text",
					},
				},
			},
		},
		configuration.Settings{
			"index": map[string]interface{}{
				"number_of_shards": 5,
			},
		},
	)
}

func createConfiguration(
	mappings map[string]interface{},
	settings map[string]interface{},
) configuration.Index {
	return configuration.New(mappings, settings)
}

func TestIndex_Diff(t *testing.T) {
	testCases := []struct {
		name                 string
		currentConfiguration configuration.Index
		newConfiguration     configuration.Index
		expectedChanges      configuration.ChangeCollection
	}{
		{
			name: "prune metadata",
			currentConfiguration: createConfiguration(
				map[string]interface{}{},
				map[string]interface{}{
					"index": map[string]interface{}{
						"creation_date":    "mar  7 apr 2020, 16.19.24",
						"provided_name":    "test-123-test",
						"uuid":             "e4232e2f-3a7d-4819-bc45-d40481d68929",
						"version":          2,
						"number_of_shards": 5,
					},
				},
			),
			newConfiguration: createConfiguration(
				map[string]interface{}{},
				map[string]interface{}{
					"index": map[string]interface{}{
						"number_of_shards": 5,
					},
				},
			),
			expectedChanges: configuration.ChangeCollection{},
		},
		{
			name:                 "change shard numbers",
			currentConfiguration: getIndexExample(),
			newConfiguration: func(index configuration.Index) configuration.Index {
				index.GetSettings()["index"].(map[string]interface{})["number_of_shards"] = 1
				return index
			}(getIndexExample()),
			expectedChanges: configuration.ChangeCollection{
				configuration.Change{
					Type: configuration.ChangeTypeUpdate,
					Path: []string{"settings", "index", "number_of_shards"},
					From: 5,
					To:   1,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			changes, err := tc.currentConfiguration.Diff(tc.newConfiguration)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedChanges, changes)
		})
	}
}

func TestIndex_Reorganize(t *testing.T) {
	index := createConfiguration(
		map[string]interface{}{},
		map[string]interface{}{
			"index": map[string]interface{}{
				"creation_date": "51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
				"provided_name": "51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
				"uuid":          "51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
				"version":       "51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
			},
			"number_of_shards":                  "b41cb716-0583-42eb-981b-9756c8ffc317",
			"shard.check_on_startup":            "9630e647-1649-439b-a7c9-1bbfe17747c4",
			"codec":                             "a3ef3b19-0330-4c1c-8326-d56b704917bb",
			"routing_partition_size":            "81459146-8af3-4a81-a8a8-1879a69c3e5a",
			"load_fixed_bitset_filters_eagerly": "643a99df-d525-4cb7-a997-3cf0709fea53",
			"number_of_replicas":                "72411a53-3dbf-4e21-b666-6178e51f86b0",
			"auto_expand_replicas":              "161db13f-5d41-4b86-b658-7578291b768f",
			"refresh_interval":                  "cba8f922-d151-4af8-a91e-531fb8f83877",
			"max_result_window":                 "91b3d6a3-3e0f-473b-a944-d8ff3362abe5",
			"max_inner_result_window":           "6a2dea2a-5002-4926-bbe6-c662ec1b9cfe",
			"max_rescore_window":                "e29ea0ad-c994-428b-9ab5-d2510350e17d",
			"max_docvalue_fields_search":        "21513e3f-481d-4442-8d3f-e3fa7fd44bf8",
			"max_script_fields":                 "c0b9f2e8-3eff-453f-a8d2-38d4882b6b93",
			"max_ngram_diff":                    "a861af1e-3a31-4453-953f-bbb55aa71296",
			"max_shingle_diff":                  "8d4030e3-cd26-457b-bdbb-97cae8c21f17",
			"blocks.read_only":                  "b10518ce-3b10-47dc-895e-e87a66a8c7fc",
			"blocks.read_only_allow_delete":     "8df4fd9e-480c-4d1d-956b-08b0941a949c",
			"blocks.read":                       "1a2dc574-62c3-4ace-9dd1-10a377f54163",
			"blocks.write":                      "79a6ac88-e23b-4c4b-b54f-fe4d718cee81",
			"blocks.metadata":                   "63e5f3f8-d71e-4736-92d9-270cd78b7da8",
			"max_refresh_listeners":             "2c684150-abc8-44a5-93c0-2584ae40cfe3",
			"highlight.max_analyzed_offset":     "13d876ca-860e-4322-88e4-4f5ace596a58",
			"max_terms_count":                   "3f301c70-5fa9-4809-95ef-1df4b27fab33",
			"routing.allocation.enable":         "26a16b66-2ead-4083-8efb-276560077384",
			"routing.rebalance.enable":          "98e78561-86ce-41d7-944f-598b2e007f4f",
			"gc_deletes":                        "be57a2ca-2d1d-4c37-aee8-d66bb99d1f0e",
			"max_regex_length":                  "d9843625-cfe7-46d8-b4e7-2433961f7496",
			"default_pipeline":                  "06ab44c8-d7f0-4d59-9bdb-74f83dc64aea",
		},
	)

	assert.Equal(
		t,
		configuration.Settings{
			"index": map[string]interface{}{
				"number_of_shards":                  "b41cb716-0583-42eb-981b-9756c8ffc317",
				"shard.check_on_startup":            "9630e647-1649-439b-a7c9-1bbfe17747c4",
				"codec":                             "a3ef3b19-0330-4c1c-8326-d56b704917bb",
				"routing_partition_size":            "81459146-8af3-4a81-a8a8-1879a69c3e5a",
				"load_fixed_bitset_filters_eagerly": "643a99df-d525-4cb7-a997-3cf0709fea53",
				"number_of_replicas":                "72411a53-3dbf-4e21-b666-6178e51f86b0",
				"auto_expand_replicas":              "161db13f-5d41-4b86-b658-7578291b768f",
				"refresh_interval":                  "cba8f922-d151-4af8-a91e-531fb8f83877",
				"max_result_window":                 "91b3d6a3-3e0f-473b-a944-d8ff3362abe5",
				"max_inner_result_window":           "6a2dea2a-5002-4926-bbe6-c662ec1b9cfe",
				"max_rescore_window":                "e29ea0ad-c994-428b-9ab5-d2510350e17d",
				"max_docvalue_fields_search":        "21513e3f-481d-4442-8d3f-e3fa7fd44bf8",
				"max_script_fields":                 "c0b9f2e8-3eff-453f-a8d2-38d4882b6b93",
				"max_ngram_diff":                    "a861af1e-3a31-4453-953f-bbb55aa71296",
				"max_shingle_diff":                  "8d4030e3-cd26-457b-bdbb-97cae8c21f17",
				"blocks.read_only":                  "b10518ce-3b10-47dc-895e-e87a66a8c7fc",
				"blocks.read_only_allow_delete":     "8df4fd9e-480c-4d1d-956b-08b0941a949c",
				"blocks.read":                       "1a2dc574-62c3-4ace-9dd1-10a377f54163",
				"blocks.write":                      "79a6ac88-e23b-4c4b-b54f-fe4d718cee81",
				"blocks.metadata":                   "63e5f3f8-d71e-4736-92d9-270cd78b7da8",
				"max_refresh_listeners":             "2c684150-abc8-44a5-93c0-2584ae40cfe3",
				"highlight.max_analyzed_offset":     "13d876ca-860e-4322-88e4-4f5ace596a58",
				"max_terms_count":                   "3f301c70-5fa9-4809-95ef-1df4b27fab33",
				"routing.allocation.enable":         "26a16b66-2ead-4083-8efb-276560077384",
				"routing.rebalance.enable":          "98e78561-86ce-41d7-944f-598b2e007f4f",
				"gc_deletes":                        "be57a2ca-2d1d-4c37-aee8-d66bb99d1f0e",
				"max_regex_length":                  "d9843625-cfe7-46d8-b4e7-2433961f7496",
				"default_pipeline":                  "06ab44c8-d7f0-4d59-9bdb-74f83dc64aea",
			},
		},
		index.GetSettings(),
	)
}

func TestIndex_UnmarshalJSON(t *testing.T) {
	rawJSON := json.RawMessage(`
{
   "settings":{
      "index":{
         "creation_date":"51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
         "provided_name":"51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
         "uuid":"51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
         "version":"51db18ff-04d4-440e-94b9-2a9e24c0b1e8",
         "max_terms_count":"3f301c70-5fa9-4809-95ef-1df4b27fab33"
      },
      "gc_deletes":"be57a2ca-2d1d-4c37-aee8-d66bb99d1f0e",
      "max_regex_length":"d9843625-cfe7-46d8-b4e7-2433961f7496",
      "default_pipeline":"06ab44c8-d7f0-4d59-9bdb-74f83dc64aea"
   },
   "mappings":{
      "properties":{
         "field1":{
            "type":"text"
         }
      }
   }
}
	`)

	index := configuration.Index{}

	err := index.UnmarshalJSON(rawJSON)
	assert.NoError(t, err)

	assert.Equal(
		t,
		configuration.Index{
			Settings: map[string]interface{}{
				"index": map[string]interface{}{
					"max_terms_count":  "3f301c70-5fa9-4809-95ef-1df4b27fab33",
					"gc_deletes":       "be57a2ca-2d1d-4c37-aee8-d66bb99d1f0e",
					"max_regex_length": "d9843625-cfe7-46d8-b4e7-2433961f7496",
					"default_pipeline": "06ab44c8-d7f0-4d59-9bdb-74f83dc64aea",
				},
			},
			Mappings: map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": map[string]interface{}{
						"type": "text",
					},
				},
			},
		},
		index,
	)
}

func TestIndex_UnmarshalYAML(t *testing.T) {
	rawYaml := `
settings:
  index:
    creation_date: "51db18ff-04d4-440e-94b9-2a9e24c0b1e8"
    provided_name: "51db18ff-04d4-440e-94b9-2a9e24c0b1e8"
    uuid: "51db18ff-04d4-440e-94b9-2a9e24c0b1e8"
    version: "51db18ff-04d4-440e-94b9-2a9e24c0b1e8"
    max_terms_count: "3f301c70-5fa9-4809-95ef-1df4b27fab33"
  gc_deletes: be57a2ca-2d1d-4c37-aee8-d66bb99d1f0e
  max_regex_length: d9843625-cfe7-46d8-b4e7-2433961f7496
  default_pipeline: 06ab44c8-d7f0-4d59-9bdb-74f83dc64aea
mappings:
  properties:
    field1:
      type: text
`
	index := configuration.Index{}

	r := bytes.NewReader([]byte(rawYaml))

	err := yaml.NewDecoder(r).Decode(&index)
	assert.NoError(t, err)

	assert.Equal(
		t,
		configuration.Index{
			Settings: configuration.Settings{
				"index": map[string]interface{}{
					"default_pipeline": "06ab44c8-d7f0-4d59-9bdb-74f83dc64aea",
					"gc_deletes":       "be57a2ca-2d1d-4c37-aee8-d66bb99d1f0e",
					"max_terms_count":  "3f301c70-5fa9-4809-95ef-1df4b27fab33",
					"max_regex_length": "d9843625-cfe7-46d8-b4e7-2433961f7496",
				},
			},
			Mappings: configuration.Mappings{
				"properties": configuration.Mappings{
					"field1": configuration.Mappings{
						"type": "text",
					},
				},
			},
		},
		index,
	)
}

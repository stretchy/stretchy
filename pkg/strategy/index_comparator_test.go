package strategy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/strategy"
)

func getIndexExample() *configuration.Index {
	i := configuration.New(
		map[string]interface{}{
			"properties": map[string]interface{}{
				"field1": map[string]interface{}{
					"type": "text",
				},
			},
		},
		map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_replicas":   1,
				"auto_expand_replicas": 2,
				"search": map[string]interface{}{
					"idle.after": 3,
				},
				"refresh_interval":           4,
				"max_result_window":          5,
				"max_inner_result_window":    6,
				"max_rescore_window":         7,
				"max_docvalue_fields_search": 8,
				"max_script_fields":          9,
				"max_ngram_diff":             10,
				"max_shingle_diff":           11,
				"blocks": map[string]interface{}{
					"read_only":              12,
					"read_only_allow_delete": 13,
					"read":                   14,
					"write":                  15,
					"metadata":               16,
				},
				"max_refresh_listeners": 17,
				"analyze": map[string]interface{}{
					"max_token_count": 18,
				},
				"highlight": map[string]interface{}{
					"max_analyzed_offset": 19,
				},
				"max_terms_count":  20,
				"max_regex_length": 21,
				"routing": map[string]interface{}{
					"allocation": map[string]interface{}{
						"enable": 22,
					},
					"rebalance": map[string]interface{}{
						"enable": 23,
					},
				},
				"gc_deletes":       24,
				"default_pipeline": 25,
				"final_pipeline":   26,
			},
		},
	)

	return &i
}

func TestIndexActionVoter_Compare_None(t *testing.T) {
	indexActionVoter := strategy.NewIndexActionVoter(true)

	result, err := indexActionVoter.Compare(getIndexExample(), getIndexExample())
	assert.NoError(t, err)

	assert.IsType(t, strategy.IndexVoterResult{}, result)
	assert.Equal(t, strategy.IndexDecisionNone, result.Action())
	assert.Nil(t, result.Changes())
}

func TestIndexActionVoter_Compare_Create(t *testing.T) {
	indexActionVoter := strategy.NewIndexActionVoter(true)
	result, err := indexActionVoter.Compare(nil, getIndexExample())
	assert.NoError(t, err)

	assert.IsType(t, strategy.IndexVoterResult{}, result)
	assert.Equal(t, strategy.IndexDecisionCreate, result.Action())
	assert.Nil(t, result.Changes())
}

func TestIndexActionVoter_Compare_Update(t *testing.T) {
	mappingWithANewField := getIndexExample()

	newFieldMapping := map[string]interface{}{
		"type": "keyword",
	}

	err := mappingWithANewField.GetMappings().
		Merge(
			map[string]interface{}{
				"properties": map[string]interface{}{
					"field2": newFieldMapping,
				},
			},
		)

	assert.NoError(t, err)

	testCases := []struct {
		name                     string
		allowSoftUpdate          bool
		newIndexConfiguration    *configuration.Index
		expectedChangeCollection configuration.ChangeCollection
		expectedDecision         strategy.IndexAction
	}{
		{
			name:                  "NewField SoftUpdate enabled",
			allowSoftUpdate:       true,
			newIndexConfiguration: mappingWithANewField,
			expectedChangeCollection: configuration.ChangeCollection{
				configuration.Change{
					Type: configuration.ChangeTypeCreate,
					Path: []string{"mappings", "properties", "field2"},
					From: nil,
					To:   newFieldMapping,
				},
			},
			expectedDecision: strategy.IndexDecisionUpdate,
		},
		{
			name:                  "NewField SoftUpdate disabled",
			allowSoftUpdate:       false,
			newIndexConfiguration: mappingWithANewField,
			expectedChangeCollection: configuration.ChangeCollection{
				configuration.Change{
					Type: configuration.ChangeTypeCreate,
					Path: []string{"mappings", "properties", "field2"},
					From: nil,
					To:   newFieldMapping,
				},
			},
			expectedDecision: strategy.IndexDecisionMigrate,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			indexActionVoter := strategy.NewIndexActionVoter(tc.allowSoftUpdate)
			result, err := indexActionVoter.Compare(getIndexExample(), tc.newIndexConfiguration)
			assert.NoError(t, err)

			assert.IsType(t, strategy.IndexVoterResult{}, result)
			assert.Equal(t, tc.expectedDecision, result.Action())
			assert.Equal(
				t,
				tc.expectedChangeCollection,
				result.Changes(),
			)
		})
	}
}

func TestIndexActionVoter_Compare_Migrate(t *testing.T) {
	newIndexConfigurationWithPropertyTypeChange := getIndexExample()

	newFieldMapping := map[string]interface{}{
		"type": "keyword",
	}

	err := newIndexConfigurationWithPropertyTypeChange.GetMappings().
		Merge(
			map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": newFieldMapping,
				},
			},
		)

	assert.NoError(t, err)

	newIndexConfigurationWithChangedSettings := getIndexExample()

	err = newIndexConfigurationWithChangedSettings.GetSettings().Merge(
		map[string]interface{}{
			"index": map[string]interface{}{
				"number_of_shards": 5,
			},
		},
	)
	assert.NoError(t, err)

	testCases := []struct {
		name                  string
		newIndexConfiguration *configuration.Index
		expectedDecision      strategy.IndexAction
		expectedChanges       configuration.ChangeCollection
	}{
		{
			name:                  "Update property type",
			expectedDecision:      strategy.IndexDecisionMigrate,
			newIndexConfiguration: newIndexConfigurationWithPropertyTypeChange,
			expectedChanges: configuration.ChangeCollection{
				configuration.Change{
					Type: configuration.ChangeTypeUpdate,
					Path: []string{"mappings", "properties", "field1", "type"},
					From: "text",
					To:   "keyword",
				},
			},
		},
		{
			name:                  "Change settings",
			expectedDecision:      strategy.IndexDecisionMigrate,
			newIndexConfiguration: newIndexConfigurationWithChangedSettings,
			expectedChanges: configuration.ChangeCollection{
				configuration.Change{
					Type: configuration.ChangeTypeCreate,
					Path: []string{"settings", "index", "number_of_shards"},
					From: nil,
					To:   5,
				},
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			indexActionVoter := strategy.NewIndexActionVoter(true)
			result, err := indexActionVoter.Compare(getIndexExample(), testCase.newIndexConfiguration)
			assert.NoError(t, err)

			assert.IsType(t, strategy.IndexVoterResult{}, result)
			assert.Equal(t, testCase.expectedDecision, result.Action())
			assert.Equal(
				t,
				testCase.expectedChanges,
				result.Changes(),
			)
		})
	}
}

func TestIndexDecision_String(t *testing.T) {
	testCases := []struct {
		actionName    string
		indexDecision strategy.IndexAction
	}{
		{
			actionName:    "None",
			indexDecision: strategy.IndexDecisionNone,
		},
		{
			actionName:    "Create",
			indexDecision: strategy.IndexDecisionCreate,
		},
		{
			actionName:    "Migrate",
			indexDecision: strategy.IndexDecisionMigrate,
		},
		{
			actionName:    "Update",
			indexDecision: strategy.IndexDecisionUpdate,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.actionName, func(t *testing.T) {
			assert.Equal(t, testCase.actionName, testCase.indexDecision.String())
		})
	}
}

func TestIndexActionVoter_Compare_ErrorOnWrongPropertyChange(t *testing.T) {
	index := getIndexExample()

	err := index.GetMappings().
		Merge(
			map[string]interface{}{
				"properties": map[string]interface{}{
					"field1": "wrong",
				},
			},
		)

	assert.NoError(t, err)

	indexActionVoter := strategy.NewIndexActionVoter(true)
	_, err = indexActionVoter.Compare(getIndexExample(), index)
	assert.Error(t, err)
}

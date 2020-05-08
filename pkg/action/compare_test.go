package action_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/action"
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/elasticsearch"
	"github.com/stretchy/stretchy/pkg/strategy"
)

const prefix = "myPrefix"
const indexName1 = "index-example"
const indexName2 = "another-index"

func getConfiguration1() configuration.Index {
	return configuration.New(
		configuration.Mappings{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "integer",
				},
				"uuid": map[string]interface{}{
					"type": "keyword",
				},
				"name": map[string]interface{}{
					"type": "text",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"coordinates": map[string]interface{}{
					"type": "geo_point",
				},
			},
		},
		configuration.Settings{
			"number_of_shards":   "5",
			"max_result_window":  "100000",
			"number_of_replicas": "1",
		},
	)
}

func getConfiguration2() configuration.Index {
	return configuration.New(
		configuration.Mappings{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "integer",
				},
				"uuid": map[string]interface{}{
					"type": "keyword",
				},
				"name": map[string]interface{}{
					"type": "text",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
				"coordinates": map[string]interface{}{
					"type": "geo_point",
				},
			},
		},
		configuration.Settings{
			"number_of_shards":   "5",
			"max_result_window":  "100000",
			"number_of_replicas": "1",
		},
	)
}

func TestCompare_Compare_AliasDoesNotExist(t *testing.T) {
	client := elasticsearch.NewMockClient()

	compareAction := action.NewCompare(
		client,
		prefix,
		true,
	)

	client.On(
		"AliasExist",
		elasticsearch.ResolveAliasName(prefix, indexName1),
	).Return(
		false,
		nil,
	)

	compareResult, err := compareAction.Compare(indexName1, getConfiguration1())

	assert.NoError(t, err)

	assert.Equal(
		t,
		action.CompareResult{
			AliasName:        elasticsearch.ResolveAliasName(prefix, indexName1),
			CurrentIndexName: "",
			CurrentConfig:    configuration.Index{},
			NewConfig:        getConfiguration1(),
			Result:           strategy.NewIndexVoterResult(strategy.IndexDecisionCreate, nil),
		},
		compareResult,
	)
}

func TestCompare_CompareAll(t *testing.T) {
	client := elasticsearch.NewMockClient()

	compareAction := action.NewCompare(
		client,
		prefix,
		true,
	)

	client.On(
		"AliasExist",
		elasticsearch.ResolveAliasName(prefix, indexName1),
	).Return(
		false,
		nil,
	)

	client.On(
		"AliasExist",
		elasticsearch.ResolveAliasName(prefix, indexName2),
	).Return(
		false,
		nil,
	)

	compareResult, err := compareAction.CompareAll(
		configuration.IndexCollection{
			indexName1: getConfiguration1(),
			indexName2: getConfiguration2(),
		},
	)

	assert.NoError(t, err)

	assert.ElementsMatch(
		t,
		action.CompareResultCollection{
			action.CompareResult{
				AliasName:        elasticsearch.ResolveAliasName(prefix, indexName1),
				CurrentIndexName: "",
				CurrentConfig:    configuration.Index{},
				NewConfig:        getConfiguration1(),
				Result:           strategy.NewIndexVoterResult(strategy.IndexDecisionCreate, nil),
			},
			action.CompareResult{
				AliasName:        elasticsearch.ResolveAliasName(prefix, indexName2),
				CurrentIndexName: "",
				CurrentConfig:    configuration.Index{},
				NewConfig:        getConfiguration2(),
				Result:           strategy.NewIndexVoterResult(strategy.IndexDecisionCreate, nil),
			},
		},
		compareResult,
	)
}

func TestCompare_Compare_AliasExist(t *testing.T) {
	client := elasticsearch.NewMockClient()

	compareAction := action.NewCompare(
		client,
		prefix,
		true,
	)

	aliasName := elasticsearch.ResolveAliasName(prefix, indexName1)

	client.On(
		"AliasExist",
		aliasName,
	).Return(
		true,
		nil,
	)

	aliasedIndexName := "aliased-index"

	client.On(
		"GetAliasedIndex",
		aliasName,
	).Return(
		aliasedIndexName,
		nil,
	)

	client.On(
		"GetIndexConfiguration",
		aliasedIndexName,
	).Return(
		getConfiguration1(),
		nil,
	)

	compareResult, err := compareAction.Compare(indexName1, getConfiguration1())

	assert.NoError(t, err)

	assert.Equal(
		t,
		action.CompareResult{
			AliasName:        elasticsearch.ResolveAliasName(prefix, indexName1),
			CurrentIndexName: aliasedIndexName,
			CurrentConfig:    getConfiguration1(),
			NewConfig:        getConfiguration1(),
			Result:           strategy.NewIndexVoterResult(strategy.IndexDecisionNone, nil),
		},
		compareResult,
	)
}

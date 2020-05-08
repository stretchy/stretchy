package action_test

import (
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchy/stretchy/pkg/action"
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/elasticsearch"
	"github.com/stretchy/stretchy/pkg/strategy"
)

const noneAliasName = "index-none"

const createAliasName = "index-create"

const updateAliasName = "index-update"
const currentUpdateIndexName = "index-update-current"

const migrateAliasName = "index-migrate"
const currentMigrateIndexName = "index-migrate-current"

func TestApply_ApplyAll(t *testing.T) {
	client := elasticsearch.NewMockClient()
	now := time.Now()

	patch := monkey.Patch(time.Now, func() time.Time { return now })
	defer patch.Unpatch()

	// Create Mocks
	client.On(
		"CreateIndex",
		elasticsearch.CreateIndexName(createAliasName),
		createConfig(),
	).Return(nil)

	client.On(
		"CreateAlias",
		createAliasName,
		elasticsearch.CreateIndexName(createAliasName),
	).Return(nil)

	// Update Mocks
	client.On(
		"UpdateIndexConfiguration",
		currentUpdateIndexName,
		updateConfig(),
	).Return(nil)

	// Migrate Mocks

	client.On(
		"CreateIndex",
		elasticsearch.CreateIndexName(migrateAliasName),
		migrateConfig(),
	).Return(nil)

	client.On(
		"Reindex",
		currentMigrateIndexName,
		elasticsearch.CreateIndexName(migrateAliasName),
	).Return(nil)

	client.On(
		"UpdateAlias",
		migrateAliasName,
		elasticsearch.CreateIndexName(migrateAliasName),
	).Return(nil)

	applyAction := action.NewApply(client)

	compareResultCollection := action.CompareResultCollection{
		action.CompareResult{
			AliasName: noneAliasName,
			Result:    strategy.NewIndexVoterResult(strategy.IndexDecisionNone, nil),
		}, // NONE
		action.CompareResult{
			AliasName: createAliasName,
			NewConfig: createConfig(),
			Result:    strategy.NewIndexVoterResult(strategy.IndexDecisionCreate, nil),
		}, // CREATE
		action.CompareResult{
			AliasName:        updateAliasName,
			NewConfig:        updateConfig(),
			CurrentIndexName: currentUpdateIndexName,
			Result:           strategy.NewIndexVoterResult(strategy.IndexDecisionUpdate, nil),
		}, // UPDATE
		action.CompareResult{
			AliasName:        migrateAliasName,
			NewConfig:        migrateConfig(),
			CurrentIndexName: currentMigrateIndexName,
			Result:           strategy.NewIndexVoterResult(strategy.IndexDecisionMigrate, nil),
		}, // MIGRATE
	}

	err := applyAction.ApplyAll(compareResultCollection)
	assert.NoError(t, err)

	mock.AssertExpectationsForObjects(t, client)
}

func createConfig() configuration.Index {
	return configuration.New(
		configuration.Mappings{
			"type": "create-index",
		},
		configuration.Settings{
			"type": "create-index",
		},
	)
}

func updateConfig() configuration.Index {
	return configuration.New(
		configuration.Mappings{
			"type": "update-index",
		},
		configuration.Settings{
			"type": "update-index",
		},
	)
}

func migrateConfig() configuration.Index {
	return configuration.New(
		configuration.Mappings{
			"type": "migrate-index",
		},
		configuration.Settings{
			"type": "migrate-index",
		},
	)
}

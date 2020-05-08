package action

import (
	"fmt"

	"github.com/stretchy/stretchy/pkg/elasticsearch"
	"github.com/stretchy/stretchy/pkg/strategy"
)

type Apply struct {
	client elasticsearch.Client
}

func NewApply(
	client elasticsearch.Client,
) *Apply {
	return &Apply{
		client: client,
	}
}

func (a *Apply) Apply(compareResult CompareResult) error {
	switch compareResult.Result.Action() {
	case strategy.IndexDecisionNone:
		return nil
	case strategy.IndexDecisionCreate:
		newIndexName := elasticsearch.CreateIndexName(compareResult.AliasName)
		if err := a.client.CreateIndex(newIndexName, compareResult.NewConfig); err != nil {
			return err
		}

		return a.client.CreateAlias(compareResult.AliasName, newIndexName)
	case strategy.IndexDecisionUpdate:
		return a.client.UpdateIndexConfiguration(compareResult.CurrentIndexName, compareResult.NewConfig)
	case strategy.IndexDecisionMigrate:
		newIndexName := elasticsearch.CreateIndexName(compareResult.AliasName)
		if err := a.client.CreateIndex(newIndexName, compareResult.NewConfig); err != nil {
			return err
		}

		if err := a.client.Reindex(compareResult.CurrentIndexName, newIndexName); err != nil {
			return err
		}

		return a.client.UpdateAlias(compareResult.AliasName, newIndexName)
	}

	return fmt.Errorf(
		"unknown decision '%s' on index '%s'",
		compareResult.Result.Action().String(),
		compareResult.AliasName,
	)
}

func (a *Apply) ApplyAll(compareResultCollection CompareResultCollection) error {
	for _, compareResult := range compareResultCollection {
		if err := a.Apply(compareResult); err != nil {
			return err
		}
	}

	return nil
}

package action

import (
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/elasticsearch"
	"github.com/stretchy/stretchy/pkg/strategy"
)

type Compare struct {
	client           elasticsearch.Client
	indexPrefix      string
	indexActionVoter *strategy.IndexActionVoter
}

type CompareResult struct {
	AliasName        string
	CurrentIndexName string
	CurrentConfig    configuration.Index
	NewConfig        configuration.Index
	Result           strategy.IndexVoterResult
}

type CompareResultCollection []CompareResult

func NewCompare(
	client elasticsearch.Client,
	indexPrefix string,
	updateEnabled bool,
) *Compare {
	return &Compare{
		client:           client,
		indexPrefix:      indexPrefix,
		indexActionVoter: strategy.NewIndexActionVoter(updateEnabled),
	}
}

func (c *Compare) Compare(indexName string, index configuration.Index) (CompareResult, error) {
	aliasName := elasticsearch.ResolveAliasName(c.indexPrefix, indexName)

	aliasExist, err := c.client.AliasExist(aliasName)
	if err != nil {
		return CompareResult{}, err
	}

	if !aliasExist {
		action, err := c.indexActionVoter.Compare(nil, &index)
		if err != nil {
			return CompareResult{}, err
		}

		return CompareResult{
			AliasName:        aliasName,
			CurrentIndexName: "",
			CurrentConfig:    configuration.Index{},
			NewConfig:        index,
			Result:           action,
		}, nil
	}

	currentIndexName, err := c.client.GetAliasedIndex(aliasName)
	if err != nil {
		return CompareResult{}, err
	}

	currentIndex, err := c.client.GetIndexConfiguration(currentIndexName)
	if err != nil {
		return CompareResult{}, err
	}

	action, err := c.indexActionVoter.Compare(&currentIndex, &index)
	if err != nil {
		return CompareResult{}, err
	}

	return CompareResult{
		AliasName:        aliasName,
		CurrentIndexName: currentIndexName,
		CurrentConfig:    currentIndex,
		NewConfig:        index,
		Result:           action,
	}, nil
}

func (c *Compare) CompareAll(indexCollection configuration.IndexCollection) (CompareResultCollection, error) {
	compareResultCollection := CompareResultCollection{}

	for indexName, index := range indexCollection {
		compareResult, err := c.Compare(indexName, index)
		if err != nil {
			return nil, err
		}

		compareResultCollection = append(compareResultCollection, compareResult)
	}

	return compareResultCollection, nil
}

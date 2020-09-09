package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/stretchy/stretchy/pkg/configuration"
)

// V7Client is an Elasticsearch client for elasticsearch v6. Create one by calling NewV7Client.
type V7Client struct {
	client *elastic.Client
}

// NewV7Client creates a V6Client instance.
func NewV7Client(
	options Options,
) (*V7Client, error) {
	client, err := newOlivereV7(options)
	if err != nil {
		return nil, err
	}

	return &V7Client{
		client: client,
	}, nil
}

func newOlivereV7(
	options Options,
) (*elastic.Client, error) {
	buildOptions := []elastic.ClientOptionFunc{
		elastic.SetURL(options.Host),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	}

	if options.User != "" && options.Password != "" {
		buildOptions = append(buildOptions, elastic.SetBasicAuth(options.User, options.Password))
	}

	if options.Debug == true {
		buildOptions = append(buildOptions, elastic.SetTraceLog(newLogger()))
	}

	return elastic.NewClient(buildOptions...)
}

func (c *V7Client) IndexExist(indexName string) (bool, error) {
	exist, err := c.client.
		IndexExists(indexName).
		Do(context.Background())

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (c *V7Client) AliasExist(aliasName string) (bool, error) {
	_, err := c.client.Aliases().Alias(aliasName).Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *V7Client) CreateIndex(indexName string, mapping configuration.Index) error {
	_, err := c.client.CreateIndex(indexName).BodyJson(mapping).Do(context.Background())
	return err
}

func (c *V7Client) CreateAlias(aliasName string, indexName string) error {
	return c.UpdateAlias(aliasName, indexName)
}

func (c *V7Client) UpdateAlias(aliasName string, indexName string) error {
	_, err := c.client.
		Alias().
		Remove("*", aliasName).
		Add(indexName, aliasName).
		Do(context.Background())

	return err
}

func (c *V7Client) GetAliasedIndex(aliasName string) (string, error) {
	aliasResult, err := c.client.
		Aliases().
		Alias(aliasName).
		Do(context.Background())

	if err != nil {
		return "", err
	}

	if len(aliasResult.Indices) > 1 {
		return "", fmt.Errorf("alias '%s' targets more than 1 index. currently not supported", aliasName)
	}

	for index := range aliasResult.Indices {
		return index, nil
	}

	return "", fmt.Errorf("alias '%s' doesn't target any index", aliasName)
}

func (c *V7Client) GetIndexConfiguration(indexName string) (configuration.Index, error) {
	indexResult, err := c.client.
		IndexGet(indexName).
		Do(context.Background())

	if err != nil {
		return configuration.Index{}, err
	}

	return configuration.New(
		indexResult[indexName].Mappings,
		indexResult[indexName].Settings,
	), nil
}

func (c *V7Client) Reindex(sourceIndexName string, targetIndexName string) error {
	_, err := c.client.
		Reindex().
		SourceIndex(sourceIndexName).
		DestinationIndex(targetIndexName).
		WaitForCompletion(true).
		Refresh("true").
		Do(context.Background())

	return err
}

func (c *V7Client) UpdateIndexConfiguration(indexName string, configuration configuration.Index) error {
	if _, err := c.client.
		PutMapping().
		Index(indexName).
		BodyJson(configuration.GetMappings()).
		Do(context.Background()); err != nil {
		return err
	}

	return nil
}

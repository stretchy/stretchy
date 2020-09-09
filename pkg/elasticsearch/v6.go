package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/stretchy/stretchy/pkg/configuration"
)

// V6Client is an Elasticsearch client for elasticsearch v6. Create one by calling NewV6Client.
type V6Client struct {
	client *elastic.Client
}

// NewV6Client creates a V6Client instance.
func NewV6Client(
	options Options,
) (*V6Client, error) {
	client, err := newOlivereV6(options)
	if err != nil {
		return nil, err
	}

	return &V6Client{
		client: client,
	}, nil
}

func newOlivereV6(
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

func (c *V6Client) IndexExist(indexName string) (bool, error) {
	exist, err := c.client.
		IndexExists(indexName).
		Do(context.Background())

	if err != nil {
		return false, err
	}

	return exist, nil
}

func (c *V6Client) AliasExist(aliasName string) (bool, error) {
	_, err := c.client.Aliases().Alias(aliasName).Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c *V6Client) CreateIndex(indexName string, mapping configuration.Index) error {
	_, err := c.client.CreateIndex(indexName).
		BodyJson(mapping).
		IncludeTypeName(false).
		Do(context.Background())

	return err
}

func (c *V6Client) CreateAlias(aliasName string, indexName string) error {
	return c.UpdateAlias(aliasName, indexName)
}

func (c *V6Client) UpdateAlias(aliasName string, indexName string) error {
	_, err := c.client.
		Alias().
		Remove("*", aliasName).
		Add(indexName, aliasName).
		Do(context.Background())

	return err
}

func (c *V6Client) GetAliasedIndex(aliasName string) (string, error) {
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

func (c *V6Client) GetIndexConfiguration(indexName string) (configuration.Index, error) {
	indexResult, err := c.client.
		IndexGet(indexName).
		IncludeTypeName(false).
		Do(context.Background())

	if err != nil {
		return configuration.Index{}, err
	}

	return configuration.New(
		indexResult[indexName].Mappings,
		indexResult[indexName].Settings,
	), nil
}

func (c *V6Client) Reindex(sourceIndexName string, targetIndexName string) error {
	return c.reindex(sourceIndexName, targetIndexName, true)
}

func (c *V6Client) reindex(sourceIndexName string, targetIndexName string, shouldRetry bool) error {
	_, err := c.client.
		Reindex().
		SourceIndex(sourceIndexName).
		DestinationIndexAndType(targetIndexName, "_doc").
		WaitForCompletion(true).
		Refresh("true").
		Do(context.Background())

	if err != nil && shouldRetry == true {
		return c.reindex(sourceIndexName, targetIndexName, false)
	}

	return err
}

func (c *V6Client) UpdateIndexConfiguration(indexName string, configuration configuration.Index) error {
	if _, err := c.client.
		PutMapping().
		Index(indexName).
		BodyJson(configuration.GetMappings()).
		IncludeTypeName(false).
		Do(context.Background()); err != nil {
		return err
	}

	return nil
}

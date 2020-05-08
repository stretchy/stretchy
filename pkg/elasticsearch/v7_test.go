//nolint:dupl
package elasticsearch

import (
	"context"
	"github.com/olivere/elastic/v7"
	"reflect"
)

type V7ClientExtended V7Client

func NewV7ClientExtended(client *V7Client) *V7ClientExtended {
	c := V7ClientExtended(*client)

	return &c
}

func (v7e *V7ClientExtended) GetClient() Client {
	return (*V7Client)(v7e)
}

func (v7e *V7ClientExtended) Load(indexName string, documents ...map[string]interface{}) error {
	requests := make([]elastic.BulkableRequest, len(documents))

	for i, d := range documents {
		requests[i] = elastic.NewBulkIndexRequest().
			Index(indexName).
			Doc(d).
			RetryOnConflict(7)
	}

	response, err := v7e.client.
		Bulk().
		Add(requests...).
		Refresh("true").
		Do(context.Background())

	if err != nil {
		return err
	}

	if response.Errors == false {
		return nil
	}

	for _, items := range response.Items {
		for _, item := range items {
			if item.Error != nil {
				return &elastic.Error{
					Status:  item.Status,
					Details: item.Error,
				}
			}
		}
	}

	return nil
}

func (v7e *V7ClientExtended) CleanupIndex(indexName string) error {
	_, err := v7e.client.
		DeleteByQuery(indexName).
		ProceedOnVersionConflict().
		Query(elastic.NewMatchAllQuery()).
		Refresh("true").
		Do(context.Background())

	return err
}

func (v7e *V7ClientExtended) Cleanup() error {
	_, err := v7e.client.DeleteIndex("*").Do(context.Background())

	return err
}

func (v7e *V7ClientExtended) GetAll(indexName string) ([]map[string]interface{}, error) {
	searchResult, err := v7e.client.
		Search(indexName).
		Query(elastic.NewMatchAllQuery()).
		Size(1000).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	return v7e.parseSearchResults(searchResult), nil
}

func (v7e *V7ClientExtended) parseSearchResults(searchResult *elastic.SearchResult) []map[string]interface{} {
	results := make([]map[string]interface{}, len(searchResult.Hits.Hits))
	for i, item := range searchResult.Each(reflect.TypeOf(map[string]interface{}{})) {
		results[i] = item.(map[string]interface{})
	}

	return results
}

//nolint:dupl
package elasticsearch

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

func Test_newLoggerV6(t *testing.T) {
	logger := newLogger()

	assert.IsType(t, &log.Logger{}, logger)
}

type V6ClientExtended V6Client

func NewV6ClientExtended(client *V6Client) *V6ClientExtended {
	c := V6ClientExtended(*client)

	return &c
}

func (v6e *V6ClientExtended) GetClient() Client {
	return (*V6Client)(v6e)
}

func (v6e *V6ClientExtended) Load(indexName string, documents ...map[string]interface{}) error {
	requests := make([]elastic.BulkableRequest, len(documents))

	for i, d := range documents {
		requests[i] = elastic.NewBulkIndexRequest().
			Index(indexName).
			Doc(d).
			Type("_doc").
			RetryOnConflict(7)
	}

	response, err := v6e.client.
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

func (v6e *V6ClientExtended) CleanupIndex(indexName string) error {
	_, err := v6e.client.
		DeleteByQuery(indexName).
		ProceedOnVersionConflict().
		Query(elastic.NewMatchAllQuery()).
		Refresh("true").
		Do(context.Background())

	return err
}

func (v6e *V6ClientExtended) Cleanup() error {
	_, err := v6e.client.DeleteIndex("*").Do(context.Background())

	return err
}

func (v6e *V6ClientExtended) GetAll(indexName string) ([]map[string]interface{}, error) {
	searchResult, err := v6e.client.
		Search(indexName).
		Query(elastic.NewMatchAllQuery()).
		Size(1000).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	return v6e.parseSearchResults(searchResult), nil
}

func (v6e *V6ClientExtended) parseSearchResults(searchResult *elastic.SearchResult) []map[string]interface{} {
	results := make([]map[string]interface{}, len(searchResult.Hits.Hits))
	for i, item := range searchResult.Each(reflect.TypeOf(map[string]interface{}{})) {
		results[i] = item.(map[string]interface{})
	}

	return results
}

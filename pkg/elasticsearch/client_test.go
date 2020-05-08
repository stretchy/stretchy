// +build integration

package elasticsearch_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
	"github.com/stretchy/stretchy/pkg/elasticsearch"
	"github.com/stretchy/stretchy/pkg/utils"
)

func TestElasticSearch_NewV6(t *testing.T) {
	elasticsearchHost := getElasticSearchHost(t, 6)
	assert.NotEmpty(t, elasticsearchHost)
	client, err := elasticsearch.New(elasticsearch.Options{
		Host:     elasticsearchHost,
		User:     "",
		Password: "",
		Debug:    false,
	})

	assert.NoError(t, err)
	assert.IsType(t, &elasticsearch.V6Client{}, client)
}

func TestElasticSearch_NewV7(t *testing.T) {
	elasticsearchHost := getElasticSearchHost(t, 7)
	assert.NotEmpty(t, elasticsearchHost)
	client, err := elasticsearch.New(elasticsearch.Options{
		Host:     elasticsearchHost,
		User:     "",
		Password: "",
		Debug:    false,
	})

	assert.NoError(t, err)
	assert.IsType(t, &elasticsearch.V7Client{}, client)
}

func TestElasticSearch_NewErrors(t *testing.T) {
	testCases := []struct {
		name string
		host string
	}{
		{
			name: "unknown version",
			host: "http://www.mocky.io/v2/5e9057e0330000741327d6d4",
		},
		{
			name: "missing scheme",
			host: "www.myhost.com",
		},
		{
			name: "missing host",
			host: "http://",
		},
		{
			name: "empty host",
			host: "",
		},
		{
			name: "404 error",
			host: "http://www.mocky.io/v2/5ea2971a310000dad61ef208",
		},
		{
			name: "500 error",
			host: "http://www.mocky.io/v2/5ea299943100008ccd1ef216",
		},
		{
			name: "domain does not exist",
			host: "http://www.notexist.tech/",
		},
		{
			name: "json syntax errors",
			host: "http://www.mocky.io/v2/5ea29a103100008ccd1ef21b",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			client, err := elasticsearch.New(elasticsearch.Options{
				Host:     testCase.host,
				User:     "",
				Password: "",
				Debug:    false,
			})

			assert.Nil(t, client)
			assert.Error(t, err)
		})
	}
}

const (
	existingIndexName    = "test-index-a"
	notExistingIndex     = "test-index-not-existing"
	existingAliasName    = "alias-test-index-a"
	notExistingAliasName = "alias-test-index-not-existing"
)

func createBaseDocument() map[string]interface{} {
	return map[string]interface{}{
		"id":          rand.Intn(2147483647),
		"uuid":        utils.RandomString(32),
		"name":        utils.RandomString(100),
		"created_at":  time.Now().Format("2006-01-02T15:04:05-07:00"),
		"coordinates": []float64{rand.Float64(), rand.Float64()},
	}
}

func TestClient_IndexExist(t *testing.T) {
	for _, clientTestCase := range getClientTestCases(t) {
		clientTestCase := clientTestCase
		t.Run(clientTestCase.name, func(t *testing.T) {
			loadTestScenario(t, clientTestCase.extendedClient)

			exist, err := clientTestCase.client.IndexExist(existingIndexName)
			assert.NoError(t, err)
			assert.True(t, exist)

			exist, err = clientTestCase.client.IndexExist(notExistingIndex)
			assert.NoError(t, err)
			assert.False(t, exist)
		})
	}
}

func TestClient_AliasExist(t *testing.T) {
	for _, clientTestCase := range getClientTestCases(t) {
		clientTestCase := clientTestCase
		t.Run(clientTestCase.name, func(t *testing.T) {
			loadTestScenario(t, clientTestCase.extendedClient)

			exist, err := clientTestCase.client.AliasExist(existingAliasName)
			assert.NoError(t, err)
			assert.True(t, exist)

			exist, err = clientTestCase.client.AliasExist(notExistingAliasName)
			assert.NoError(t, err)
			assert.False(t, exist)
		})
	}
}

func TestClient_GetAliasedIndex(t *testing.T) {
	for _, clientTestCase := range getClientTestCases(t) {
		clientTestCase := clientTestCase
		t.Run(clientTestCase.name, func(t *testing.T) {
			loadTestScenario(t, clientTestCase.extendedClient)

			indexName, err := clientTestCase.client.GetAliasedIndex(existingAliasName)
			assert.NoError(t, err)
			assert.Equal(t, existingIndexName, indexName)

			indexName, err = clientTestCase.client.GetAliasedIndex(notExistingAliasName)
			assert.Error(t, err)
			assert.Equal(t, "", indexName)
		})
	}
}

func TestClient_GetIndexConfiguration(t *testing.T) {
	for _, clientTestCase := range getClientTestCases(t) {
		clientTestCase := clientTestCase
		t.Run(clientTestCase.name, func(t *testing.T) {
			loadTestScenario(t, clientTestCase.extendedClient)

			indexConfiguration, err := clientTestCase.client.GetIndexConfiguration(existingIndexName)
			assert.NoError(t, err)
			assert.Equal(t, getBaseConfiguration(t), indexConfiguration)

			indexConfiguration, err = clientTestCase.client.GetIndexConfiguration(notExistingIndex)
			assert.Error(t, err)
			assert.Equal(t, configuration.Index{}, indexConfiguration)
		})
	}
}

func TestClient_UpdateIndexConfiguration(t *testing.T) {
	for _, clientTestCase := range getClientTestCases(t) {
		clientTestCase := clientTestCase
		t.Run(clientTestCase.name, func(t *testing.T) {
			loadTestScenarioWithDocuments(t, clientTestCase.extendedClient)

			indexConfiguration, err := clientTestCase.client.GetIndexConfiguration(existingIndexName)
			assert.NoError(t, err)
			assert.Equal(t, getBaseConfiguration(t), indexConfiguration)

			err = clientTestCase.client.UpdateIndexConfiguration(existingIndexName, getExtendedBaseConfiguration(t))
			assert.NoError(t, err)

			indexConfiguration, err = clientTestCase.client.GetIndexConfiguration(existingIndexName)
			assert.NoError(t, err)
			assert.Equal(t, getExtendedBaseConfiguration(t), indexConfiguration)

			err = clientTestCase.client.UpdateIndexConfiguration(notExistingIndex, getExtendedBaseConfiguration(t))
			assert.Error(t, err)
		})
	}
}

func TestClient_Reindex(t *testing.T) {
	for _, clientTestCase := range getClientTestCases(t) {
		clientTestCase := clientTestCase
		t.Run(clientTestCase.name, func(t *testing.T) {
			loadTestScenarioWithDocuments(t, clientTestCase.extendedClient)

			indexConfiguration, err := clientTestCase.client.GetIndexConfiguration(existingIndexName)
			assert.NoError(t, err)
			assert.Equal(t, getBaseConfiguration(t), indexConfiguration)

			newIndexName := "new-index-test"
			err = clientTestCase.client.CreateIndex(newIndexName, getExtendedBaseConfiguration(t))
			assert.NoError(t, err)

			indexConfiguration, err = clientTestCase.client.GetIndexConfiguration(newIndexName)
			assert.NoError(t, err)
			assert.Equal(t, getExtendedBaseConfiguration(t), indexConfiguration)

			err = clientTestCase.client.Reindex(existingIndexName, newIndexName)
			assert.NoError(t, err)

			assertSameDocuments(t, clientTestCase.extendedClient, existingIndexName, newIndexName)
		})
	}
}

func assertSameDocuments(t *testing.T, client extendedClient, index1 string, index2 string) {
	index1Documents, err := client.GetAll(index1)
	assert.NoError(t, err)

	index2Documents, err := client.GetAll(index2)
	assert.NoError(t, err)

	assert.Len(t, index2Documents, len(index1Documents))

	assert.Equal(t, index1Documents, index2Documents)
}

type extendedClient interface {
	Load(indexName string, documents ...map[string]interface{}) error
	CleanupIndex(indexName string) error
	Cleanup() error
	GetAll(indexName string) ([]map[string]interface{}, error)
	GetClient() elasticsearch.Client
}

type ClientTestCase struct {
	name           string
	client         elasticsearch.Client
	extendedClient extendedClient
}

func getElasticSearchHost(t *testing.T, version int) string {
	switch version {
	case 6:
		return os.Getenv("TEST_ELASTICSEARCH_HOST_v6")
	case 7:
		return os.Getenv("TEST_ELASTICSEARCH_HOST_v7")
	}

	t.Fatalf("Missing elasticsearch version '%d' test configuration", version)

	return ""
}

func getBaseConfiguration(t *testing.T) configuration.Index {
	return configuration.New(
		map[string]interface{}{
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
		map[string]interface{}{
			"number_of_shards":   "5",
			"max_result_window":  "100000",
			"number_of_replicas": "1",
		},
	)
}

func getExtendedBaseConfiguration(t *testing.T) configuration.Index {
	newConfiguration := getBaseConfiguration(t)

	err := newConfiguration.GetMappings().
		Merge(
			map[string]interface{}{
				"properties": map[string]interface{}{
					"new_field": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
		)

	assert.NoError(t, err)

	return newConfiguration
}

func getClientTestCases(t *testing.T) []ClientTestCase {
	clientV6, err := elasticsearch.NewV6Client(elasticsearch.Options{
		Host:     getElasticSearchHost(t, 6),
		User:     "",
		Password: "",
		Debug:    false,
	})

	assert.NoError(t, err)

	clientV7, err := elasticsearch.NewV7Client(elasticsearch.Options{
		Host:     getElasticSearchHost(t, 7),
		User:     "",
		Password: "",
		Debug:    false,
	})

	assert.NoError(t, err)

	return []ClientTestCase{
		{
			name:           "v6",
			client:         clientV6,
			extendedClient: elasticsearch.NewV6ClientExtended(clientV6),
		},
		{
			name:           "v7",
			client:         clientV7,
			extendedClient: elasticsearch.NewV7ClientExtended(clientV7),
		},
	}
}

const docsNumber = 100

func loadTestScenario(t *testing.T, extendedClient extendedClient) {
	err := extendedClient.Cleanup()
	assert.NoError(t, err)

	err = extendedClient.GetClient().CreateIndex(existingIndexName, getBaseConfiguration(t))
	assert.NoError(t, err)

	err = extendedClient.GetClient().CreateAlias(existingAliasName, existingIndexName)
	assert.NoError(t, err)
}

func loadTestScenarioWithDocuments(t *testing.T, extendedClient extendedClient) {
	loadTestScenario(t, extendedClient)
	docs := make([]map[string]interface{}, docsNumber)

	for n := 0; n < docsNumber; n++ {
		docs[n] = createBaseDocument()
	}

	err := extendedClient.Load(existingIndexName, docs...)
	assert.NoError(t, err)
}

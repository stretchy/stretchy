package elasticsearch

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchy/stretchy/pkg/configuration"
)

type MockClient struct {
	mock.Mock
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (mc *MockClient) IndexExist(indexName string) (bool, error) {
	args := mc.Called(indexName)
	return args.Bool(0), args.Error(1)
}

func (mc *MockClient) AliasExist(aliasName string) (bool, error) {
	args := mc.Called(aliasName)
	return args.Bool(0), args.Error(1)
}

func (mc *MockClient) CreateIndex(indexName string, mapping configuration.Index) error {
	args := mc.Called(indexName, mapping)
	return args.Error(0)
}

func (mc *MockClient) CreateAlias(aliasName string, indexName string) error {
	args := mc.Called(aliasName, indexName)
	return args.Error(0)
}

func (mc *MockClient) UpdateAlias(aliasName string, newIndexName string) error {
	args := mc.Called(aliasName, newIndexName)
	return args.Error(0)
}

func (mc *MockClient) GetAliasedIndex(aliasName string) (string, error) {
	args := mc.Called(aliasName)
	return args.String(0), args.Error(1)
}

func (mc *MockClient) GetIndexConfiguration(indexName string) (configuration.Index, error) {
	args := mc.Called(indexName)
	return args.Get(0).(configuration.Index), args.Error(1)
}

func (mc *MockClient) UpdateIndexConfiguration(indexName string, configuration configuration.Index) error {
	args := mc.Called(indexName, configuration)
	return args.Error(0)
}

func (mc *MockClient) Reindex(sourceIndexName string, targetIndexName string) error {
	args := mc.Called(sourceIndexName, targetIndexName)
	return args.Error(0)
}

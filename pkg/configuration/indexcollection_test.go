package configuration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchy/stretchy/pkg/configuration"
)

func assertIndexExampleExist(t *testing.T, indexCollection configuration.IndexCollection, index configuration.Index) {
	assert.True(t, indexCollection.Exist(indexExampleName))

	storedIndex, err := indexCollection.Get(indexExampleName)
	assert.NoError(t, err)
	assert.Equal(t, index, storedIndex)
}

func TestIndexCollection_Load(t *testing.T) {
	indexCollection := configuration.NewIndexCollection()
	assert.Len(t, indexCollection, 0)

	indexExample := getIndexExample()

	_, exist := indexCollection[indexExampleName]
	assert.False(t, exist)

	indexCollection.Load(indexExampleName, indexExample)

	assertIndexExampleExist(t, indexCollection, indexExample)
}

func TestIndexCollection_Load_Overwrite(t *testing.T) {
	indexCollection := configuration.NewIndexCollection()
	assert.Len(t, indexCollection, 0)

	indexExample := getIndexExample()

	_, exist := indexCollection[indexExampleName]
	assert.False(t, exist)

	indexCollection.Load(indexExampleName, indexExample)
	assertIndexExampleExist(t, indexCollection, indexExample)

	newIndex := configuration.Index{}

	indexCollection.Load(indexExampleName, newIndex)

	assertIndexExampleExist(t, indexCollection, newIndex)
}

func TestIndexCollection_Get(t *testing.T) {
	indexCollection := configuration.NewIndexCollection()
	assert.Len(t, indexCollection, 0)

	indexExample := getIndexExample()

	indexCollection.Load(indexExampleName, indexExample)

	assertIndexExampleExist(t, indexCollection, indexExample)

	index, err := indexCollection.Get(indexExampleName)
	assert.NoError(t, err)
	assert.Equal(t, indexExample, index)
}

func TestIndexCollection_Get_NotExist(t *testing.T) {
	indexCollection := configuration.NewIndexCollection()

	index, err := indexCollection.Get(indexExampleName)
	assert.Error(t, err)
	assert.Equal(t, configuration.Index{}, index)
}

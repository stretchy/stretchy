package configuration

import "fmt"

type IndexCollection map[string]Index

func NewIndexCollection() IndexCollection {
	return IndexCollection{}
}

func (mc IndexCollection) Get(indexName string) (Index, error) {
	if mc.Exist(indexName) {
		return mc[indexName], nil
	}

	return Index{}, fmt.Errorf("cannot find configuration for index '%s'", indexName)
}

func (mc IndexCollection) Load(indexName string, index Index) {
	mc[indexName] = index
}

func (mc IndexCollection) Exist(indexName string) bool {
	_, exist := mc[indexName]

	return exist
}

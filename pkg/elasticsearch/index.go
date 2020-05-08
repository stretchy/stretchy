package elasticsearch

import (
	"fmt"
	"time"
)

func CreateIndexName(aliasName string) string {
	return fmt.Sprintf("%s-%d", aliasName, time.Now().Unix())
}

func ResolveAliasName(prefix string, mappingName string) string {
	if prefix == "" {
		return mappingName
	}

	return fmt.Sprintf("%s-%s", prefix, mappingName)
}

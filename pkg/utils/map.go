package utils

import (
	"fmt"
	"reflect"
)

func MapHasKey(key interface{}, haystack interface{}) (bool, error) {
	kind := reflect.TypeOf(haystack).Kind()
	if kind != reflect.Map {
		return false, fmt.Errorf("expected '%s', got '%s'", reflect.Map.String(), kind.String())
	}

	keys := reflect.ValueOf(haystack).MapKeys()

	for _, k := range keys {
		if k.Interface() == key {
			return true, nil
		}
	}

	return false, nil
}

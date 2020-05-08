package utils

import (
	"fmt"
	"reflect"
)

func InSlice(element interface{}, slice interface{}) (bool, error) {
	kind := reflect.TypeOf(slice).Kind()
	if kind != reflect.Slice {
		return false, fmt.Errorf("expected '%s', got '%s'", reflect.Slice.String(), kind.String())
	}

	s := reflect.ValueOf(slice)
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == element {
			return true, nil
		}
	}

	return false, nil
}

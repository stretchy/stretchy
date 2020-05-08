package utils

import (
	"fmt"
	"reflect"
)

func IsAScalar(value interface{}) bool {
	kind := reflect.ValueOf(value).Kind()

	isAScalar, _ := InSlice(kind, []reflect.Kind{
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String,
	})

	return isAScalar
}

func PrintScalar(value interface{}) string {
	if !IsAScalar(value) {
		return ""
	}

	return fmt.Sprint(value)
}

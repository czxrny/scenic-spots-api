package generics

import (
	"reflect"
	"strings"
)

func DereferenceAll[T any](in []*T) []T {
	out := make([]T, 0, len(in))
	for _, ptr := range in {
		if ptr != nil {
			out = append(out, *ptr)
		}
	}
	return out
}

func StructToMapLower[T any](item T) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	structValue := reflect.ValueOf(item)
	structType := reflect.TypeOf(item)
	numberOfFields := structValue.NumField()

	for i := 0; i < numberOfFields; i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// skip private fields
		if !fieldValue.CanInterface() {
			continue
		}

		lowerCaseName := strings.ToLower(field.Name[:1]) + field.Name[1:]

		result[lowerCaseName] = fieldValue.Interface()
	}

	return result, nil
}

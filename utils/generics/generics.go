package generics

import (
	"reflect"
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

func HasEmptyFields(data interface{}) bool {
	structure := reflect.ValueOf(data)

	if structure.Kind() == reflect.Ptr {
		structure = structure.Elem()
	}

	for i := 0; i < structure.NumField(); i++ {
		field := structure.Field(i)

		if !field.CanInterface() {
			continue
		}

		zero := reflect.Zero(field.Type()).Interface()
		current := field.Interface()

		if reflect.DeepEqual(current, zero) {
			return true
		}
	}

	return false
}

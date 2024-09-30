package util

import "fmt"

func TypeAssert[T any](value interface{}, typeName string) (T, error) {
	var castedValue T
	castedValue, ok := value.(T)
	if !ok {
		return castedValue, fmt.Errorf("%s is not of type %T", typeName, castedValue)
	}
	return castedValue, nil
}
package util

import "fmt"

// MapGetKeyAsType retrieves a key from a map and converts it to a given type
func MapGetKeyAsType[T any](key string, m map[string]interface{}) (T, error) {
	if valueRaw, found := m[key]; found {
		if value, ok := valueRaw.(T); ok {
			return value, nil
		}
	}

	return *new(T), fmt.Errorf("unable to find %s key in %v", key, m)
}

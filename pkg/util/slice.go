// Package util provides some common utility functions
package util

// SliceMap maps a slice of type T to a slice of type U
func SliceMap[T any, U any](input []T, mapFunc func(T) U) []U {
	v := make([]U, len(input))
	for i, item := range input {
		v[i] = mapFunc(item)
	}
	return v
}

// SliceFilter returns a new slice consisting of elements that passed the filter function
func SliceFilter[T any](input []T, filter func(T) bool) []T {
	var filtered []T
	for _, item := range input {
		if filter(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

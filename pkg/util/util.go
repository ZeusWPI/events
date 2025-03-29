// Package util provides some common utility functions
package util

// If is a ternary if check
func If[T any](cond bool, vTrue, vFalse T) T {
	if cond {
		return vTrue
	}

	return vFalse
}

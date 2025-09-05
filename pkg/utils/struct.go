// Package utils provides some general usefull utilities
package utils

import "reflect"

// Merge populates the zero fields of the destination with the values of the patch struct
// It will ignore non exported fields
func Merge[T any](dest *T, patch T) {
	if dest == nil {
		return
	}

	destVal := reflect.ValueOf(dest).Elem()
	patchVal := reflect.ValueOf(patch)

	if destVal.Kind() != reflect.Struct || patchVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < patchVal.NumField(); i++ {
		patchField := patchVal.Field(i)
		if !patchField.IsZero() {
			destVal.Field(i).Set(patchField)
		}
	}
}

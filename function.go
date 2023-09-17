package jnoronhautils

import (
	"reflect"
)

func InArray[T any](arr []T, element T) bool {
	if len(arr) > 0 {
		for _, data := range arr {
			if reflect.DeepEqual(data, element) {
				return true
			}
		}
	}
	return false
}

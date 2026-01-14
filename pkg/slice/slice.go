package slice

import (
	"cmp"
	"slices"
	"strings"
)

func FilterArray[T any](array []T, fun func(T) bool) []T {
	var newArr []T
	for _, element := range array {
		if fun(element) {
			newArr = append(newArr, element)
		}
	}
	return newArr
}

func ArrayToStringBySep(arr []string, sep string) string {
	return strings.Join(arr, sep)
}

func ArrayToString(arr []string) string {
	return ArrayToStringBySep(arr, " ")
}

func RemoveDuplicate[T cmp.Ordered](sliceList []T) []T {
	if len(sliceList) == 0 {
		return nil
	}
	newSliceList := make([]T, len(sliceList))
	copy(newSliceList, sliceList)
	slices.Sort(newSliceList)
	return slices.Compact(newSliceList)
}

func MapToValues[K comparable, V any](input map[K]V) []V {
	values := make([]V, 0, len(input))
	for _, v := range input {
		values = append(values, v)
	}
	return values
}

func MapToKeys[K comparable, V any](input map[K]V) []K {
	keys := make([]K, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

package slice

import (
	"cmp"
	"golangutils/pkg/conv"
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

func ObjArrayToStringBySep[T any](arr []T, sep string) string {
	dataArr := []string{}
	for _, data := range arr {
		dataArr = append(dataArr, conv.ToString(data))
	}
	return ArrayToStringBySep(dataArr, sep)
}

func ObjArrayToString[T any](arr []T) string {
	return ObjArrayToStringBySep(arr, " ")
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

func MapExistKey[K comparable, V any](input map[K]V, searchKey K) bool {
	_, exists := input[searchKey]
	return exists
}

func MapExistValue[K comparable, V comparable](input map[K]V, searchValue V) bool {
	for _, v := range input {
		if v == searchValue {
			return true
		}
	}
	return false
}

func IsMapEmpty[K comparable, V any](input map[K]V) bool {
	return len(input) == 0
}

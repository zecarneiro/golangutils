package slice

import (
	"fmt"
	"strings"

	"golangutils/pkg/conv"
	"golangutils/pkg/str"
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

func RemoveDuplicateByOrder[T comparable](sliceList []T, isFromStart bool) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(sliceList))
	if len(sliceList) == 0 {
		return result
	}
	if isFromStart {
		for _, val := range sliceList {
			if !seen[val] {
				result = append(result, val)
				seen[val] = true
			}
		}
	} else {
		for i := len(sliceList) - 1; i >= 0; i-- {
			val := sliceList[i]
			if !seen[val] {
				result = append([]T{val}, result...)
				seen[val] = true
			}
		}
	}
	return result
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	return RemoveDuplicateByOrder(sliceList, true)
}

func MapToValues[K comparable, V any](input map[K]V) []V {
	values := make([]V, 0, len(input))
	if input == nil {
		return values
	}
	for _, v := range input {
		values = append(values, v)
	}
	return values
}

func MapToKeys[K comparable, V any](input map[K]V) []K {
	keys := make([]K, 0, len(input))
	if input == nil {
		return keys
	}
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

func MapExistKey[K comparable, V any](input map[K]V, searchKey K) bool {
	if input == nil {
		return false
	}
	_, exists := input[searchKey]
	return exists
}

func MapExistValue[K comparable, V comparable](input map[K]V, searchValue V) bool {
	if input == nil {
		return false
	}
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

func ConcatMap[T comparable](map1, map2 map[T]any) map[T]any {
	result := make(map[T]any)
	for k, v := range map1 {
		result[k] = v
	}
	for k, v := range map2 {
		result[k] = v
	}
	return result
}

func RemoveAllEmpty(input []string) []string {
	newInput := []string{}
	for _, data := range input {
		if !str.IsEmpty(data) {
			newInput = append(newInput, data)
		}
	}
	return newInput
}

func MapToString[T comparable, V any](input map[T]V) string {
	data := ""
	for k, v := range input {
		if str.IsEmpty(data) {
			data = fmt.Sprintf(`%s=%s`, conv.ToString(k), conv.ToString(v))
		}
	}
	return data
}

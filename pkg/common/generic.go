package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TaskFunc func()

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

func FilterArray[T any](array []T, fun func(T) bool) []T {
	var newArr []T
	for _, element := range array {
		if fun(element) {
			newArr = append(newArr, element)
		}
	}
	return newArr
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ObjectToString(data any) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	return string(jsonData), err
}

func ObjectToStringEscapeHtml(data any) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func StringToObject[T any](data string) (T, error) {
	var object T
	err := json.Unmarshal([]byte(data), &object)
	return object, err
}

func IsNil[T any](arg *T) bool {
	return arg == nil
}

func Ternary[T any](condition bool, resultTrue, resultFalse T) T {
	if condition {
		return resultTrue
	}
	return resultFalse
}

func StringReplaceAll(data string, replacer map[string]string) string {
	var newData = data
	if len(replacer) > 0 {
		for key, value := range replacer {
			newData = strings.Replace(newData, key, value, -1)
		}
	}
	return newData
}

func GetSubstring(str string, start int, end int) string {
	newStr := ""
	index := 0
	for _, character := range str {
		if index > end {
			break
		}
		if index >= start {
			byteChar := byte(character)
			newStr += string(byteChar)
		}
		index++
	}
	return newStr
}

func StringContains(original string, search string, isIgnoreCase bool) bool {
	if isIgnoreCase {
		return strings.Contains(strings.ToLower(original), strings.ToLower(search))
	}
	return strings.Contains(original, search)
}

func HasLastChar(data string, checkChar string) bool {
	if len(data) <= 0 {
		return false
	}
	return strings.HasSuffix(data, checkChar)
}

func DeleteLastChar(original string) string {
	if len(original) > 0 {
		return original[:len(original)-1]
	}
	return ""
}

func StringToInt(data string) (int, error) {
	return strconv.Atoi(data)
}

func IntToString(data int) string {
	return strconv.Itoa(data)
}

// Example: map[string]interface{}{"FUNC_NAME": FUNC, "FUNC_NAME_1": FUNC_1, ....}
func FuncCall[T interface{}](caller interface{}, params ...interface{}) (T, error) {
	var in []reflect.Value = []reflect.Value{}
	var result T
	var err error
	funcRef := reflect.ValueOf(caller)
	if len(params) > 0 {
		if len(params) != funcRef.Type().NumIn() {
			err = errors.New("The number of params is out of index.")
		}
	}
	if err == nil {
		in = make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		res := funcRef.Call(in)
		if res != nil {
			result = res[0].Interface().(T)
		}
	}
	return result, err
}

func IsValidJSON(s string) bool {
	var js any
	return json.Unmarshal([]byte(s), &js) == nil
}

func BytesToMB(b int64) float64 {
	return float64(b) / 1024 / 1024
}

func Sleep(second int) {
	time.Sleep(time.Second * time.Duration(second))
}

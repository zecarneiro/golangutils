package jnoronha_golangutils

import (
	"encoding/json"
	"errors"
	"io"
	"jnoronha_golangutils/entities"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
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

func FilterArray[T any](array []T, fun func(T) bool) []T {
	var newArr []T
	for _, element := range array {
			if fun(element) {
				newArr = append(newArr, element)
			}
	}
	return newArr
}

func ObjectToString(data any) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	return string(jsonData), err
}

func StringToObject[T any](data string) (T, error) {
	var object T
	err := json.Unmarshal([]byte(data), &object)
	return object, err
}

func IsNil[T any](arg *T) bool {
	return arg == nil
}

func ProcessError(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func Download(url string, destFile string) entities.Response[bool] {
	// Create a GET request to fetch the file
	response, err := http.Get(url)
	if err != nil {
		return entities.Response[bool]{Data: false, Error: err}
	}
	defer response.Body.Close()

	// Create the file to which the downloaded content will be written
	file, err := os.Create(destFile)
	if err != nil {
		return entities.Response[bool]{Data: false, Error: err}
	}
	defer file.Close()

	// Copy the response body (file content) to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return entities.Response[bool]{Data: false, Error: err}
	}
	return entities.Response[bool]{Data: true}
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

// Example: map[string]interface{}{"FUNC_NAME": FUNC, "FUNC_NAME_1": FUNC_1, ....}
func FuncCall[T interface{}](funcName string, callerStorage map[string]interface{}, params ... interface{}) (T, error) {
	var in []reflect.Value = []reflect.Value{}
	var result T
	var err error
	funcRef := reflect.ValueOf(callerStorage[funcName])
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

func HasInternet() bool {
	timeout := time.Duration(5000 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	//default url to check connection is http://google.com
	_, err := client.Get("https://google.com")
	if err != nil {
		return false
	}
	return true
}

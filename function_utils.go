package jnoronha_golangutils

import (
	"encoding/json"
	"fmt"
	"io"
	"jnoronha_golangutils/entities"
	"log"
	"net/http"
	"os"
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

func Confirm(message string) bool {
	fmt.Printf("%s [y/n]: ", message)
	var response string
	fmt.Scanln(&response)
	if response == "Y" || response == "y" {
		return true
	}
	return false
}

package obj

import (
	"bytes"
	"encoding/json"
)

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
	var target T
	err := json.Unmarshal([]byte(data), &target)
	if err != nil {
		return *new(T), err
	}
	return target, nil
}

func IsValidJSON(s string) bool {
	var js any
	return json.Unmarshal([]byte(s), &js) == nil
}

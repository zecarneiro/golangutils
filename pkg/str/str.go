package str

import (
	"fmt"
	"strings"
	"unicode"
)

func StringReplaceAll(data string, replacer map[string]string) string {
	newData := data
	if len(replacer) > 0 {
		for key, value := range replacer {
			newData = strings.ReplaceAll(newData, key, value)
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

func Contains(original string, search string, caseInsensitive bool) bool {
	if caseInsensitive {
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

func IsEmpty(data string) bool {
	data = strings.Trim(data, " ")
	return len(data) == 0 || data == " "
}

func GetInDoubleQuotes(data string) string {
	if IsEmpty(data) {
		return data
	}
	return fmt.Sprintf("\"%s\"", data)
}

func GetInSingleQuotes(data string) string {
	if IsEmpty(data) {
		return data
	}
	return fmt.Sprintf("'%s'", data)
}

func ToCamelCase(data string, keepFirstWord bool) string {
	words := strings.Fields(data)
	if len(words) <= 1 {
		return data
	}
	newWord := ""
	for index, word := range words {
		if index == 0 {
			if !keepFirstWord {
				word = strings.ToLower(word)
			}
		} else {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			word = string(runes)
		}
		newWord = fmt.Sprintf("%s%s", newWord, word)
	}
	return newWord
}

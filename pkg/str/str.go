package str

import "strings"

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

func IsEmpty(data string) bool {
	data = strings.Trim(data, " ")
	return len(data) == 0 || data == " "
}

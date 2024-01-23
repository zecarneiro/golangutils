package jnoronha_golangutils

import "fmt"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"
)

func LogLog(data string, keepLine bool) {
	if keepLine {
		fmt.Print(data)
	} else {
		fmt.Println(data)
	}
}

func DebugLog(data string, keepLine bool) {
	LogLog("[DEBUG] ", true)
	LogLog(data, keepLine)
}

func WarnLog(data string, keepLine bool) {
	LogLog(fmt.Sprintf("[%sWARN%s] ", string(ColorYellow), string(ColorReset)), true)
	LogLog(data, keepLine)
}

func ErrorLog(data string, keepLine bool) {
	LogLog(fmt.Sprintf("[%sERROR%s] ", string(ColorRed), string(ColorReset)), true)
	LogLog(data, keepLine)
}

func InfoLog(data string, keepLine bool) {
	LogLog(fmt.Sprintf("[%sINFO%s] ", string(ColorBlue), string(ColorReset)), true)
	LogLog(data, keepLine)
}

func SuccessLog(data string, keepLine bool) {
	LogLog(fmt.Sprintf("[%sOK%s] ", string(ColorGreen), string(ColorReset)), true)
	LogLog(data, keepLine)
}

func PromptLog(data string) {
	LogLog(fmt.Sprintf("%s>>>%s ", string(ColorGray), string(ColorReset)), true)
	LogLog(data, false)
}

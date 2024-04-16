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

func OkLog(data string, keepLine bool) {
	LogLog(fmt.Sprintf("[%sOK%s] ", string(ColorGreen), string(ColorReset)), true)
	LogLog(data, keepLine)
}

func PromptLog(data string) {
	LogLog(fmt.Sprintf("%s>>>%s ", string(ColorGray), string(ColorReset)), true)
	LogLog(data, false)
}

func TitleLog(data string) {
	messageLen := len(data)
	separator := ""
	for i := 1; i < messageLen + 8; i++ {
		separator = "#" + separator
	}
	LogLog(separator, false)
	LogLog("## " + data + " ##", false)
	LogLog(separator, false)
}

func HeaderLog(data string) {
	LogLog("## " + data + " ##", false)
}

func Separatorlog(length int) {
	data := "# "
	if (length < 6) {
		length = 6
	}
	for i := 1; i < length - 4; i++ {
		data += "-"
	}
	data += " #"
	LogLog(data, false)
}

func LogHelp(appName string, description string, usages []string,  args []string, others []string) {
	LogLog(appName + " - " + description, false)
	if len(usages) > 0 {
		LogLog("USAGE: ", false)
		for _, usage := range usages {
			LogLog("\t" + appName + " " + usage, false)
		}
	}
	if len(args) > 0 {
		LogLog("ARGS: ", false)
		args = append(args, "--help|-h: Show Help")
		for _, arg := range args {
			LogLog("\t" + arg, false)
		}
	}
	if len(others) > 0 {
		for _, other := range others {
			LogLog("\t" + other, false)
		}
	}
}

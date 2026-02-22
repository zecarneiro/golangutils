package logger

import (
	"fmt"

	"golangutils/pkg/enums"
)

func FormatDataWithColor(data string, color string) string {
	return fmt.Sprintf("%s%s%s", color, data, enums.Reset.String())
}

func Log(data any) {
	logToFile(data, "")
	print(data)
}

func Debug(data any) {
	logToFile(data, "debug")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", "DEBUG"))
	print(data)
}

func Warn(data any) {
	logToFile(data, "warn")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", FormatDataWithColor("WARN", enums.Yellow.String())))
	print(data)
}

func Error(data error) {
	if data != nil {
		dataStr := fmt.Sprintf("%v", data)
		logToFile(dataStr, "error")
		WithKeepLine(true)
		print(fmt.Sprintf("[%s] ", FormatDataWithColor("ERROR", enums.Red.String())))
		print(dataStr)
	}
}

func ErrorStr(data string) {
	Error(fmt.Errorf("%s", data))
}

func Info(data any) {
	logToFile(data, "info")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", FormatDataWithColor("INFO", enums.Blue.String())))
	print(data)
}

func Ok(data any) {
	logToFile(data, "ok")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", FormatDataWithColor("OK", enums.Green.String())))
	print(data)
}

func Prompt(data string) {
	logToFile(data, "prompt")
	WithKeepLine(true)
	print(fmt.Sprintf("%s ", FormatDataWithColor(">>>", enums.Gray.String())))
	print(data)
}

func Title(data string) {
	messageLen := len(data)
	prefix := "##    "
	suffix := "    ##"
	plusLoop := len(prefix) + len(suffix) + 1
	separator := ""
	for i := 1; i < messageLen+plusLoop; i++ {
		separator = "#" + separator
	}
	print(separator)
	print(fmt.Sprintf("%s%s%s", prefix, data, suffix))
	print(separator)
}

func Header(data string) {
	length := len(data) + headerLength
	data = fmt.Sprintf(" %s ", data)
	if len(data) < length-2 {
		newLength := length - len(data)
		middle := newLength / 2
		if newLength%2 != 0 {
			middle++
		}
		for i := 1; i < middle-1; i++ {
			data = fmt.Sprintf("-%s-", data)
			if len(data) >= length-2 {
				break
			}
		}
	}
	print(fmt.Sprintf("#%s#", data))
}

func Separator() {
	data := "# "
	if separatorLength < 6 {
		separatorLength = 6
	}
	for i := 1; i < separatorLength-4; i++ {
		data += "-"
	}
	data += " #"
	print(data)
}

func Help(appName string, description string, usages []string, args []string, others []string) {
	print(appName + " - " + description)
	if len(usages) > 0 {
		print("USAGE: ")
		for _, usage := range usages {
			print("\t" + appName + " " + usage)
		}
	}
	if len(args) > 0 {
		print("ARGS: ")
		args = append(args, "--help|-h: Show Help")
		for _, arg := range args {
			print("\t" + arg)
		}
	}
	if len(others) > 0 {
		for _, other := range others {
			print("\t" + other)
		}
	}
}

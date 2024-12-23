package golangutils

import (
	"fmt"
	"golangutils/enum"
)

const ()

type LoggerUtils struct {
	keepLine bool
}

func NewLoggerUtils() LoggerUtils {
	return LoggerUtils{keepLine: false}
}

func (l *LoggerUtils) disableKeepLine() {
	l.keepLine = false
}

func (l *LoggerUtils) EnableKeepLine() {
	l.keepLine = true
}

func FormatDataWithColor(data string, color string) string {
	return fmt.Sprintf("%s%s%s", string(color), data, string(enum.COLOR_RESET))
}

func (l *LoggerUtils) Log(data string) {
	if l.keepLine {
		fmt.Print(data)
	} else {
		fmt.Println(data)
	}
	l.disableKeepLine()
}

func (l *LoggerUtils) Debug(data string) {
	l.EnableKeepLine()
	l.Log("[DEBUG] ")
	l.Log(data)
}

func (l *LoggerUtils) Warn(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%s] ", FormatDataWithColor("WARN", enum.COLOR_YELLOW)))
	l.Log(data)
}

func (l *LoggerUtils) Error(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%s] ", FormatDataWithColor("ERROR", enum.COLOR_RED)))
	l.Log(data)
}

func (l *LoggerUtils) Info(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%s] ", FormatDataWithColor("INFO", enum.COLOR_BLUE)))
	l.Log(data)
}

func (l *LoggerUtils) Ok(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("[%s] ", FormatDataWithColor("OK", enum.COLOR_GREEN)))
	l.Log(data)
}

func (l *LoggerUtils) Prompt(data string) {
	l.EnableKeepLine()
	l.Log(fmt.Sprintf("%s ", FormatDataWithColor(">>>", enum.COLOR_GRAY)))
	l.Log(data)
}

func (l *LoggerUtils) Title(data string) {
	messageLen := len(data)
	prefix := "##    "
	suffix := "    ##"
	plusLoop := len(prefix) + len(suffix) + 1
	separator := ""
	for i := 1; i < messageLen+plusLoop; i++ {
		separator = "#" + separator
	}
	l.Log(separator)
	l.Log(fmt.Sprintf("%s%s%s", prefix, data, suffix))
	l.Log(separator)
}

func (l *LoggerUtils) Header(data string, length int) {
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
	l.Log(fmt.Sprintf("#%s#", data))
}

func (l *LoggerUtils) Separator(length int) {
	data := "# "
	if length < 6 {
		length = 6
	}
	for i := 1; i < length-4; i++ {
		data += "-"
	}
	data += " #"
	l.Log(data)
}

func (l *LoggerUtils) Help(appName string, description string, usages []string, args []string, others []string) {
	l.Log(appName + " - " + description)
	if len(usages) > 0 {
		l.Log("USAGE: ")
		for _, usage := range usages {
			l.Log("\t" + appName + " " + usage)
		}
	}
	if len(args) > 0 {
		l.Log("ARGS: ")
		args = append(args, "--help|-h: Show Help")
		for _, arg := range args {
			l.Log("\t" + arg)
		}
	}
	if len(others) > 0 {
		for _, other := range others {
			l.Log("\t" + other)
		}
	}
}

func (l *LoggerUtils) HasAndLogError(err error) bool {
	if err != nil {
		l.Error(err.Error())
		return true
	}
	return false
}

package golangutils

import (
	"fmt"
	"golangutils/enum"
	"log"
	"os"
	"strings"
	"time"
)

type LoggerUtils struct {
	keepLine bool
	logFile  string
}

func NewLoggerUtils() LoggerUtils {
	return LoggerUtils{
		keepLine: false,
	}
}

func (l *LoggerUtils) logToFile(data string, logType string) {
	if len(l.logFile) > 0 {
		f, err := os.OpenFile(l.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		f.WriteString(l.GetTimeAndDataForLogFile(data, logType))
		f.Sync()
	}
}

func (l *LoggerUtils) disableKeepLine() {
	l.keepLine = false
}

func (l *LoggerUtils) EnableKeepLine() {
	l.keepLine = true
}

func (l *LoggerUtils) GetTimeAndDataForLogFile(data string, logType string) string {
	t := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d  %s  %s\n", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), strings.ToUpper(logType), data)
}

func FormatDataWithColor(data string, color string) string {
	return fmt.Sprintf("%s%s%s", string(color), data, string(enum.COLOR_RESET))
}

func (l *LoggerUtils) log(data string) {
	if l.keepLine {
		fmt.Print(data)
	} else {
		fmt.Println(data)
	}
	l.disableKeepLine()
}

func (l *LoggerUtils) SetLogFile(logFile string) {
	l.logFile = logFile
}

func (l *LoggerUtils) Log(data string) {
	l.logToFile(data, "")
	l.log(data)
}

func (l *LoggerUtils) Debug(data string) {
	l.logToFile(data, "debug")
	l.EnableKeepLine()
	l.log("[DEBUG] ")
	l.log(data)
}

func (l *LoggerUtils) Warn(data string) {
	l.logToFile(data, "warn")
	l.EnableKeepLine()
	l.log(fmt.Sprintf("[%s] ", FormatDataWithColor("WARN", enum.COLOR_YELLOW)))
	l.log(data)
}

func (l *LoggerUtils) Error(data string) {
	l.logToFile(data, "error")
	l.EnableKeepLine()
	l.log(fmt.Sprintf("[%s] ", FormatDataWithColor("ERROR", enum.COLOR_RED)))
	l.log(data)
}

func (l *LoggerUtils) Info(data string) {
	l.logToFile(data, "info")
	l.EnableKeepLine()
	l.log(fmt.Sprintf("[%s] ", FormatDataWithColor("INFO", enum.COLOR_BLUE)))
	l.log(data)
}

func (l *LoggerUtils) Ok(data string) {
	l.logToFile(data, "ok")
	l.EnableKeepLine()
	l.log(fmt.Sprintf("[%s] ", FormatDataWithColor("OK", enum.COLOR_GREEN)))
	l.log(data)
}

func (l *LoggerUtils) Prompt(data string) {
	l.logToFile(data, "prompt")
	l.EnableKeepLine()
	l.log(fmt.Sprintf("%s ", FormatDataWithColor(">>>", enum.COLOR_GRAY)))
	l.log(data)
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
	l.log(separator)
	l.log(fmt.Sprintf("%s%s%s", prefix, data, suffix))
	l.log(separator)
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
	l.log(fmt.Sprintf("#%s#", data))
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
	l.log(data)
}

func (l *LoggerUtils) Help(appName string, description string, usages []string, args []string, others []string) {
	l.log(appName + " - " + description)
	if len(usages) > 0 {
		l.log("USAGE: ")
		for _, usage := range usages {
			l.log("\t" + appName + " " + usage)
		}
	}
	if len(args) > 0 {
		l.log("ARGS: ")
		args = append(args, "--help|-h: Show Help")
		for _, arg := range args {
			l.log("\t" + arg)
		}
	}
	if len(others) > 0 {
		for _, other := range others {
			l.log("\t" + other)
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

package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golangutils/pkg/common"
	"log"
	"os"
	"strings"
	"time"
)

func logToFile(data string, logType string) {
	if len(logFile) > 0 {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		f.WriteString(getTimeAndDataForLogFile(data, logType))
		f.Sync()
	}
}

func getTimeAndDataForLogFile(data string, logType string) string {
	t := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d  %s  %s\n", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), strings.ToUpper(logType), data)
}

func print(data string) {
	if keepLine {
		fmt.Print(data)
	} else {
		fmt.Println(data)
	}
	WithKeepLine(false)
}

func FormatDataWithColor(data string, color string) string {
	return fmt.Sprintf("%s%s%s", color, data, common.Reset.String())
}

func Log(data string) {
	logToFile(data, "")
	print(data)
}

func Debug(data string) {
	logToFile(data, "debug")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", "DEBUG"))
	print(data)
}

func Warn(data string) {
	logToFile(data, "warn")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", FormatDataWithColor("WARN", common.Yellow.String())))
	print(data)
}

func Error(data any) {
	var dataStr string
	if data == nil {
		log.Fatal("Receive nil")
	}
	switch v := data.(type) {
	case error:
		dataStr = v.Error()
	case string:
		dataStr = v
	default:
		dataStr = fmt.Sprintf("%v", v)
	}
	logToFile(dataStr, "error")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", FormatDataWithColor("ERROR", common.Red.String())))
	print(dataStr)
}

func Info(data string) {
	logToFile(data, "info")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", FormatDataWithColor("INFO", common.Blue.String())))
	print(data)
}

func Ok(data string) {
	logToFile(data, "ok")
	WithKeepLine(true)
	print(fmt.Sprintf("[%s] ", FormatDataWithColor("OK", common.Green.String())))
	print(data)
}

func Prompt(data string) {
	logToFile(data, "prompt")
	WithKeepLine(true)
	print(fmt.Sprintf("%s ", FormatDataWithColor(">>>", common.Gray.String())))
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

func Json(data string) {
	var formatado bytes.Buffer
	err := json.Indent(&formatado, []byte(data), "", "  ")
	if err == nil {
		Error(err)
	} else {
		print(formatado.String())
	}
}

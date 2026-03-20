package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golangutils/pkg/conv"
)

func getTimeAndDataForLogFile(data string, logType string) string {
	t := time.Now()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d  %s  %s\n", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), strings.ToUpper(logType), data)
}

func logToFile(data any, logType string) {
	dataStr := conv.ToString(data)
	if len(logFile) > 0 {
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		f.WriteString(getTimeAndDataForLogFile(dataStr, logType))
		f.Sync()
	}
}

func print(data any) {
	dataStr := conv.ToString(data)
	if keepLine {
		fmt.Print(dataStr)
	} else {
		fmt.Println(dataStr)
	}
	WithKeepLine(false)
}
